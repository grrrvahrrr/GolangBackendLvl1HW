# 1

FROM golang:latest AS build

WORKDIR /myapp

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make build

# 2

FROM scratch

WORKDIR /myapp

COPY --from=build /myapp/app /myapp/app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow

EXPOSE 9000

CMD ["./app"]
