FROM alpine:3.6

RUN mkdir -p /opt
ADD ./e2e-harness /opt/e2e-harness

RUN mkdir -p /opt/resources
ADD resources/templates/ /opt/resources/templates

RUN \
    apk -Uuv add --update --no-cache \
      build-base=0.5-r0 \
      git=2.13.5-r0 \
      jq=1.5-r3 \
      less=487-r0 \
      libffi-dev=3.2.1-r3 \
      openssh-client=7.5_p1-r1 \
      openssl=1.0.2k-r0

# TODO: add kubectl and helm binaries

WORKDIR /opt

ENTRYPOINT ["/opt/e2e-harness"]
