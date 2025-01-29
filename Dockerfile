FROM golang:1.23.5-alpine AS build

LABEL maintainer="info@reinkrul.nl"

ENV GO111MODULE=on
ENV GOPATH=/

ARG TARGETARCH
ARG TARGETOS

COPY . .

RUN mkdir /app
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/bin .

# ## Deploy
FROM alpine:latest

COPY --from=build /app/bin /app/bin

WORKDIR /app

# Run the app binary when we run the container
ENTRYPOINT ["/app/bin"]
