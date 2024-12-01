FROM golang:1.23.2-alpine3.20

WORKDIR /go/src/app

COPY . .

EXPOSE 8081

RUN go build -o main main.go

CMD [ "./main" ]