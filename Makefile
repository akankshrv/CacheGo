build:
	go build -o bin/CacheGo

run: build
	./bin/CacheGo

runfollower: build
	./bin/CacheGo --listenaddr :4000 --leaderaddr :3000

test: 
	@go test -v ./...