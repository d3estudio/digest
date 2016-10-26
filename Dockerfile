# docker build --rm -t d3estudio/digest .
FROM centurylink/ca-certs
ENV GODEBUG=netdns=go

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="Digest" \
      org.label-schema.description="SlackBot that watch channels looking for links and reactions, and generates digests based on those reactions" \
      org.label-schema.url="http://digest.d3.do" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/d3estudio/digest" \
      org.label-schema.vendor="D3 Estudio" \
      org.label-schema.version=$VERSION \
      org.label-schema.schema-version="1.0"

ADD release/linux/amd64/collector     /collector
ADD release/linux/amd64/emoji-manager /emoji-manager
ADD release/linux/amd64/prefetcher    /prefetcher
ADD release/linux/amd64/processor     /processor
