# builder image
FROM golang:alpine as builder
COPY  . /build/
WORKDIR /build
RUN apk update && apk add git
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -o clientdb .

# generate clean, final image for end users
FROM scratch
WORKDIR /app
COPY --from=builder /build/ .

EXPOSE 9000

# executable
ENTRYPOINT [ "./clientdb" ]