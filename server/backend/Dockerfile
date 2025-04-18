FROM golang:1.24.2-alpine AS build

RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.sum ./
COPY go.mod ./
RUN go mod download

COPY . .

# Compilare il binario principale
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main cmd/api/main.go

RUN for file in $(find ./protocols -name '*.go' ! -name 'protocols.go'); do \
    sed -i "1s|^|// |" "$file"; \
    sed -i "2s|^// ||" "$file"; \
done

# Compilare tutti i plugin tranne protocols.go
RUN for file in $(find ./protocols -name '*.go' ! -name 'protocols.go'); do \
    go build -buildmode=plugin -o "protocols/$(basename ${file%.go}.so)" "$file"; \
done

FROM alpine:3.20.1 AS prod

WORKDIR /app

RUN apk add --no-cache libc6-compat

# Copia il binario principale e i plugin compilati
COPY --from=build /app/main /app/main
COPY --from=build /app/protocols /app/protocols

# Altri setup
RUN mkdir -p ./internal/database && touch ./internal/database/cookiefarm.db

EXPOSE ${PORT}

CMD ["./main"]
