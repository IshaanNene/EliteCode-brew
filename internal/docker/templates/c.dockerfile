FROM gcc:latest

WORKDIR /workspace

# Install basic utilities
RUN apt-get update && apt-get install -y \
    time \
    && rm -rf /var/lib/apt/lists/*

# Set resource limits
USER nobody

CMD ["gcc", "--version"]