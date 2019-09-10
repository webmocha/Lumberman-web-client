FROM golang as build

WORKDIR /app
ADD . /app
RUN cd /app && \
  make build-linux64

FROM alpine

COPY --from=build /app/bin/lumberman-web-client /

EXPOSE 80

ENTRYPOINT ["/lumberman-web-client"]
