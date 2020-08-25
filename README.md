# Data Infuser GATEWAY

## 개발환경
* Golang 1.14.4
  * gin-gonic (https://github.com/gin-gonic/gin)
  * grpc-go (https://github.com/grpc/grpc-go)
  
## Proto Buffer 공통 모듈 다운로드
```sh
$ git clone git@gitlab.com:promptech1/data-infuser/infuser-protobuf.git
```

## Configuration
* config/config-sample.yaml 참고하여 config/config-dev.yaml 생성

## 개발 환경 실행
* Gateway Server
```sh
go run main.go -logtostderr=true
```