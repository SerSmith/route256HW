generate-all:
	cd checkout && GOOS=linux GOARCH=amd64 make generate 
	cd loms && GOOS=linux GOARCH=amd64 make generate

build-all:
	cd checkout && GOOS=linux GOARCH=amd64 make build
	cd loms && GOOS=linux GOARCH=amd64 make build
	cd notifications && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	sudo docker compose up --force-recreate --build  -d
	cd checkout && make goose-up
	cd loms && make goose-up


precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit