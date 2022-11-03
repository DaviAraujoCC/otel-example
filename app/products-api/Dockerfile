FROM golang:1.18 AS base


    WORKDIR /usr/app

    COPY go.mod go.sum ./
    RUN go mod download -x

FROM base AS development
    RUN go install github.com/cortesi/modd/cmd/modd@latest && \
        go install github.com/go-delve/delve/cmd/dlv@latest
    CMD modd -f ./cmd/server/modd.conf


FROM base AS compiler
    COPY . ./
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/app/main ./cmd/server

FROM gcr.io/distroless/static:nonroot AS release
WORKDIR /
COPY --from=compiler /usr/app/main .
USER 65532:65532

ENTRYPOINT ["/main"]