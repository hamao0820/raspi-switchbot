FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu-24.04

RUN set -x \
    && ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo 'Asia/Tokyo' >/etc/timezone

ARG TARGETARCH \
    GO_VERSION=1.25.0
RUN set -x && curl -fsSL "https://golang.org/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz" | tar -xz -C /usr/local
ENV PATH=$PATH:/usr/local/go/bin
ENV PATH=$PATH:/home/vscode/go/bin
