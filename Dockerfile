FROM alpine:3.12
RUN apk update \
  && apk add --no-cache ca-certificates tzdata \
  && update-ca-certificates

COPY csac /usr/bin/csac

ENTRYPOINT [ "/usr/bin/csac" ]