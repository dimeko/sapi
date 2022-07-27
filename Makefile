include .env

run:
	docker-compose up -d app postgres memcached migrate
	docker-compose up -d postgres && \
	docker-compose up -d migrate  && \
	docker-compose up -d memcached  && \
	docker-compose up app

d-stop:
	docker-compose stop

d-rm:
	docker-compose rm

cleanup: d-stop d-rm

test:
	docker-compose up -d postgres_test && \
	docker-compose up -d migrate_test  && \
	docker-compose up apptest
	docker-compose down -d postgres_test migrate_test apptest
	
migrate-run:
	docker-compose run migrate || docker-compose up migrate
