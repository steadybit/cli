FROM node:18-alpine as builder

COPY ./ ./build
WORKDIR ./build
RUN npm ci && \
    npm run build && \
    npm pack

FROM node:18-alpine

COPY --from=builder ./build/steadybit-*.tgz steadybit.tgz

# Git is required to automatically add tags to a policy binding
RUN apk add --no-cache git && \
  npm -g install ./steadybit.tgz && \
  rm ./steadybit.tgz

ENTRYPOINT ["steadybit"]
