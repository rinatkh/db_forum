FROM golang:latest AS build

WORKDIR /app

COPY . ./
RUN go build cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get -y update && apt-get install -y postgresql-12 && rm -rf /var/lib/apt/lists/*
USER postgres

RUN /etc/init.d/postgresql start && \
  psql --command "CREATE USER root WITH SUPERUSER PASSWORD 'password';" && \
  createdb -O root docker && \
  /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /cmd

COPY ./db/db.sql ./db.sql
COPY --from=build /app/main .

EXPOSE 5000
ENV PGPASSWORD password
CMD service postgresql start && psql -h localhost -d docker -U root -p 5432 -a -q -f ./db.sql && ./main
