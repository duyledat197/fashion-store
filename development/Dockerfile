#* build:
FROM golang:1.21.5 as build-stage

WORKDIR /app
COPY .. /app
RUN CGO_ENABLED=0 && go build -o app-exe

#* production: 
FROM alpine:3.14 as production-stage
COPY --from=build-stage /app/app-exe /bin
RUN chmod +x /bin/app-exe
