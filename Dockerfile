FROM golang:latest AS build

MAINTAINER Nikita Lobaev

RUN mkdir /go/src/BMSTU-DB

COPY . /go/src/BMSTU-DB

WORKDIR /go/src/BMSTU-DB/app/main

RUN go build -o bmstu-db .

FROM ubuntu:20.04 AS release

MAINTAINER Nikita Lobaev

RUN apt-get update -y && apt-get install -y locales gnupg2
RUN locale-gen en_US.UTF-8
RUN update-locale LANG=en_US.UTF-8

ENV PGVER 12
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

USER postgres

COPY postgres/. /home

WORKDIR /home

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forums_user WITH SUPERUSER PASSWORD 'forums_user';" &&\
    psql --command "\i '/home/default.sql'" &&\
    createdb -E UTF8 forums_1nf &&\
    psql --command "\i '/home/1nf.sql'" &&\
    createdb -E UTF8 forums_2nf &&\
    psql --command "\i '/home/2nf.sql'" &&\
    createdb -E UTF8 forums_3nf &&\
    psql --command "\i '/home/3nf.sql'" &&\
    /etc/init.d/postgresql stop

RUN echo "listen_addresses='*'\n" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

COPY --from=build /go/src/BMSTU-DB/app/main/bmstu-db /usr/bin/BMSTU-DB/
COPY --from=build /go/src/BMSTU-DB/config.json /usr/bin/BMSTU-DB/

EXPOSE 5432
EXPOSE 5000

WORKDIR /usr/bin/BMSTU-DB

CMD service postgresql start && ls && ./bmstu-db config.json
