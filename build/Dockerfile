FROM       golang:1.15.7
WORKDIR    /stream-server
COPY       go.mod .
COPY       go.sum .
RUN        go mod download
COPY       . .
RUN        go build -o stream-server ./cmd/main.go
ENTRYPOINT ["./stream-server"]