From debian:latest
ADD templates templates
ADD templates templates
ADD script script
RUN mkdir /data
CMD ./fish