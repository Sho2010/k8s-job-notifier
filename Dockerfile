FROM golang:1.19.4 as build-env

ENV GO111MODULE=on
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

ADD . /go/src/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app -v

#---

FROM gcr.io/distroless/base
COPY --chown=nonroot:nonroot --from=build-env /go/bin/app /
USER nonroot

CMD ["/app"]
