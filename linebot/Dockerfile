FROM golang:alpine AS build

WORKDIR /app

COPY ./felice/go.mod .

COPY ./felice/go.sum .

RUN go mod download

COPY ./felice .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp

FROM alpine:3.18.2

COPY --from=build /app/myapp /usr/local/bin/myapp

CMD ["myapp"]
