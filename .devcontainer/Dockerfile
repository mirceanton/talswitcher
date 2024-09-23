## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=1.31.1@sha256:da4a9868e20d941636087cb8624a4bb441f5249d69e8f3d27e53c7d4d280a5f3
ARG TALOSCTL_VERSION=v1.8.0@sha256:08b8260e162391b7714466f36105902beb3dcac3d4fa1408b4f4838ddbb5b8ea
ARG TALHELPER_VERSION=v3.0.5@sha256:1d5ea10b83e5bce0a32907ada9866267abf1854215000ead86538ffe779c1357
ARG TASKFILE_VERSION=v3.38.0@sha256:308c4f5be86bffae3f956cbd7225c4fec69b0e7a89012120b818a10df45d7c59
ARG GOLINT_VERSION=v1.61.0-alpine@sha256:61e2d68adc792393fcb600340fe5c28059638d813869d5b4c9502392a2fb4c96
ARG GORELEASER_VERSION=v2.3.2@sha256:d62b4a18dfe3af7bd4da9e5954b496548ef04e73ae8f98cd75ba63a9ed4d73e5


## ================================================================================================
# "Build" stage for utilities with docker images already present
## ================================================================================================
FROM bitnami/kubectl:${KUBECTL_VERSION} AS kubectl
FROM ghcr.io/siderolabs/talosctl:${TALOSCTL_VERSION} AS talosctl
FROM ghcr.io/budimanjojo/talhelper:${TALHELPER_VERSION} AS talhelper
FROM ghcr.io/mirceanton/taskfile:${TASKFILE_VERSION} AS taskfile
FROM golangci/golangci-lint:${GOLINT_VERSION} AS golangci-lint
FROM goreleaser/goreleaser:${GORELEASER_VERSION} AS goreleaser


## ================================================================================================
## Main image
## ================================================================================================
FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm@sha256:089dc0f3ae6e3a07b0896c28c0bde042fbb1e31fe87908ca55f32e37f56ff0c6 AS workspace

COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl
RUN kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl.bash > /dev/null

COPY --from=talosctl /talosctl /usr/local/bin/talosctl
RUN talosctl completion bash | sudo tee /etc/bash_completion.d/talosctl.bash > /dev/null

COPY --from=talhelper /usr/local/bin/talhelper /usr/local/bin/talhelper
RUN talhelper completion bash | sudo tee /etc/bash_completion.d/talhelper.bash > /dev/null

COPY --from=golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
RUN golangci-lint completion bash | sudo tee /etc/bash_completion.d/golangci-lint.bash > /dev/null

COPY --from=goreleaser /usr/bin/goreleaser /usr/local/bin/goreleaser
RUN goreleaser completion bash | sudo tee /etc/bash_completion.d/goreleaser.bash > /dev/null

COPY --from=taskfile /task /usr/local/bin/task

USER vscode
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
