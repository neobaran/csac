FROM golang:1-alpine as builder

RUN apk --no-cache --no-progress add make git

WORKDIR /go/csac

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.12
RUN apk update \
  && apk add --no-cache ca-certificates tzdata \
  && update-ca-certificates

COPY --from=builder /go/csac/csac /usr/bin/csac

ENTRYPOINT [ "/usr/bin/csac" ]