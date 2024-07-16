FROM golang:1.22-alpine as builder

WORKDIR /usr/src/app
COPY . .
RUN go build -o importer cmd/main.go

FROM scratch

COPY --from=builder /usr/src/app/importer /usr/local/bin/importer

CMD ["importer"]
