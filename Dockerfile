FROM golang:1.13 as build

COPY . code
WORKDIR code
RUN CGO_ENABLED=0 go build -o build/apiserver cmd/api/main.go

FROM alpine:3.7
COPY --from=build /go/code/build/apiserver /usr/local/bin/apiserver

EXPOSE 8080

CMD ["apiserver"]
