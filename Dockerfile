FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app/cmd
RUN go test -v main_test.go