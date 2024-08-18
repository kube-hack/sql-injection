FROM golang

COPY main.go go.mod go.sum ./

RUN go mod download && \
    go build -o main main.go

ENTRYPOINT [ "./main" ]