FROM golang:1.14.6-alpine3.12 as builder
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /build
ADD go.mod .
ADD go.sum .
COPY . .
RUN go mod download
RUN go run ./cmd/app/main.go
# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/study study/
# go run ./cmd/app/main.go

# FROM alpine
# COPY --from=builder /go/src/study/build/study /usr/bin/study
# EXPOSE 8080 8080
# ENTRYPOINT ["/usr/bin/study"]