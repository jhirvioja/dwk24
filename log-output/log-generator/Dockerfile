FROM golang:1.22 AS build-stage
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /log-generator

FROM gcr.io/distroless/base-debian11 AS release-stage
WORKDIR /
COPY --from=build-stage /log-generator /log-generator
USER nonroot:nonroot

ENTRYPOINT ["/log-generator"]