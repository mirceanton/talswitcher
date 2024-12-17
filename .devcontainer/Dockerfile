## ================================================================================================
# Utility versions
## ================================================================================================
ARG TASKFILE_VERSION=v3.38.0@sha256:308c4f5be86bffae3f956cbd7225c4fec69b0e7a89012120b818a10df45d7c59
ARG GOLINT_VERSION=v1.62.2-alpine@sha256:0f3af3929517ed4afa1f1bcba4eae827296017720e08ecd5c68b9f0640bc310d
ARG GORELEASER_VERSION=v2.5.0@sha256:8a70dbaf65cef49cd3ae9b564395fe875fa8136a3d45908f84c6e35ff8925ec9

ARG KUBECTL_VERSION=1.31.4@sha256:64614ef8290f3fb27fed5164b338debeeb79a1e5e26c93eb920770b71abd7c48
ARG TALOSCTL_VERSION=v1.9.0@sha256:938c0536721d023de67c45906dafeff692eb19f65122da6cbf543087222a8b54


## ================================================================================================
# "Build" stage for utilities with docker images already present
## ================================================================================================
FROM ghcr.io/mirceanton/taskfile:${TASKFILE_VERSION} AS taskfile
FROM golangci/golangci-lint:${GOLINT_VERSION} AS golangci-lint
FROM goreleaser/goreleaser:${GORELEASER_VERSION} AS goreleaser

FROM bitnami/kubectl:${KUBECTL_VERSION} AS kubectl
FROM ghcr.io/siderolabs/talosctl:${TALOSCTL_VERSION} AS talosctl

## ================================================================================================
## Main image
## ================================================================================================
FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm@sha256:a417a341a2a8648db7bf8527b86364848362a2c8dc150993c8a4cc2b53b4ec47 AS workspace

# Install additional tools
COPY --from=taskfile /task /usr/local/bin/task
COPY --from=golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint
COPY --from=goreleaser /usr/bin/goreleaser /usr/local/bin/goreleaser
COPY --from=talosctl /talosctl /usr/local/bin/talosctl
COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin/kubectl

# Generate bash completions
RUN golangci-lint completion bash | sudo tee /etc/bash_completion.d/golangci-lint.bash > /dev/null
RUN goreleaser completion bash | sudo tee /etc/bash_completion.d/goreleaser.bash > /dev/null
RUN kubectl completion bash | sudo tee /etc/bash_completion.d/kubectl.bash > /dev/null
RUN talosctl completion bash | sudo tee /etc/bash_completion.d/talosctl.bash > /dev/null

USER vscode
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
