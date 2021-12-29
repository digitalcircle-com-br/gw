FROM alpine:latest
RUN apk update && apk add ca-certificates iputils curl vim iputils && rm -rf /var/cache/apk/*
RUN update-ca-certificates
VOLUME /srv
COPY ./gw /bin/gw
COPY ./gw.yaml /srv/gw.yaml

ENV DEBUG 1
ENV TRACE 1
ENV DEBUG_GW 1
ENV DEBUG_STATIC 1
ENV DEBUG_CORS 1
ENV DEBUG_ENV 1
ENV CONFIG /srv/gw.yaml
ENV CONTROL_KEY "please set something to prevent unnatended control"
WORKDIR /srv
CMD "/bin/gw"
