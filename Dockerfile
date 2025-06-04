FROM golang:1.24-bullseye
WORKDIR /app
COPY . .
ENV CGO_ENABLED=1
RUN go mod download
RUN go build -o whatsapp-bridge
EXPOSE 8080
CMD ["./whatsapp-bridge"]
