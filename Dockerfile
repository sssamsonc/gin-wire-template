# Compile stage
FROM golang:1.21-alpine3.18 AS build-env
ENV CGO_ENABLED 0
ADD . /go/src/gin-wire-template
WORKDIR /go/src/gin-wire-template
RUN go build -o gin-wire-template main.go wire_gen.go

# Final stage
FROM alpine:3.18
EXPOSE 8080
RUN addgroup -S nonroot && adduser -u 65530 -S nonroot -G nonroot

# install the timezone
RUN apk add -U tzdata

#golang GIN framework
ENV GIN_MODE=release

WORKDIR /
COPY .env /
COPY --from=build-env /go/src/gin-wire-template/gin-wire-template /
USER 65530
CMD ["/gin-wire-template"]