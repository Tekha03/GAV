FROM golang:1.25 AS build

WORKDIR /app

COPY go.work go.work.sum ./
COPY api ./api
COPY docs ./docs
COPY messenger ./messenger
COPY shared ./shared
COPY social_network ./social_network

ARG SERVICE_DIR
ARG COMMAND_PATH

RUN cd "${SERVICE_DIR}" && go build -o /out/app "${COMMAND_PATH}"

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]
