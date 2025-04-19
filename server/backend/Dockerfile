FROM golang:1.24.2-alpine AS build

RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.sum ./
COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main cmd/api/main.go

RUN for file in $(find ./protocols -name '*.go' ! -name 'protocols.go'); do \
    go build -buildmode=plugin -o "protocols/$(basename ${file%.go}.so)" "$file"; \
done

FROM alpine:3.20.1 AS prod

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=build /app/main /app/main
COPY --from=build /app/protocols /app/protocols

RUN mkdir -p ./internal/database && touch ./internal/database/cookiefarm.db

EXPOSE ${PORT}

CMD ["./main"]
