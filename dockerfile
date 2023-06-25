FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o bin/app cmd/main.go 

EXPOSE 9871
ENTRYPOINT [ "./bin/app" ]