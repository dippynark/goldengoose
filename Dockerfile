FROM golang:1.16 as builder

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o goldengoose main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /go/goldengoose .
USER nonroot:nonroot

ENTRYPOINT ["/goldengoose"]
