FROM golang:alpine AS builder
ENV CGO_ENABLED=0
COPY . /build/
WORKDIR /build
RUN go build -a -installsuffix docker -ldflags='-w -s' -o /build/bin/abx /build

FROM ghcr.io/acrobox/docker/minimal:latest
ENV ACROBOX_HOME="/machines"
VOLUME /machines
COPY --from=builder /build/bin/abx /usr/local/bin/abx
USER user
ENTRYPOINT ["/usr/local/bin/abx"]
CMD ["help"]

LABEL org.opencontainers.image.source https://github.com/acrobox/abx
