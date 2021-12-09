FROM alpine:latest

RUN apk update
RUN apk add postgresql14-client postgresql13-client mysql-client bash

WORKDIR /scripte
COPY /src/pgdump.sh /scripte
RUN chmod u+x /scripte/pgdump.sh
