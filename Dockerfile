## Build
FROM golang:1.20-alpine AS build

# установка рабочей директории и копирование в нее объектов
WORKDIR .
RUN git submodule update --init --recursive

# загружаем зависимости
COPY go.mod .
COPY go.sum .

# билд
RUN CGO_ENABLED=0 GOOS=linux go build -o build/ags ./cmd/ags

## Deploy
FROM alpine

# копируем билд в alpine контейнер с заметно меньшим размером
COPY --from=build /go/src/srvgitlab.imatic.pro/imatic/ags/build/ags /usr/bin/ags

ENTRYPOINT ["/usr/bin/ags", "start", "--cfgname", "app/config.yml"]
