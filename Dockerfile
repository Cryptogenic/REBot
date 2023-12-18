FROM ubuntu:22.04

# Install dependencies
RUN apt update && apt install -y git golang-go cmake python3 python3-dev python3-pip libcapstone4 libcapstone-dev

# Install Keystone from source
WORKDIR /opt/keystone
RUN git clone https://github.com/keystone-engine/keystone.git /opt/keystone
RUN mkdir /opt/keystone/build
RUN cd /opt/keystone/build
RUN ../make-share.sh
RUN cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED_LIBS=ON -DLLVM_TARGETS_TO_BUILD="AArch64;X86" -G "Unix Makefiles" ..
RUN make -j8
RUN make install
RUN ldconfig

# REBot setup
WORKDIR /opt/rebot
RUN git clone https://github.com/Cryptogenic/REBot.git /opt/rebot
RUN cd /opt/rebot

## Replace 13371337 with your Discord user id for developer commands - TODO: Use environment variables or args instead
RUN sed -i 's/165177089035599873/13371337/g' main.go
RUN sed -i 's/gapstone.New(arch, uint(mode))/gapstone.New(arch, int(mode))/g' /opt/rebot/commands-asm.go
RUN go mod init github.com/Cryptogenic/rebot
RUN go get github.com/bnagy/gapstone
RUN go get github.com/bwmarrin/discordgo
RUN go get github.com/go-ini/ini
RUN go get github.com/keystone-engine/keystone/bindings/go/keystone
RUN go build

# Add REBot config with Discord Bot token
ADD config.ini /opt/rebot/config.ini

# Run
ENTRYPOINT ["/opt/rebot/rebot"]
