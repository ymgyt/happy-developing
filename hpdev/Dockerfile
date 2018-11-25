FROM golang:1.11.2-alpine3.8 as build

WORKDIR /hpdev

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o hpdev


FROM alpine:3.8

# TODO specify USER !!

WORKDIR /root/

COPY --from=build /hpdev .
COPY --from=build /hpdev/static .
COPY --from=build /hpdev/templates .

EXPOSE 8123
ENV APP_ROOT=/root

ENTRYPOINT ["./hpdev"]