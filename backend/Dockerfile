FROM golang:1.21.3-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download
RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

FROM golang:1.21.3-alpine as builder
RUN apk update && apk add --no-cache nodejs npm
RUN npm install -g @redocly/cli
COPY --from=modules /go/pkg /go/pkg
COPY --from=modules /go/bin/oapi-codegen /go/bin/oapi-codegen
COPY --from=modules /go/bin/dlv /go/bin/dlv
COPY . /app
WORKDIR /app
RUN sh generate-openapi.sh
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /bin/app ./cmd/app

FROM alpine
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/sql /sql
COPY --from=builder /bin/app /app
COPY --from=builder /go/bin/dlv /bin/dlv
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/.env .
COPY --from=builder /tmp /tmp
CMD ["/app"]
