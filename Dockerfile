# =================================================================================================
# BUILDER STAGE
# =================================================================================================
FROM golang:1.23.0-alpine@sha256:44a2d64f00857d544048dd31d8e1fbd885bb90306819f4313d7bc85b87ca04b0 AS builder

ARG PKG=github.com/mirceanton/talswitcher
ARG VERSION=dev

WORKDIR /build
COPY . .

RUN go build -ldflags "-s -w -X github.com/mirceanton/talswitcher/cmd.version=${VERSION}" -o talswitcher


# =================================================================================================
# PRODUCTION STAGE
# =================================================================================================
FROM scratch
USER 8675:8675
COPY --from=builder --chmod=555 /build/talswitcher /talswitcher
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/talswitcher"]