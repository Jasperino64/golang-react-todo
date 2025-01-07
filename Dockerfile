FROM golang:1.23

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN ENV=production
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp .
CMD ["/myapp"]

EXPOSE 5000