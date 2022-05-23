FROM node:16-alpine as builder

COPY ./ ./build
WORKDIR ./build
RUN npm ci && \
    npm run build && \
    npm pack

FROM node:16-alpine

COPY --from=builder ./build/steadybit-*.tgz steadybit.tgz

# Git is required to automatically add tags to a service definition
RUN apk add --no-cache git && \
  npm -g install ./steadybit.tgz && \
  rm ./steadybit.tgz

ENTRYPOINT ["steadybit"]
