FROM golang:1.22

WORKDIR /app

COPY main.go go.mod go.sum ./

RUN go mod download

RUN go build -o main main.go

EXPOSE 8080

ENTRYPOINT [ "./main" ]