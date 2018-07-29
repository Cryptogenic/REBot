# REBot
[![Release Mode](https://img.shields.io/badge/Release%20Mode-Stable-green.svg)]()  [![Maintenance](https://img.shields.io/badge/Maintained%3F-Partially-yellow.svg)]()  [![Version](https://img.shields.io/badge/Version-1.0-brightgreen.svg)]()

REBot is a Discord bot programmed in the Golang programming language to provide useful commands to reverse engineers and exploit developers. It provides features such as on-the-fly assembly and disassembly for common (and even less common) architectures, such as x86/64, ARM/AARCH64, PPC, and MIPS. It also provides other features such as a technical dictionary, CVE look-up, and giving tips and tricks on reverse engineering and exploit development practices.

You can join the official REBot to your server using [this discord invite link](https://discordapp.com/oauth2/authorize?client_id=472921462328524831&permissions=0&scope=bot). REBot doesn't need any special permissions - only text read and send permissions so it can interact with you.

## Getting Started
Below are some instructions on how to get the project building on your machine. I personally built this on Windows 10. The project itself is platform independent - however the way you build and configure keystone and capstone will vary depending on your target system.

### Prerequisites
The following software is required to built and use REBot.
- Golang
- Keystone Assembler Engine
- Capstone Disassembler Engine

### Go Dependencies
The following dependencies are required to build the project using Go.
- [bwmarrin/discordgo](http://github.com/bwmarrin/discordgo)
- [go-ini/ini](http://github.com/go-ini/ini)
- [keystone go bindings](http://github.com/keystone-engine/keystone/bindings/go/keystone)
- [gapstone - capstone go bindings](http://github.com/bnagy/gapstone)

## Building
### Installing prerequisites
#### Windows
You should install Golang using their official Windows installer provided [here](https://golang.org/dl/). As for Keystone and Capstone, for Windows I recommend that you use the pre-compiled libraries they provide [here](https://github.com/keystone-engine/keystone/releases/download/0.9.1/keystone-0.9.1-win64.zip) (for keystone) and [here](https://github.com/aquynh/capstone/releases/download/3.0.5/capstone-3.0.5-win64.zip) (for capstone).

The libraries and header files for Keystone and Capstone should go in dedicated folders, I personally chose `C:\src\keystone` and `C:\src\capstone` respectively. The `.lib` files can be disgarded/deleted/renamed, as if they're present in the folder when you go to build, you may encounter linker errors. The `.dll` file is sufficient.

You may also need to edit the bindings to build the project successfully - changing the `cgo` directives to reflect the directories you setup the libraries in. For Keystone, I edited `go/src/github.com/keystone-engine/keystone/bindings/go/keystone/keystone-binding.go` and set the `cgo` directives to the following:

```golang
//#cgo windows CFLAGS: -IC:/src/keystone/include
//#cgo windows LDFLAGS: -LC:/src/keystone -lkeystone -lstdc++ -lm
```

For Capstone, nearly all the `.go` files for the binding had to be edited for the `cgo` directives. I edited them to the following:

```golang
//#cgo windows CFLAGS: -IC:/src/capstone/include
//#cgo windows LDFLAGS: -LC:/src/capstone -lcapstone
```

#### Linux (Debian)
Firstly, you should install Golang. You can do this via the following commands in the terminal:

```
sudo apt-get update
sudo apt-get install golang-go
```

For Keystone and Capstone, you may have to build the libraries yourself. Follow [Keystone's](https://github.com/keystone-engine/keystone/blob/master/docs/COMPILE-NIX.md) and [Capstone's](https://github.com/aquynh/capstone/blob/master/COMPILE.TXT) COMPILE instructions for your given system.

### Building the project
Finally, you can build the project by simply running `go build`. You will however need to add your Discord app authentication token to the `config.ini` file - or REBot won't be able to connect to discord.

## License
Specter (Cryptogenic) - [@SpecterDev](https://twitter.com/SpecterDev)

This project is licensed under the WTFPL license - see the [LICENSE.md](LICENSE.md) file for details.
