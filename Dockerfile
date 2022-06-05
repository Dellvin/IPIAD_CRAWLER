FROM golang:1.16-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

WORKDIR cmd
RUN go build -o main .

WORKDIR ../
EXPOSE 80

ENTRYPOINT [ "./main" ]