FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/windsea ./cmd/server

FROM debian:bookworm-slim
RUN useradd --create-home --uid 10001 windsea
WORKDIR /app
COPY --from=build /out/windsea /app/windsea
USER windsea
ENV HTTP_ADDR=:8080 DATABASE_URL=file:/data/windsea.db?_pragma=foreign_keys(1)
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["/app/windsea"]

