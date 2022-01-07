#build stage
FROM golang:1.16.13-alpine3.15 AS builder
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# RUN go get -d -v ./...
# RUN go build -o /go/bin/app -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /app/main .
RUN cp -R /app/config .

#final stage
FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
# COPY --from=builder /dist/main /app
COPY --from=builder /dist /app
# RUN chmod +x /app/main
ENTRYPOINT /app/main
LABEL Name=binancechartsvg Version=1.0.0
EXPOSE 8000
