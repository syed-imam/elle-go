FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /elle-go .

FROM postgres:16
COPY --from=build /elle-go /usr/local/bin/elle-go
ENV POSTGRES_PASSWORD=elle
