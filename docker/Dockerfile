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

WORKDIR /app
COPY --from=build /app/kvwmap-backup  .
EXPOSE 8082/tcp

CMD /app/runtime
