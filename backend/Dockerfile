FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o beer-festival-backend .

FROM alpine:latest

WORKDIR /app

ENV PORT=1337

COPY --from=build /app/beer-festival-backend .

COPY --from=build /app/festivals.json .

EXPOSE 1337

CMD ["./beer-festival-backend"]
