FROM scratch

COPY bin /

WORKDIR /

ENTRYPOINT ["/helper-reset-password"]
