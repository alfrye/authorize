FROM ubuntu:22.04

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID


RUN apt update && apt install -y \
apt-utils \
build-essential \
git \
jq \
net-tools \
openssl \
vim \  
wget \
curl \
sudo \
&& rm -rf /var/lib/apt/lists/*

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME



# Install Go 
ARG GOLANG_VERSION="1.24.0"

RUN curl --connect-timeout 15 --retry-delay 5 -O https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
 tar -C /usr/local -xf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
 rm go${GOLANG_VERSION}.linux-amd64.tar.gz



#Install Go tooling
ARG GO_PLS_VERSION="v0.18.0"
ARG GO_DELV_DLV_VERSION="v1.24.0"
ARG GO_FUMPT="v0.7.0"
ARG GOIMPORTS_VERSION="v0.30.0"
ENV GOPATH="/go"
ENV PATH="${PATH}:/usr/local/go/bin:/go/bin"
ENV GOCACHE=/home/vscode/.cache/go-build

RUN mkdir -p $GOPATH $GOCACHE
RUN chown -R vscode:vscode $GOPATH $GOCACHE

RUN go install golang.org/x/tools/gopls@${GO_PLS_VERSION} \
 && go install github.com/go-delve/delve/cmd/dlv@${GO_DELV_DLV_VERSION} \
 && go install mvdan.cc/gofumpt@${GO_FUMPT} \
 && go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION}

 # Install K3d 
 RUN curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=v5.0.0 bash


# Install Tilt

RUN curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | sudo bash

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
USER vscode

