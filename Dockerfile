FROM golang:1.15-buster
RUN apt-get install git
WORKDIR /app

COPY . .

RUN go build helloworld.go

EXPOSE 51101

CMD ["./app/helloworld"]
