FROM alpine:3.7

RUN adduser -D -u 1001 e2e-harness
ENV HOME=/home/e2e-harness
ENV WORKDIR=/workdir

RUN mkdir -p ${WORKDIR}

RUN chown -R e2e-harness:e2e-harness ${WORKDIR} ${HOME}

RUN apk -Uuv add --update --no-cache \
      bash=4.4.19-r1 \
      build-base=0.5-r0 \
      git=2.15.2-r0 \
      jq=1.5-r5 \
      less=520-r0 \
      libffi-dev=3.2.1-r4 \
      openssh-client=7.5_p1-r8 \
      openssl=1.0.2o-r0 \
      sudo=1.8.21_p2-r1 \
      iptables=1.6.1-r1

ENV KUBECTL_VERSION=v1.10.1
ENV HELM_VERSION=v2.8.2
ENV SHIPYARD_VERSION=v0.1.0

RUN wget https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
  chmod a+x ./kubectl && \
  mv ./kubectl /usr/local/bin

RUN wget https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz && \
  tar zxfv helm-${HELM_VERSION}-linux-amd64.tar.gz && \
  chmod a+x ./linux-amd64/helm && \
  mv ./linux-amd64/helm /usr/local/bin && \
  rm -rf ./linux-amd64 helm-${HELM_VERSION}-linux-amd64.tar.gz

RUN wget https://github.com/giantswarm/shipyard/releases/download/${SHIPYARD_VERSION}/shipyard && \
  chmod a+x ./shipyard && \
  mv ./shipyard /usr/local/bin

RUN echo '%e2e-harness ALL=(root) NOPASSWD: /sbin/iptables' | tee -a /etc/sudoers.d/e2e-harness

USER e2e-harness

RUN mkdir -p ${HOME}/.helm/plugins/ && \
  git clone https://github.com/app-registry/appr-helm-plugin.git ${HOME}/.helm/plugins/registry && \
  helm registry --help

WORKDIR ${WORKDIR}

ENTRYPOINT ["/bin/true"]
