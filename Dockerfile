# =================================================================================================
# BUILDER STAGE
# =================================================================================================
FROM golang:1.22.6-alpine@sha256:1a478681b671001b7f029f94b5016aed984a23ad99c707f6a0ab6563860ae2f3 AS builder

ARG PKG=github.com/mirceanton/talswitcher
ARG VERSION=dev

WORKDIR /build
COPY . .

RUN go build -o talswitcher


# =================================================================================================
# PRODUCTION STAGE
# =================================================================================================
FROM scratch
USER 8675:8675
COPY --from=builder --chmod=555 /build/talswitcher /talswitcher
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/talswitcher"]