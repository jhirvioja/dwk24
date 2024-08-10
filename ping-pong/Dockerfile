FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /ping-pong

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /ping-pong /ping-pong

USER nonroot:nonroot

ENTRYPOINT ["/ping-pong"]