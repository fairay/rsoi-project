# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p logs
RUN mkdir -p temp

RUN go build -o ./app

ENTRYPOINT [ "./app" ]
