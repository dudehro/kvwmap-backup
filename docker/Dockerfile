##
## first stage build
##
FROM golang:alpine as build
MAINTAINER Georg Kämmert georg.kaemmert@gdi-service.de

WORKDIR /app

RUN apk update && apk add git
RUN git clone https://github.com/dudehro/kvwmap-backup.git 

WORKDIR kvwmap-backup
RUN go mod download && go build -o runtime main.go

##
## second stage deploy
##
FROM alpine:latest as deploy
MAINTAINER Georg Kämmert georg.kaemmert@gdi-service.de

RUN apk update && apk add htop

WORKDIR /app
COPY ./src/config.csv .
COPY --from=build /app/kvwmap-backup/runtime .

CMD /app/runtime
