FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o beer-festival-backend .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/beer-festival-backend .

COPY --from=build /app/festivals.json .

EXPOSE 8080

ENV PORT=8080

CMD ["./beer-festival-backend"]
