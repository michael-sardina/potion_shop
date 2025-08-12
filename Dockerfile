FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o potion-shop

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/potion-shop /app/potion-shop
COPY templates /app/templates
VOLUME ["/app"]
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/potion-shop"]
