# comment it for local development
ARG GO_LANG_TAG
FROM golang:${GO_LANG_TAG} as Builder
RUN apk add --update  && \
    apk add --no-cache alpine-conf tzdata git openssh

ADD ./src /go/src/app
ADD ./src/log /go/log
ADD ./src/config /go/config

RUN  mkdir -p /root/.ssh
COPY id_rsa /root/.ssh/id_rsa
COPY known_hosts /root/.ssh/known_hosts
RUN chmod -R 400 /root/.ssh


RUN git config --global url."git@***:".insteadOf "***"


RUN cd /go/src/app && \
    go install app
RUN rm -rf /root/.ssh
# comment it for local development
ARG ${ALPINE_TAG}
FROM alpine:${ALPINE_TAG} as App

COPY --from=Builder /go/bin/* /go/bin/app
COPY --from=Builder /go/log /go/log
COPY --from=Builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Nicosia
USER nonroot:nonroot

WORKDIR "/go"
ENTRYPOINT ["/go/bin/app"]
