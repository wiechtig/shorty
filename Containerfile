FROM golang:1.24 AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/shorty

FROM gcr.io/distroless/static-debian12
WORKDIR /app

COPY ./db/migrations ./db/migrations
COPY --from=build /go/bin/app ./bin

EXPOSE 4242
EXPOSE 4343
CMD ["./bin"]
