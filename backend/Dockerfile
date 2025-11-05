FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o beer-festival-backend .

FROM alpine:3.22.2

WORKDIR /app
ENV PORT=1337

COPY --from=build /app/beer-festival-backend .

RUN adduser -D beer-festival
USER beer-festival

EXPOSE 1337
CMD ["./beer-festival-backend"]
