FROM alpine:3.6

RUN mkdir -p /opt
ADD ./e2e-harness /opt/e2e-harness

RUN mkdir -p /opt/resources
ADD resources/templates/ /opt/resources/templates

RUN \
    apk -Uuv add --update \
      build-base \
      git \
      jq \
      less \
      libffi-dev \
      openssh-client \
      openssl && \
    rm /var/cache/apk/*

# TODO: add kubectl and helm binaries

WORKDIR /opt

ENTRYPOINT ["/opt/e2e-harness"]
