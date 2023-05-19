FROM golang:1.18-alpine AS build

RUN apk add build-base
WORKDIR /app

COPY cmd internal go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /ctf-kartochki2-backend cmd/main.go

FROM golang:1.18-alpine

ENV APP_HOME /app/ctf-kartochki2-backend
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

RUN mkdir ./cmd ./internal ../db/
COPY cmd ./cmd
COPY internal ./internal
COPY favicon.ico go.mod go.sum ./
COPY sqlite.db/ ../db/

COPY --from=build /ctf-kartochki2-backend ./

EXPOSE 80

ENTRYPOINT ["./ctf-kartochki2-backend"]