FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . /app
RUN go build -o /service ./services/catalog

FROM alpine:3.18
COPY --from=build /service /service
EXPOSE 8082
ENTRYPOINT ["/service"]
