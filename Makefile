dt:=$(shell date)

ver:
	echo "package gw\nconst VER= \"${dt}\"\n" > src/gw/ver.go

docker: ver
	#GOOS=linux go build -ldflags "-linkmode external -extldflags -static" -o ./bin/gw_linux -tags prd ./src
	CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/gw_linux -tags prd ./src
	docker build -t reg.digitalcircle.com.br/gw .
	#docker push reg.digitalcircle.com.br/gw

reload:
	docker rm -f gw
	docker run -d --restart=always -p 80:80 -p 443:443 -v /srv/gw/srv:/srv --name gw --network net reg.digitalcircle.com.br/gw
pub: docker
	docker save reg.digitalcircle.com.br/gw > /tmp/gw.tar
	scp /tmp/gw.tar dc01:/tmp
	ssh dc01 docker load < /tmp/gw.tar
	- ssh dc01 docker rm -f gw
	ssh dc01 docker run -d --restart=always -p 80:80 -p 443:443 -v /srv/gw/srv:/srv --name gw --network net reg.digitalcircle.com.br/gw

local: ver
	go build -o ~/.go/bin/gw ./src
