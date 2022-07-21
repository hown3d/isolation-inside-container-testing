FROM golang:1.18 as build

ENV CGO_ENABLED=1

WORKDIR /workspace
COPY . .
RUN \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  go build -o /libcontainer-test .

FROM gcr.io/distroless/base:debug
COPY --from=build /libcontainer-test /main

ENTRYPOINT [ "/main" ]
