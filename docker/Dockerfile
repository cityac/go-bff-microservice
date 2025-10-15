FROM golang:1.23.4-bookworm AS builder

ARG CI_JOB_TOKEN
ARG buildenv

WORKDIR /build

COPY . .
RUN git submodule update --init
RUN git config --global url."https://gitlab-ci-token:${CI_JOB_TOKEN}@gitlab.HIDDEN.com/".insteadOf "https://gitlab.HIDDEN.com/"
RUN go mod download

RUN echo "Building with build environment: ${buildenv}"
RUN GOOS=linux go build -tags "${buildenv}" -o bff cmd/*

FROM debian:bookworm
WORKDIR /app

COPY --from=builder /build/bff .
COPY --from=builder /build/scripts/launcher.sh .

RUN apt update && apt install -y curl jq && apt clean

CMD ./launcher.sh

