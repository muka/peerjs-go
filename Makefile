
IMAGENAME ?= opny/peer-server
BUILDPATH ?= ./build

server/run:
	go run cmd/server/main.go

kill-server:
	kill -9 $(lsof -t -i tcp:9000)

peerjs/interop/js:
	cd interop/js && npm run serve

peerjs/server/run:
	docker stop peerjs-server || true
	docker run --rm --name peerjs-server -p 9000:9000 -d peerjs/peerjs-server --port 9000 --path /
	docker logs -f peerjs-server

docker/run:
	docker run --name peer-server -it --rm -p 9000:9000 $(IMAGENAME)

build: build/amd64 build/arm64 build/arm

build/amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILDPATH}/amd64 .

build/arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${BUILDPATH}/arm64 .

build/arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o ${BUILDPATH}/arm7 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o ${BUILDPATH}/arm6 .

docker/build: build docker/build/amd64 docker/build/arm64 docker/build/arm

docker/build/manifest:

	docker manifest push --purge ${IMAGENAME} || true

	docker manifest create \
		${IMAGENAME} \
		--amend ${IMAGENAME}-amd64 \
		--amend ${IMAGENAME}-arm64 \
		--amend ${IMAGENAME}-arm6 \
		--amend ${IMAGENAME}-arm7
	
	docker manifest annotate ${IMAGENAME} ${IMAGENAME}-amd64 --arch amd64 --os linux
	docker manifest annotate ${IMAGENAME} ${IMAGENAME}-arm64 --arch arm64 --os linux
	docker manifest annotate ${IMAGENAME} ${IMAGENAME}-arm6 --arch arm --variant v6 --os linux
	docker manifest annotate ${IMAGENAME} ${IMAGENAME}-arm7 --arch arm --variant v7 --os linux

	docker manifest push ${IMAGENAME}

docker/build/amd64:
	docker build . -t ${IMAGENAME}-amd64 --build-arg ARCH=amd64

docker/build/arm64:
	docker build . -t ${IMAGENAME}-arm64 --build-arg ARCH=arm64

docker/build/arm:
	docker build . -t ${IMAGENAME}-arm6 --build-arg ARCH=arm --build-arg ARM=7
	docker build . -t ${IMAGENAME}-arm7 --build-arg ARCH=arm --build-arg ARM=6

docker/push: docker/build docker/push/amd64 docker/push/arm64 docker/push/arm docker/build/manifest
	docker manifest push ${IMAGENAME}

docker/push/amd64:
	docker push ${IMAGENAME}-amd64

docker/push/arm64:
	docker push ${IMAGENAME}-arm64

docker/push/arm:
	docker push ${IMAGENAME}-arm6
	docker push ${IMAGENAME}-arm7

