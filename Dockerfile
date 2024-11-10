FROM golang:1.23.2-alpine3.20

WORKDIR /go/src/app

COPY . .

EXPOSE 8081

RUN go build -o main main.go

CMD [ "./main" ]

# FROM golang:1.23.2-alpine3.20 AS builder

# ENV GOPATH /go
# ENV GOROOT /usr/local/go
# ENV PACKAGE github.com/alaimcosta/go-api-task
# ENV BUILD_DIR ${GOPATH}/src/${PACKAGE}

# COPY . ${BUILD_DIR}
# WORKDIR ${BUILD_DIR}

# COPY go.mod go.sum ./

# RUN go build -o app .

# RUN cp -v app /usr/bin/app

# FROM alpine:latest
# COPY --from=builder /usr/bin/app /usr/bin/app
# COPY /entrypoint.sh /entrypoint.sh

# RUN chmod 755 /entrypoint.sh

# EXPOSE 8081
# ENTRYPOINT [ "/entrypoint.sh" ]