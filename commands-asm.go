package main

import(
	"time"
	"strconv"
	"strings"
	"math/rand"
	"encoding/hex"

	"github.com/bwmarrin/discordgo"
	"github.com/keystone-engine/keystone/bindings/go/keystone"
	"github.com/bnagy/gapstone"
)

// Returns the proper keystone architecture based on the user input string
func parseArchitectureKeystone(arch string) (keystone.Architecture, keystone.Mode) {
	switch arch {
	case "x86":
		return keystone.ARCH_X86, keystone.MODE_32
	case "x64", "x86_64", "x86-64":
		return keystone.ARCH_X86, keystone.MODE_64
	case "arm":
		return keystone.ARCH_ARM, keystone.MODE_ARM
	case "aarch64", "arm64":
		return keystone.ARCH_ARM64, keystone.MODE_LITTLE_ENDIAN
	case "ppc", "ppc32":
		return keystone.ARCH_PPC, keystone.MODE_PPC32 | keystone.MODE_BIG_ENDIAN
	case "ppc64":
		return keystone.ARCH_PPC, keystone.MODE_PPC64
	case "mips", "mips32":
		return keystone.ARCH_MIPS, keystone.MODE_MIPS32 | keystone.MODE_BIG_ENDIAN
	case "mips64":
		return keystone.ARCH_MIPS, keystone.MODE_MIPS64
	default:
		return 0, 0
	}
}

// Returns the proper capstone architecture based on the user input string
func parseArchitectureCapstone(arch string) (int, int) {
	switch arch {
	case "x86":
		return gapstone.CS_ARCH_X86, gapstone.CS_MODE_32
	case "x64", "x86_64", "x86-64":
		return gapstone.CS_ARCH_X86, gapstone.CS_MODE_64
	case "arm":
		return gapstone.CS_ARCH_ARM, gapstone.CS_MODE_ARM
	case "aarch64", "arm64":
		return gapstone.CS_ARCH_ARM64, gapstone.CS_MODE_ARM
	case "ppc", "ppc32":
		return gapstone.CS_ARCH_PPC, gapstone.CS_MODE_BIG_ENDIAN
	case "ppc64":
		return gapstone.CS_ARCH_PPC, gapstone.CS_MODE_LITTLE_ENDIAN
	case "mips", "mips32":
		return gapstone.CS_ARCH_MIPS, gapstone.CS_MODE_MIPS32 | gapstone.CS_MODE_BIG_ENDIAN
	case "mips64":
		return gapstone.CS_ARCH_MIPS, gapstone.CS_MODE_MIPS64 | gapstone.CS_MODE_LITTLE_ENDIAN
	default:
		return 0, 0
	}
}

// Assembles the given instructions into opcodes via the given architecture
func cmdAssemble(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	asmArch  	 := args[1]
	instructions 	 := ""

	// Stitch together the rest of the arguments for the instructions
	if len(args) > 2 {
		for i := 2; i < len(args); i++ {
			instructions += args[i] + " "
		}
	}

	if arch, mode := parseArchitectureKeystone(asmArch); arch != 0 && mode != 0 {
		// Output Message
		outMsg := "Assembly: ```\n"

		// Split by ';' deliminator (termination character in assembly)
		ins := strings.Split(instructions, ";")

		// Use the keystone library for assembly
		if ks, err := keystone.New(arch, mode); err == nil {
			defer ks.Close()

			// Get each instruction's opcodes individually to format nicely
			for _, i := range ins {
				// Use intel syntax for x86 because AT&T syntax is ugly
				if arch == keystone.ARCH_X86 {
					if err := ks.Option(keystone.OPT_SYNTAX, keystone.OPT_SYNTAX_INTEL); err != nil {
						s.ChannelMessageSend(m.ChannelID, "Failed to set keystone option")
						return
					}
				}

				if ops, _, ok := ks.Assemble(i, 0); ok {
					opcodes := ""

					// Encode a string of the bytes
					for _, op := range ops {
						// Format to hex representation, and pad to 2 chars.
						opcodes += padLeft(strconv.FormatInt(int64(op), 16), "0", 2) + " "
					}

					// Allow some space between the opcodes and instructions
					outMsg += padRight(opcodes, " ", 30) + i + "\n"
				} else {
					// Keystone assembler failed
					s.ChannelMessageSend(m.ChannelID, "Could not assemble the given assembly. Are the instructions valid?")
					return
				}
			}

			// Keystone assembler succeeded, give the user the output
			s.ChannelMessageSend(m.ChannelID, outMsg + "```")
			return
		}

		// If we reached this point, it's because keystone's engine failed to initialize
		s.ChannelMessageSend(m.ChannelID, "Keystone engine is not working! :(")
	} else {
		supportedArchs := "```"
		supportedArchs += "x86, x86_64/x64, arm, arm64/aarch64, ppc/ppc32, ppc64, mips/mips32, mips64"
		supportedArchs += "```"

		s.ChannelMessageSend(m.ChannelID, "Architecture not supported! Supported architectures: " + supportedArchs)
	}
}

