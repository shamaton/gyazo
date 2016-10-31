From golang:latest

MAINTAINER shamaton

COPY ./sh sh
COPY ./src src

RUN ./sh/run.sh build

EXPOSE 8080
CMD ["./bin/gyazo"]