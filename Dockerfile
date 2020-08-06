FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o intro-detection-info .

FROM gcr.io/distroless/static
COPY --from=builder /build/intro-detection-info /app/
WORKDIR /app
CMD ["./intro-detection-info"]