// Disassembles the given opcodes into instructions via the architecture
func cmdDisassemble(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	asmArch  	 := args[1]
	opcodes  	 := ""

	// Stitch together the rest of the arguments for the opcodes
	if len(args) > 2 {
		for i := 2; i < len(args); i++ {
			opcodes += args[i]
		}
	}
	
	// Allow some flexibility in input (ie. allow 0x, ;)
	opcodes = strings.Replace(opcodes, ";", "")
	opcodes = strings.Replace(opcodes, "0x", "")

	if arch, mode := parseArchitectureCapstone(asmArch); arch != 0 && mode != 0 {
		// Output Message
		outMsg := "Disassembly: ```\n"

		// Use the keystone library for assembly
		if gs, err := gapstone.New(arch, uint(mode)); err == nil {
			defer gs.Close()

			// Use intel syntax for x86 because AT&T syntax is ugly
			if arch == gapstone.CS_ARCH_X86 {
				if err := gs.SetOption(gapstone.CS_OPT_SYNTAX, gapstone.CS_OPT_SYNTAX_INTEL); err != nil {
					s.ChannelMessageSend(m.ChannelID, "Failed to set gapstone option")
					return
				}
			}

			// We need to decode the string as capstone only accepts raw binary data for input
			if opcodesBinary, err := hex.DecodeString(opcodes); err == nil {
				if ins, err := gs.Disasm(opcodesBinary, 0, 0); err == nil {
					for _, i := range ins {
						instructionOpCodes := ""
						ops := i.Bytes

						// Encode a string of the bytes
						for _, op := range ops {
							instructionOpCodes += padLeft(strconv.FormatInt(int64(op), 16), "0", 2) + " "
						}

						// Beautify the output
						outMsg += padRight(instructionOpCodes, " ", 30) + i.Mnemonic + " " + i.OpStr + "\n"
					}
				} else {
					// Capstone disassembler failed
					s.ChannelMessageSend(m.ChannelID, "Could not disassemble the given opcodes. Are the opcodes valid?")
					return
				}

				// Disassembler succeeded, give the user the output
				s.ChannelMessageSend(m.ChannelID, outMsg + "```")
				return
			} else {
				// Failed to decode the string into raw binary data - must be invalid hex
				s.ChannelMessageSend(m.ChannelID, "Invalid opcodes.")
				return
			}
		}

		// If we reached this point, it's because capstone's engine failed to initialize
		s.ChannelMessageSend(m.ChannelID, "Capstone engine is not working! :(")
	} else {
		// Unsupported architecture

		supportedArchs := "```"
		supportedArchs += "x86, x86_64/x64, arm, arm64/aarch64, ppc/ppc32, ppc64, mips/mips32, mips64"
		supportedArchs += "```"

		s.ChannelMessageSend(m.ChannelID, "Architecture not supported! Supported architectures: " + supportedArchs)
	}
}

