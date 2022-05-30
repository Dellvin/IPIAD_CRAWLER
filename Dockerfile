FROM golang:1.16-buster AS build

WORKDIR /app

COPY ./go.mod go.sum ./

COPY . .

WORKDIR /app/cmd

RUN go mod vendor
RUN go mod download
RUN go build -o main .

WORKDIR /app

EXPOSE 80

ENTRYPOINT [ "./cmd/main" ]