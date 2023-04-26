ARG GOLANG_VERSION=1.20
FROM golang:${GOLANG_VERSION} as builder

ARG COMPILE_CMD

WORKDIR /workspace

COPY . .

ENV CGO_ENABLED=0

RUN go mod download

# Build
RUN go build -a -o main "./cmd/${COMPILE_CMD}"

# Use distroless as minimal base image to package the manager binary
FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /

LABEL org.opencontainers.image.source https://github.com/ydataai/aws-adapter

COPY --from=builder /workspace/main .

ENTRYPOINT ["/main"]
