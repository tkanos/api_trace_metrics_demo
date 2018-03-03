.PHONY: build install test deploy smoke es_mem_issue all

es_mem_issue:
	sudo sysctl -w vm.max_map_count=262144

install:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/bluehoodie/smoke
	go get -u golang.org/x/lint/golint
	dep ensure

build:
	tput setaf 3; echo "\t\t===== BUILD ====="; tput sgr0;
	CGO_ENABLED=0 go build -o ./pixel-api -a -ldflags '-s' -installsuffix cgo main.go

test:
	tput setaf 3; echo "\t\t===== UNIT TEST ====="; tput sgr0;
	go test -i
	go test -race -v `go list ./... | grep -v /vendor/`
	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status
	go vet `go list ./... | grep -v /vendor/`

deploy: build
	tput setaf 3; echo "\t\t===== DEPLOY ====="; tput sgr0;
	docker build -t pixel-api .
	rm pixel-api
smoke: deploy
	tput setaf 3; echo "\t\t===== SMOKE / INTEGRATION TEST ====="; tput sgr0;
	docker run -d --name=smoke -p 5000:5000 pixel-api
	smoke -file smoke.json -url http://localhost:5000 -v
	docker stop smoke
	docker rm smoke
all: test smoke
	tput setaf 3; echo "\t\t===== LAUNCH SERVER ====="; tput sgr0;
	docker-compose up
