FROM golang:1.15
RUN mkdir /app
COPY . /app/
WORKDIR /app/cmd
RUN go test -v main_test.go
