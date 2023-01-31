FROM golang:1.19.5-alpine3.16 as builder
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/study github.com/loskutovanl/study


FROM alpine
COPY --from=builder /app/build/study /study
EXPOSE 8080 8080
ENTRYPOINT ["/study"]