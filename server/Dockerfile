FROM golang:1.23.3-alpine3.20 AS build

WORKDIR /app

ENV CGO_ENABLED=1

COPY . .

RUN apk update && \
    apk add --no-cache gcc musl-dev sqlite-dev

RUN go mod tidy

RUN go build cmd/server/main.go

FROM golang:1.23.3-alpine3.20

WORKDIR /app

COPY --from=build /app/main .
COPY db.sqlite3 .

CMD [ "./main" ]
