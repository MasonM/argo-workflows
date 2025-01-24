FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu-24.04

ARG INSTALL_K3S_VERSION=v1.31.2+k3s1
ARG KIT_VERSION=v0.1.21

# install protocol buffer compiler (protoc)
RUN apt-get update && apt-get install -y protobuf-compiler

# install k3s
COPY hack/k3s.sh /etc/init.d/k3s
RUN /etc/init.d/k3s install

# install kit
RUN curl -q https://raw.githubusercontent.com/kitproj/kit/main/install.sh | tag=${KIT_VERSION} sh