docker-run:
	docker-compose up -d app postgres memcached migrate

docker-stop:
	docker-compose stop

docker-rm:
	docker-compose rm

docker-cleanup: docker-stop docker-rm

docker-test:
	docker-compose up -d apptest postgres_test memcached migrate --rm
	docker-compose stop
	docker-compose rm
	
host-build:
	go build main.go

host-run:
	go run main.go server

test:
	go test github.com/dimoiko100/bl-assignment/store github.com/dimoiko100/bl-assignment/api
