FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
ARG VERSION=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.builtAt=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X main.version=${VERSION}" -o intro-detection-info .

FROM gcr.io/distroless/static
COPY --from=builder /build/intro-detection-info /app/
WORKDIR /app
ENTRYPOINT ["./intro-detection-info"]
