## ================================================================================================
# Utility versions
## ================================================================================================
ARG TASKFILE_VERSION=v3.38.0@sha256:308c4f5be86bffae3f956cbd7225c4fec69b0e7a89012120b818a10df45d7c59
ARG GOLINT_VERSION=v1.63.4-alpine@sha256:e1768ec667a10f7712d38abc484d18660d8c32caf3a252b8531ea77ff09de575
ARG GORELEASER_VERSION=v2.7.0@sha256:41247b711fc423eddcd8395b998b9de5c7d169bfe99e6a474ba35e87896c1343

ARG KUBECTL_VERSION=1.32.1@sha256:a7d3b17920b36a4e146b0e3fbfc73db45cceb2eeb94b7af6be672d34f499c86f
ARG TALOSCTL_VERSION=v1.9.3@sha256:d1d324131e46d915ab16a24aba8e5aabd3361c0749f78fb2cf16703eafae8a2f


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
