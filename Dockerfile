ARG GOLANG_VERSION=1.18
FROM golang:${GOLANG_VERSION} as builder

WORKDIR /aws-quota-provider

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN cd /aws-quota-provider && go mod download

# Build
RUN go build -a -o quota ./cmd/quota

# Use distroless as minimal base image to package the manager binary
FROM gcr.io/distroless/base:latest-amd64
WORKDIR /

LABEL org.opencontainers.image.source https://github.com/ydataai/aws-quota-provider

COPY --from=builder /aws-quota-provider/quota .

ENTRYPOINT ["/quota"]
