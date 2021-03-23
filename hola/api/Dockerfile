FROM golang
WORKDIR /
COPY . .
RUN go build api-server.go
EXPOSE 5051
ENTRYPOINT ["go","run","api-server"]