# SPDX-License-Identifier: MIT
# SPDX-FileCopyrightText: 2026 Steadybit GmbH

# Built natively for the build platform and cross-compiled, so a multi-arch build does
# not run the Go toolchain under emulation.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder
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
ENTRYPOINT ["steadybit"]
