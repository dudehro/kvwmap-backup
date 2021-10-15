FROM golang:bullseye
MAINTAINER Georg Kämmert georg.kaemmert@gdi-service.de

#RUN mkdir -p /backup-interface/backup-config
#ADD ./sources/* /backup-interface/
#RUN apt-get install glibc
#VOLUME /backup-interface/backup-config

WORKDIR /go/src/github.com
#RUN apt update && apt install -y git

RUN git clone https://github.com/dudehro/kvwmap-backup.git
#WORKDIR /go/src/github.com/kvwmap-backup
RUN go build -pkgdir /go/src/github.com/kvwmap-backup

EXPOSE 8082/tcp

CMD ls -al /go/src; go env
