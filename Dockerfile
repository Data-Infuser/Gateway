FROM amd64/golang:1.14

WORKDIR /go/src/gateway

COPY go.sum go.mod ./

RUN go mod download

COPY . .

# 배포 환경 설정
ARG GATEWAY_ENV=dev

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    GATEWAY_ENV=$GATEWAY_ENV

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

