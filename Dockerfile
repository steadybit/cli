# SPDX-License-Identifier: MIT
# SPDX-FileCopyrightText: 2026 Steadybit GmbH

# Built natively for the build platform and cross-compiled, so a multi-arch build does
# not run the Go toolchain under emulation.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY api ./api
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath \
    -ldflags "-s -w -X github.com/steadybit/cli/v6/internal/platform.Version=${VERSION}" \
    -o /steadybit ./cmd/steadybit

# Alpine rather than scratch: pipelines use the image with a shell, and the e2e suite
# installs expect into it.
FROM alpine:3
RUN apk upgrade --no-cache
COPY --from=builder /steadybit /usr/local/bin/steadybit
# A container is started fresh for each run, often in a pipeline whose CI variables are
# not passed in, and is updated by pulling a newer image, so the daily release check
# would only ask GitHub every time and give the wrong advice.
ENV STEADYBIT_NO_UPDATE_CHECK=1
ENTRYPOINT ["steadybit"]
