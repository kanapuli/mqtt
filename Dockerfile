FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Update and install basic tools
RUN apt-get update && apt-get install -y \
    neovim \
    git \
    ripgrep \
    fd-find \
    bat \
    jq \
    nodejs \
    npm \
    elixir \
    ruby \
    python3 \
    python3-pip \
    netcat-traditional \
    cargo \
    default-jre \
    tmux

# Install golang
RUN wget https://go.dev/dl/go1.23.2.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.23.2.linux-arm64.tar.gz && \
    rm go1.23.2.linux-arm64.tar.gz
ENV PATH /usr/local/go/bin:$PATH
ENV GOPROXY=https://golangproxy.umh.app,https://proxy.golang.org,direct

# Install Ginkgo for Go testing
RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest && \
    go install github.com/onsi/gomega/...@latest

# Install zig 
RUN wget -O zig.tar.xz https://ziglang.org/builds/zig-linux-aarch64-0.14.0-dev.1860+2e2927735.tar.xz && \
    tar -C /usr/local -xf zig.tar.xz && \
    rm zig.tar.xz
ENV PATH /usr/local/zig-linux-aarch64-0.14.0-dev.1860+2e2927735:$PATH

# Clean up
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set up Go environment
ENV GOPATH /home/vscode/go
ENV PATH $GOPATH/bin:$PATH

# Set up Rust environment
ENV PATH /home/vscode/.cargo/bin:$PATH

# Configure git to skip Git LFS
RUN git config --global lfs.fetchexclude '*'
RUN git config --global core.excludesfile ~/.gitignore

# Set some aliases in zshrc file
RUN echo "alias ll='ls -alF'" >> ~/.zshrc && \
    echo "alias gs='git status'" >> ~/.zshrc && \
    echo "alias gc='git commit'" >> ~/.zshrc && \
    echo "alias gp='git push'" >> ~/.zshrc && \
    echo "alias gpl='git pull'" >> ~/.zshrc && \
    echo "alias gco='git checkout'" >> ~/.zshrc && \
    echo "alias gcb='git checkout -b'" >> ~/.zshrc && \
    echo "alias gcm='git checkout main'" >> ~/.zshrc && \
    echo "alias ga='git add '" >> ~/.zshrc && \
    echo "alias cr='cargo run'" >> ~/.zshrc && \
    echo "alias cb='cargo build'" >> ~/.zshrc && \
    echo "alias ct='cargo test'" >> ~/.zshrc && \
    echo "alias vim='nvim'" >> ~/.zshrc

# Set the default shell to zsh
SHELL ["/bin/zsh", "-c"]

# Set the working directory
WORKDIR /workspace

# Set the default command
CMD ["/bin/zsh"]
