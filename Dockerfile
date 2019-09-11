FROM golang as build-api

WORKDIR /app
ADD go.mod go.sum /app/
RUN go mod download
ADD Makefile *.go /app/
RUN make build-linux64

FROM node:alpine as build-ui

WORKDIR /app
ADD ./ui/package.json ./ui/yarn.lock /app/
RUN yarn
ADD ./ui /app
RUN yarn build

FROM scratch

COPY --from=build-api /app/bin/lumberman-web-client /
COPY --from=build-ui /app/public /ui/public

EXPOSE 80

ENTRYPOINT ["/lumberman-web-client"]
