FROM golang:1.15
WORKDIR /project
RUN go get github.com/cheekybits/genny
COPY . .
CMD bash
