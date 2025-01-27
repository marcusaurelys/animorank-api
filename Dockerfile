FROM golang:latest

VOLUME /var/run/docker.sock

COPY go.mod go.sum main.go ./ 

RUN apt-get update && apt-get install -y docker.io
RUN go build main.go 

EXPOSE 8080

CMD ["./main"]
