FROM golang:1.25.3 AS build
WORKDIR /go/src
COPY . .

# чтобы не было зависимостей от Си-библиотек в конечном бинарнике,
# потому что их не будет в контейнере с альпайн
ENV CGO_ENABLED=0 

RUN go build -o reviewer-assignment-server ./server

FROM scratch
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/reviewer-assignment-server ./
EXPOSE 8080/tcp
ENTRYPOINT ["./reviewer-assignment-server"]