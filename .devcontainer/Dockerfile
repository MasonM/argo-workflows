# TODO: Bump to ubuntu-24.04 once Python issues are fixed: https://github.com/argoproj/argo-workflows/pull/14302#issuecomment-2742883052
FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu-22.04
ARG TARGETARCH
ARG K3D_VERSION=v5.8.3
SHELL ["/bin/bash", "-c"]

# install k3d
RUN wget -q -O - https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=${K3D_VERSION} bash

# install protocol buffer compiler (protoc)
RUN apt-get update && apt-get install -y protobuf-compiler clang-format

# install kubectl
COPY hack/k8s-versions.sh ./
RUN . k8s-versions.sh && \
    wget -O /usr/local/bin/kubectl "https://dl.k8s.io/release/${K8S_VERSIONS[min]}/bin/linux/${TARGETARCH}/kubectl" && \
    chmod +x /usr/local/bin/kubectl
