FROM golang:1.24.2-alpine AS build

RUN apk add --no-cache alpine-sdk make

WORKDIR /app

COPY go.sum go.mod Makefile ./
RUN go mod download

COPY . .

RUN make build-prod
RUN make build-plugins-prod

FROM alpine:3.20.1 AS prod

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=build /app/bin/cookieserver /app/main
COPY --from=build /app/protocols /app/protocols

RUN touch ./cookiefarm.db

EXPOSE ${PORT}

CMD ["./main"]
