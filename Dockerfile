From golang:latest

MAINTAINER shamaton

# install cron
RUN apt-get update -y && apt-get install -y cron

# copy resources
COPY ./sh sh
COPY ./src src

# setup cron
RUN echo '* 5 * * * root sh /go/sh/cleaner.sh' >> /etc/crontab
RUN /etc/init.d/cron start

# build application
RUN ./sh/run.sh build

EXPOSE 8080
CMD ["./bin/gyazo"]