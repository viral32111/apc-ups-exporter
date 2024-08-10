# syntax=docker/dockerfile:1

# Start from the provided image
ARG BASE_IMAGE=alpine:3.20
FROM ghcr.io/viral32111/${BASE_IMAGE}

# Add the binary from the context directory
COPY --chown=0:0 --chmod=755 ./* /usr/local/bin/apc-ups-exporter

# Switch to the regular user
USER ${USER_ID}:${USER_ID}

# Publish the default metrics port
EXPOSE 5000/tcp

# Launch the exporter
ENTRYPOINT [ "apc-ups-exporter" ]
