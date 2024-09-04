FROM scratch
USER 8675:8675
COPY talswitcher /
ENTRYPOINT ["/talswitcher"]