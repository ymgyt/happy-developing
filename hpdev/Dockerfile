FROM golang:1.11.2-alpine3.8 as build

ENV GO111MODULE=off

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/ymgyt/happy-developing/hpdev

COPY . ./

# RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o hpdev
RUN CGO_ENABLED=0 GOOS=linux go build -o /hpdev


FROM alpine:3.8

# TODO specify USER !!

WORKDIR /root

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /hpdev .
COPY --from=build /go/src/github.com/ymgyt/happy-developing/hpdev/static ./static
COPY --from=build /go/src/github.com/ymgyt/happy-developing/hpdev/templates ./templates
COPY --from=build /go/src/github.com/ymgyt/happy-developing/hpdev/credentials ./credentials

EXPOSE 8123

ARG GCP_PROJECT_ID

ENV GCP_PROJECT_ID $GCP_PROJECT_ID
ENV GCP_CREDENTIAL_JSON=/root/credentials/hpdev.json
ENV APP_ROOT=/root

ENTRYPOINT ["./hpdev"]