// Gives a random reverse engineering trick
func cmdReTrick(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	rand.Seed(time.Now().Unix())

	tricks := []string {
		"When possible, use a debugger to trace user input in a function",
		"Viewing strings is very helpful",
		"IDA: IDA has a quick action dropdown to the right of the breakdown bar https://i.imgur.com/TmtXE1O.png",
		"IDA: You can view functions that call a target function and functions the target function calls with View -> Open subviews -> Function calls https://i.imgur.com/bdB0Rge.png",
		"IDA: Hit 'k' to convert 'rbp+var_xxx' format in instructions into 'rbp-xxxh'",
		"IDA: When in Graph View, the 'Graph overview' window can be used to quickly navigate around large functions https://i.imgur.com/8IqPs1r.png",
		"Intel x86 can be tricky, `mov rax, rax` may seem like a NOP, but it also implicitly clears the upper 32-bits of the register",
		"Intel x86 can be tricky, `cmpxchg` instructions implicitly modify the value of the RAX register, regardless of operands",
		"When you see instructions that check the value of one offset from a register, then the register is set to a value from another offset in a loop - it's probably a linked list",
	}

	// Pick a random one
	n := rand.Int() % len(tricks)
	s.ChannelMessageSend(m.ChannelID, tricks[n])
}

// Gives a random exploit development trick
func cmdExploitTrick(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	rand.Seed(time.Now().Unix())

	tricks := []string {
		"Use infloop gadgets in ROP chains for blind debugging",
		"For use-after-free exploits, try empty heap spraying after code execution to stabilize the process if it's a critical object",
		"The power of going straight from a bug to PC/IP control is overrated, other primitives like arbitrary R/W are often easier and just as powerful",
		"When writing shellcode, use `xor [reg], [reg]` to avoid null bytes. This can also be used for patches, as `xor` instructions typically use less opcodes than `mov` instructions.",
		"When the patch is smaller than the original code, NOPS are your friend",
		"You don't always need a seperate infoleak bug to defeat ASLR, sometimes you can use the context of other registers to calculate from.",
		"When looking for vulnerabilities in large software, map out attack surface first - it sucks to find a bug then realize it's in code that's unreachable later on.",
		"When looking for integer overflows in x86, `ja/jump above or jb/jumb below` is an unsigned compare, `jg/jump greater or jl/jump lower` is a signed compare.",
		"If you're using gdb to debug exploits, use PEDA https://github.com/longld/peda",
	}

	// Pick a random one
	n := rand.Int() % len(tricks)
	s.ChannelMessageSend(m.ChannelID, tricks[n])
}

// Gives a PDF link to the manual for the given architecture
func cmdManual(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var url string
	asmArgs := args[1]

	if asmArgs == "x86" || asmArgs == "x64" || asmArgs == "x86_64" || asmArgs == "x86-64" {
		url = "https://www.intel.com/content/dam/www/public/us/en/documents/manuals/64-ia-32-architectures-software-developer-instruction-set-reference-manual-325383.pdf"
	} else if asmArgs == "arm" || asmArgs == "aarch64" || asmArgs == "arm64" {
		url = "https://static.docs.arm.com/ddi0487/ca/DDI0487C_a_armv8_arm.pdf"
	} else if asmArgs == "ppc" || asmArgs == "ppc32" || asmArgs == "ppc64" {
		url = "http://www.plantation-productions.com/Webster/www.writegreatcode.com/Vol2/wgc2_OB.pdf"
	} else if asmArgs == "mips" || asmArgs == "mips32" || asmArgs == "mips64" {
		url = "https://www.cs.cmu.edu/afs/cs/academic/class/15740-f97/public/doc/mips-isa.pdf"
	} else {
		supportedArchs := "```"
		supportedArchs += "x86, x86_64/x64, arm, arm64/aarch64, ppc/ppc32, ppc64, mips/mips32, mips64"
		supportedArchs += "```"
		s.ChannelMessageSend(m.ChannelID, "Architecture not supported! Supported architectures: " + supportedArchs)
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Here you go: " + url)
}
