# docker build --rm -t d3estudio/digest .

FROM centurylink/ca-certs
EXPOSE 8000

ENV GODEBUG=netdns=go

WORKDIR /app

ADD release/linux/amd64/collector     /app/collector
ADD release/linux/amd64/emoji-manager /app/emoji-manager
ADD release/linux/amd64/prefetcher    /app/prefetcher
ADD release/linux/amd64/processor     /app/processor
