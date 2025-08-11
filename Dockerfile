FROM golang:1.23 AS build
RUN mkdir -p /go/src/github.com/josephburnett/jd
WORKDIR /go/src/github.com/josephburnett/jd
COPY . .
RUN set -eux; \
  export GOROOT="$(go env GOROOT)"; \
  make build
RUN cd release ; ln -s jd jd-github-action
FROM scratch
COPY --from=build /go/src/github.com/josephburnett/jd/release/jd /jd
COPY --from=build /go/src/github.com/josephburnett/jd/release/jd-github-action /jd-github-action
ENTRYPOINT ["/jd"]