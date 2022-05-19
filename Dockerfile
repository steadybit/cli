FROM node:16 as builder

COPY ./ ./build
WORKDIR ./build
RUN npm ci && \
    npm run ci && \
    npm pack

FROM node:16

COPY --from=builder ./build/steadybit-*.tgz steadybit.tgz
RUN npm -g install ./steadybit.tgz && rm ./steadybit.tgz

ENTRYPOINT ["steadybit"]
