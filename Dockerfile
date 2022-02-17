# FROM archlinux

# WORKDIR /myapp

# COPY ./app ./

# CMD ["./app"]

# 1

FROM golang:latest AS build

WORKDIR /myapp

COPY go.mod .
COPY go.sum .
# COPY ./config/config.env ./config/config.env
# COPY ./logs/access.log ./logs/access.log
# COPY ./logs/error.log ./logs/error.log
# COPY ./app ./
COPY /home/deus/Documents/testData/covid_19_data.csv ./
RUN go mod download

RUN make build

# 2

FROM scratch

WORKDIR /myapp

COPY --from=build /myapp/app /myapp/app
COPY --from=build /myapp/logs/error.log /myapp/logs/error.log
COPY --from=build /myapp/logs/acess.log /myapp/logs/access.log
COPY --from=build /myapp/config/config.env /myapp/config/config.env
COPY --from=build /myapp/covid_19_data.csv /myapp/covid_19_data.csv
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow
