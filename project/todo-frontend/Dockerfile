FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-frontend

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /todo-frontend /todo-frontend

USER nonroot:nonroot

ENTRYPOINT ["/todo-frontend"]