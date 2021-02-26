

kill-server:
	kill -9 $(lsof -t -i tcp:9000)

peerjs/interop/js:
	cd interop/js && npm run serve

peerjs/server/run:
	docker run --rm --name peerjs-server -p 9000:9000 -d peerjs/peerjs-server
	docker logs -f peerjs-server

docker/build:
	mkdir -p build
	CGO_ENABLED=0 go build -o build/peer-server-amd64 cmd/server/main.go
