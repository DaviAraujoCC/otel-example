FROM golang:1.18 as base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM gcr.io/distroless/static:nonroot AS release
WORKDIR /
COPY --from=base /app/main .
USER 65532:65532

ENTRYPOINT ["/main"]


