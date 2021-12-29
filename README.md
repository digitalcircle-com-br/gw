# GATEWAY


## TL;DR

```bash
go install github.com/digitalcircle-com-br/gw
gw
```


## Sample Config

```yaml
addr: ":8443"
acme: false
helmet: false
https: true
cors: true
cert: ./etc/server.crt
key: ./etc/server.key
insecure: true

env:
  DEBUG: 1
  TRACE: 0
  MONGO: mongodb://mongo/myme
  CORS: SAME
  LOGSDIR: ./log
  DEBUG_GW: 1
  DEBUG_CONFIG: 0
  DEBUG_SERVERLESS: 0
  DEBUG_SERVERLESS_LOG: 1

procs:
  auth: ./bin/auth -a :8081
  file: ./bin/file -a :8082
  config: ./bin/config -a :8083
  myme: ./bin/myme -a :8084
  admmyme: ./bin/admmyme -a :8085
  inboundmailprocessor: ./bin/inboundmailprocessor

cron:
  - cron: "*/5 * * * * *"
    cmd: ./bin/temp

routes:
  api.dc.local/api/auth/: http://localhost:8081 #exec://./bin/auth
  api.dc.local/api/file/: http://localhost:8082 #exec://./bin/file
  api.dc.local/api/myme/: http://localhost:8084 #exec://./bin/myme
  api.dc.local/api/admmyme/: http://localhost:8085 #exec://./bin/myme

  myme.dc.local/config: http://localhost:8083 #exec://./bin/config
  myme.dc.local/: >-
    static://
    ./client/user/webroot,
    ./client/user/dist,
    ../weblib/webroot,
    ../weblibx/webroot

  admmyme.dc.local/config: http://localhost:8083 #exec://./bin/config
  admmyme.dc.local/: >-
    static://
    ./client/admin/webroot,
    ./client/admin/dist,
    ../weblib/webroot,
    ../weblibx/webroot
```


## Running in docker

```shell script
docker volume create gw

docker network create dc

docker run -d --restart=always -p 80:80 -p 443:443 -v gw:/srv --name gw --network dc gw
```

## Ref Dockerfile (/deploy)

```dockerfile
FROM alpine:latest
VOLUME /srv
COPY ./bin/gw /bin/gw
COPY ./etc/gw.yaml /srv
ENV DEBUG 1
ENV TRACE 1
ENV CONFIG /srv/gw.yaml
WORKDIR /srv
CMD "/bin/gw"
```