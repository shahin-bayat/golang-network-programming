FROM golang:1.24.5

WORKDIR /app

RUN apt-get update && apt-get install -y iproute2

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "."]
