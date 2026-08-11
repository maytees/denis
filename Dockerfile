FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /denis .

FROM alpine:3.20
WORKDIR /app

COPY --from=build /denis /app/denis

# DENIS looks for ./config in the working dir,
# so /app/config is where the TOML files get mounted
EXPOSE 53/udp
ENTRYPOINT ["/app/denis"]
