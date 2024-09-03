FROM ubuntu:latest
RUN apt-get update -y && apt install iputils-ping -y && apt install net-tools -y && apt install curl -y&& apt install telnet -y
LABEL authors="wenbin"
COPY demo /app/demo
WORKDIR /app

ENTRYPOINT ["/app/demo"]