FROM postgres:13

ENV POSTGRES_USERF=postgres
ENV POSTGRES_DB=notes

COPY ./db/init.sql /docker-entrypoint-initdb.d/
