FROM golang as builder
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN make

FROM alpine:3.6
RUN apk add --no-cache tzdata ca-certificates
COPY --from=builder /app/bin/web /usr/bin/web
COPY --from=builder /app/bin/trigger /usr/bin/trigger
ENTRYPOINT ["/usr/bin/web"]