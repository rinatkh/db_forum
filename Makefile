all: docker-build docker-run

docker-build:
	DOCKER_BUILDKIT=1 docker build -t park .

docker-run:
	docker rm -f db_forum
	docker run -p 5000:5000 -p 5432:5432 --name db_forum -t db_forum

mod:
	go mod tidy && go mod vendor && go install ./...
