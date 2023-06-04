From mcr.microsoft.com/devcontainers/go:0-1.18
ADD templates templates
ADD fish fish
RUN mkdir /data
CMD ./fish