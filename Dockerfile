FROM golang:1.26 AS build

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/aegis-mcp ./cmd/aegis-mcp

FROM gcr.io/distroless/static-debian12

COPY --from=build /out/aegis-mcp /usr/local/bin/aegis-mcp
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/aegis-mcp"]
