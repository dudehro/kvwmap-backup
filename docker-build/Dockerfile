##
## first stage build
##
FROM golang:stretch as build
MAINTAINER Georg Kämmert georg.kaemmert@gdi-service.de

WORKDIR /app

RUN apt update && apt install git
RUN git clone https://github.com/dudehro/kvwmap-backup.git 

WORKDIR kvwmap-backup
RUN go mod download && go build -o runtime main.go

CMD git pull; go build
