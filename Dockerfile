from golang:1.22

WORKDIR /app
COPY . /app
RUN go build -o gt7-social
EXPOSE 1323

ENTRYPOINT ["/app/gt7-social"]