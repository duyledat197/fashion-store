FROM postgres:13-alpine

COPY /src/init-databases.sh /docker-entrypoint-initdb.d/