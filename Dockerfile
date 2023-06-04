#docker build -t yusongwang1991/fish-accounting:v1 .
From debian:latest
ADD fish fish
ADD templates templates
ADD script script
RUN mkdir /data
CMD ./fish