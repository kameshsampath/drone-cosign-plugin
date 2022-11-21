# syntax=docker/dockerfile:1.4

FROM gcr.io/distroless/base
LABEL org.opencontainers.image.source https://github.com/kameshsampath/drone-gcloud-auth
LABEL org.opencontainers.image.authors="Kamesh Sampath<kamesh.sampath@hotmail.com>"

LABEL description="A Drone plugin to sign and verify using projectsigstore cosign"

RUN apk add -U --no-cache ca-certificates

ARG TARGETARCH

COPY plugin_linux_${TARGETARCH}/plugin /bin/plugin

EXPOSE 8080

ENTRYPOINT ["/bin/plugin"]

