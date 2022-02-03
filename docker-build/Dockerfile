##
## first stage build
##
FROM golang:stretch as build
MAINTAINER Georg Kämmert georg.kaemmert@gdi-service.de

WORKDIR /app
COPY ./go.mod /app

RUN apt update && apt install git
RUN git clone https://github.com/dudehro/kvwmap-backup.git 

WORKDIR kvwmap-backup
#RUN go mod download && go build -o /app/bin
RUN go mod download
#RUN ls -alh

CMD git pull; go build -o /app/bin
