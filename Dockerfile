# syntax=docker/dockerfile:1

# Start from the provided image
ARG BASE_IMAGE
FROM ghcr.io/viral32111/${BASE_IMAGE}

# Add the binary from the context directory
COPY --chown=0:0 --chmod=755 ./* /usr/local/bin/apc-ups-exporter

# Launch
ENTRYPOINT [ "apc-ups-exporter" ]
