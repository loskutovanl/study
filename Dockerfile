FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/github.com/loskutovanl/study/
WORKDIR /go/src/github.com/loskutovanl/study/
RUN go mod download
COPY . /go/src/github.com/loskutovanl/study/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/study github.com/loskutovanl/study/

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/loskutovanl/study/build/study /usr/bin/study
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/study"]