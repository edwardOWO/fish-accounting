#docker build -t yusongwang1991/fish-accounting:20230613v1 .
#docker build -t yusongwang1991/fish-accounting:20230616v1 .
From debian:latest
ADD fish fish
ADD templates templates
ADD script script
RUN mkdir /data
CMD ./fish