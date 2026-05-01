# syntax=docker/dockerfile:1
FROM node:24-slim AS client-builder
WORKDIR /build/client
COPY client/package.json client/package-lock.json* ./
RUN npm ci
COPY client/ ./
RUN npm run build

FROM golang:1.26-alpine AS server-builder
WORKDIR /build/server/src
COPY server/src/go.mod server/src/go.sum ./
RUN go mod download
COPY server/src/ ./
RUN CGO_ENABLED=0 go build -o /out/<BINARY_NAME> ./cmd/<BINARY_NAME>

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=server-builder /out/<BINARY_NAME> /usr/local/bin/<BINARY_NAME>
COPY --from=client-builder /build/client/dist /var/lib/<BINARY_NAME>/static
USER nonroot:nonroot
EXPOSE <HTTP_PORT>
ENTRYPOINT ["/usr/local/bin/<BINARY_NAME>"]
CMD ["--config", "/etc/pgedge/<BINARY_NAME>.yaml"]
