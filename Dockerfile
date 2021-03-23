FROM golang
WORKDIR /
COPY . . 
RUN go mod download
EXPOSE 5051
CMD ["go","run","api/api-server.go"]