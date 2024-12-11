FROM ubuntu:20.04
COPY wetravel /app/wetravel
WORKDIR /app
CMD [ "/app/wetravel" ]
