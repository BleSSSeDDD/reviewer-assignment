FROM golang:1.25.3 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .
COPY go.sum .
COPY go.mod .

# чтобы не было зависимостей от Си-библиотек в конечном бинарнике,
# потому что их не будет в контейнере FROM scratch
ENV CGO_ENABLED=0 

RUN go build -o reviewer-assignment-server .

# чтобы конечный контейнер не содержал в себе код, только самодостаточный бинарник
FROM scratch
COPY --from=build /go/src/reviewer-assignment-server ./
EXPOSE 8080/tcp
ENTRYPOINT ["./reviewer-assignment-server"]
