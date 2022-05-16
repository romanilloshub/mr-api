FROM golang:alpine AS build-env
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o app

FROM alpine

WORKDIR /

COPY --from=build-env /app/app /app
EXPOSE 8080

ENTRYPOINT [ "/app" ]