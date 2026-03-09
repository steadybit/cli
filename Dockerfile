FROM node:22-alpine as builder

COPY ./ ./build
WORKDIR ./build
RUN npm ci && \
    npm run build && \
    npm pack

FROM node:22-alpine

COPY --from=builder ./build/steadybit-*.tgz steadybit.tgz

RUN apk update && apk upgrade --no-cache && rm -rf /var/cache/apk/* && \
    npm install -g npm@11 && \
    npm -g install ./steadybit.tgz && \
  rm ./steadybit.tgz

ENTRYPOINT ["steadybit"]
