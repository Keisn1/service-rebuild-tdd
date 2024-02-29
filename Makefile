##
# Project Title
#
# @file
# @version 0.1

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

test:
	docker-compose exec golang-server go test ./...

unit_test:
	go test ./...

restart:
	docker-compose down
	docker-compose up

.PHONY: build up down restart

# end
