FROM node:24-alpine AS builder

COPY ./ ./build
WORKDIR ./build
RUN npm ci && \
    npm run build && \
    npm pack

FROM node:24-alpine

COPY --from=builder ./build/steadybit-*.tgz steadybit.tgz

RUN apk update && apk upgrade --no-cache && rm -rf /var/cache/apk/* && \
    npm -g install npm@11 && \
    npm -g install ./steadybit.tgz && \
    rm ./steadybit.tgz

ENTRYPOINT ["steadybit"]
