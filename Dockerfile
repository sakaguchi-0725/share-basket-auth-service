FROM golang:1.24-alpine

ENV TZ Asia/Tokyo
ENV GO111MODULE=on

WORKDIR /go/src/app

COPY go.sum go.mod ./

RUN apk add --no-cache git curl unzip bash protobuf

RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go install go.uber.org/mock/mockgen@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/bufbuild/buf/cmd/buf@latest

RUN go mod download

ENV PATH="/go/bin:/usr/local/bin:$PATH"
