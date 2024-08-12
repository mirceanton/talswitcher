## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=1.30.3
ARG K9S_VERSION=v0.32.4
ARG TALOSCTL_VERSION=v1.7.6
ARG HELM_VERSION=v3.15.3
ARG TASKFILE_VERSION=v3.38.0
ARG TALHELPER_VERSION=v3.0.5

## ================================================================================================
# "Build" stage for utilities with docker images already present
## ================================================================================================
FROM bitnami/kubectl:${KUBECTL_VERSION} AS kubectl
FROM derailed/k9s:${K9S_VERSION} AS k9s
FROM ghcr.io/siderolabs/talosctl:${TALOSCTL_VERSION} AS talosctl
FROM ghcr.io/budimanjojo/talhelper:${TALHELPER_VERSION} AS talhelper


## ================================================================================================
# Build stages for other utilities
## ================================================================================================
FROM alpine AS taskfile
ARG TASKFILE_VERSION
RUN wget https://github.com/go-task/task/releases/download/${TASKFILE_VERSION}/task_linux_amd64.tar.gz && tar xvf task_linux_amd64.tar.gz && mv task /bin/task
RUN wget https://raw.githubusercontent.com/go-task/task/${TASKFILE_VERSION}/completion/bash/task.bash -O /task_completion.bash

FROM alpine AS helm
ARG HELM_VERSION
RUN wget https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz && tar xvf helm-${HELM_VERSION}-linux-amd64.tar.gz && mv linux-amd64/helm /bin/helm


## ================================================================================================
## Main image
## ================================================================================================
FROM mcr.microsoft.com/devcontainers/go:1.22-bookworm AS workspace

ENV EDITOR=vim

RUN DEBIAN_FRONTEND=noninteractive \
    apt-get update && apt-get upgrade -y && \
    apt-get install -y \
    sudo \
    git \
    bash-completion \
    vim \
    curl \
    wget \
    unzip \
    htop \
    net-tools \
    iputils-ping \
    dnsutils && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
RUN mkdir -p /etc/bash_completion.d

# Enable passwordless sudo :kek:
RUN echo 'vscode ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Install k9s with no bash completion
COPY --from=k9s /bin/k9s /usr/local/bin/k9s

# Install talosctl and set up bash completion
COPY --from=talosctl /talosctl /usr/local/bin/talosctl
RUN talosctl completion bash | sudo tee /etc/bash_completion.d/talosctl.bash > /dev/null

# Install talhelper and set up bash completion
COPY --from=talhelper /usr/local/bin/talhelper /usr/local/bin/talhelper
RUN talhelper completion bash | sudo tee /etc/bash_completion.d/talhelper.bash > /dev/null

# Install taskfile and set up bash completion
COPY --from=taskfile /bin/task /usr/local/bin/task
COPY --from=taskfile /task_completion.bash /etc/bash_completion.d/task.bash

# Install kubectl and set up bash completion
COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl
RUN kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl.bash > /dev/null

# Install helm and set up bash completion
COPY --from=helm /bin/helm /usr/local/bin/helm
RUN helm completion bash | sudo tee /etc/bash_completion.d/helm.bash > /dev/null

USER vscode
WORKDIR /workspace
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
