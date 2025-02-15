FROM golang:1.18 AS build
RUN mkdir -p /go/src/github.com/josephburnett/jd
WORKDIR /go/src/github.com/josephburnett/jd
COPY . .
RUN set -eux; \
  export GOROOT="$(go env GOROOT)"; \
  make build
  make build-action
FROM scratch
COPY --from=build /go/src/github.com/josephburnett/jd/release/jd /jd
COPY --from=build /go/src/github.com/josephburnett/jd/release/jd-acion /jd-action
ENTRYPOINT ["/jd"]