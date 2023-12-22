FROM golang:1.21 AS BUILDER
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -tags=containers -o /usr/bin/ip_server .
FROM scratch
COPY --from=BUILDER /usr/bin/ip_server /usr/bin/ip_server
ENTRYPOINT [ "/usr/bin/ip_server" ]