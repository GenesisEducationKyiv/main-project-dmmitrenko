FROM golang:1.20 as builder
WORKDIR /app

COPY . .

RUN go get ./... && go mod download
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/main.go

FROM alpine:latest 

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080

CMD [ "./main" ]
