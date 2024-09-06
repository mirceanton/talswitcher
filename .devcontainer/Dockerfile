## ================================================================================================
# Utility versions
## ================================================================================================
ARG KUBECTL_VERSION=1.31.0@sha256:7779e585c588adde5c8feafd584c1ffc76ffbe8475236aa14c13b8afe18d9af6
ARG TALOSCTL_VERSION=v1.7.6@sha256:c0d0c85fc25424ec8e28d7b98db5b750ab47705af44222d3ff3afc556fac52d5
ARG TALHELPER_VERSION=v3.0.5@sha256:1d5ea10b83e5bce0a32907ada9866267abf1854215000ead86538ffe779c1357
ARG TASKFILE_VERSION=v3.38.0@sha256:308c4f5be86bffae3f956cbd7225c4fec69b0e7a89012120b818a10df45d7c59
ARG GOLINT_VERSION=v1.60.3-alpine@sha256:8326b535f6aff1995059a0025449e8fe150b3c8e48d7c310855302c638cb2ccf
ARG GORELEASER_VERSION=v2.2.0@sha256:847eba313710996f4eeeef1ac21a4fd5ea7527e3a500633ae78f40d9baefddac








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
FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm@sha256:a3403234ec01e5383f63414aabfd3e2703a4c1e44b260354495c6f1e8cdc7594 AS workspace

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
