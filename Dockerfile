FROM golang:1.19-alpine as builder

RUN apk add --no-cache git

WORKDIR /workspace
COPY go.* ./
RUN go version
RUN go mod download

COPY . /workspace

# Build
RUN CGO_ENABLED=0 go build -a -o armor ./cmd/armor

FROM alpine
WORKDIR /

COPY --from=builder /workspace/armor .
USER 65532:65532

ENTRYPOINT ["/armor"]
