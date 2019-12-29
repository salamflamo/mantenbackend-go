FROM golang:latest

ENV MONGO_USER $user
ENV MONGO_PASS $pass
ENV MONGO_HOST $host

WORKDIR /app
COPY . .

RUN go get -v -d ./
RUN go build ./

CMD ["app"]