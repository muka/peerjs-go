

kill-server:
	kill -9 $(lsof -t -i tcp:9000)

peerjs/interop/js:
	cd interop/js && npm run serve

peerjs/server/run:
	docker run --rm --name peerjs-server -p 9000:9000 -d peerjs/peerjs-server
	docker logs -f peerjs-server
