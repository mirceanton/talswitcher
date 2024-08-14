# =================================================================================================
# BUILDER STAGE
# =================================================================================================
FROM golang:1.23.0-alpine@sha256:d0b31558e6b3e4cc59f6011d79905835108c919143ebecc58f35965bf79948f4 AS builder

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