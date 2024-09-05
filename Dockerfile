FROM alpine:3.20.2
USER 8675:8675
COPY talswitcher /
ENTRYPOINT ["/talswitcher"]
