FROM ubuntu:22.04

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG NVM_VERSION=v0.39.7
ARG GOLANG_VERSION="1.22.2"
#ARG NODE_VERSION="22.0.0"

ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION="22.0.0"
ENV GOPATH="/go"
ENV PATH="${PATH}:/usr/local/go/bin:$GOPATH/bin:/home/${USERNAME}/.local/share/lua-lang-server/bin"


RUN mkdir -p $NVM_DIR

RUN apt-get update && apt-get upgrade -y
RUN apt-get update && apt-get install --no-install-recommends -y \
apt-utils \
  ca-certificates \
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

USER vscode
# install lua
#WORKDIR /tmp
RUN mkdir -p ~/.local/share/lua-lang-server  
RUN sudo curl --connect-timeout 15 --retry-delay 5 -fsSL -O https://github.com/LuaLS/lua-language-server/releases/download/3.15.0/lua-language-server-3.15.0-linux-x64.tar.gz \ 
&&  tar -xvzf lua-language-server-3.15.0-linux-x64.tar.gz -C ~/.local/share/lua-lang-server && \
sudo rm lua-language-server-3.15.0-linux-x64.tar.gz

#&& rm -rf /tmp
# #WORKDIR /
# Install nodejs 
USER root
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash - && \
sudo apt-get install -y nodejs 

#Install nvm 
# RUN curl -fsSL "https://raw.githubusercontent.com/nvm-sh/nvm/${NVM_VERSION}/install.sh" | bash 
# # 
#
 #RUN /bin/bash -c "source $NVM_DIR/nvm.sh && nvm install $NODE_VERSION && nvm use --delete-prefix $NODE_VERSION" 


#Install yarn 
RUN  sudo npm install  --global yarn

 # Install Go 
 ARG GOLANG_VERSION="1.24.0"
    
 RUN  curl --connect-timeout 15 --retry-delay 5 -O https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
  tar -C /usr/local -xf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
  rm go${GOLANG_VERSION}.linux-amd64.tar.gz


#Install Go tooling
ARG BUF_VERSION="v1.9.0"
ARG GO_PLS_VERSION="v0.18.0"
ARG GO_DELV_DLV_VERSION="v1.24.0"
ARG GO_FUMPT="v0.7.0"
ARG GOIMPORTS_VERSION="v0.30.0"
ARG GOLANG_CI_LINT_VERSION="v1.52.2"
ARG GRPC_GATEWAY_VERSION="v2.14.0"
ARG PROTOC_GEN_GO_GRPC_VERSION="v1.2.0"
ARG PROTOC_GEN_GO_VERSION="v1.28.1"

ENV GOPATH="/go"
ENV PATH="${PATH}:/usr/local/go/bin:/go/bin:~/.local/share/lua-lang-server/bin"
ENV GOCACHE=/home/vscode/.cache/go-build
ENV CACHE=/home/vscode/.cache
#USER vscode
RUN  mkdir -p $GOPATH $GOCACHE


RUN go install golang.org/x/tools/gopls@${GO_PLS_VERSION} \
  && go install github.com/go-delve/delve/cmd/dlv@${GO_DELV_DLV_VERSION} \
  && go install github.com/bufbuild/buf/cmd/buf@${BUF_VERSION} \
  && go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANG_CI_LINT_VERSION} \
  && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@${GRPC_GATEWAY_VERSION} \
  && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@${GRPC_GATEWAY_VERSION} \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION} \
  && go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION} \
  && go install mvdan.cc/gofumpt@${GO_FUMPT} \
  && go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION} \
  && go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest \
  && go clean -modcache -cache -x \
  && chown -R vscode:vscode $GOPATH $CACHE

# Install formatters used by null-ls (outside Go ecosystem)
RUN npm install -g @johnnymorganz/stylua-bin && \
    apt-get update && apt-get install -y shfmt && \
    rm -rf /var/lib/apt/lists/*

# Install kubectl
ARG KUBERNETES_VERSION="v1.29.2"
ENV K3S_VERSION="${KUBERNETES_VERSION}-k3s1"
# RUN curl --connect-timeout 15 --retry 5 --retry-delay 5 -sSL -o /usr/local/bin/kubectl https://dl.k8s.io/release/${KUBERNETES_VERSION}/bin/linux/amd64/kubectl && \
#   chmod +x /usr/local/bin/kubectl

# Install K3d 
RUN curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=v5.0.0 bash


# Install Tilt

RUN curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | sudo bash

RUN curl -LO "https://dl.k8s.io/release/$(sudo curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
  && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
USER vscode
ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"
