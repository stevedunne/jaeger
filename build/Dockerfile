FROM golang:alpine AS builder

ARG ldFlags 

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the Items
#RUN go build -o main .
RUN go build -o /dist/jaeger-agent -ldflags "$ldFlags" ./cmd/agent/main.go
#RUN go build -o /dist/jaeger-collector -ldflags "$ldFlags" ./cmd/collector/main.go
#RUN go build -o /dist/jaeger-ingester -ldflags "$ldFlags" ./cmd/ingester/main.go

# Move to /dist directory as the place for resulting binary folder
#WORKDIR /dist

# Copy binary from build to main folder
#RUN cp /build/cmd/agent/jaeger-agent .

# Build a small image
FROM scratch

COPY --from=builder /dist/jaeger-* /

# Export necessary port
EXPOSE 6831
EXPOSE 14271

# Command to run when starting the container
CMD ["/jaeger-agent", "--reporter.grpc.host-port", "DTINWR2CSVC01.aus.amer.dell.com:14250,DTINWR2CSVC04.aus.amer.dell.com:14250,DTINWWEB01.aus.amer.dell.com:14250" ]
