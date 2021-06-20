all: build run

build:
	docker build -t libregram/server -f Dockerfile .

run:
	docker-compose up -d
