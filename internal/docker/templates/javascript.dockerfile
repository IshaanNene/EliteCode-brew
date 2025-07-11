FROM node:18-slim

WORKDIR /workspace

# Install basic utilities
RUN apt-get update && apt-get install -y \
    time \
    && rm -rf /var/lib/apt/lists/*

# Set resource limits
USER nobody

CMD ["node", "--version"]