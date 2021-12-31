clean:
	rm -rd ./dist || true

mk-dist:
	mkdir dist

build: clean mk-dist
	go build -o ./dist/dailydo .

build-prod:
	GOOS="linux" GOARCH="amd64" go build -o ./dist/dailydo .

publish:
	docker build -t dailydo-server:latest -f ./Dockerfile .

release: clean mk-dist build-prod publish

run: build
	@ docker-compose up -d
	@ ZIPKIN="localhost:9411" MONGODB="mongodb://localhost:27017/users" ./dist/Auth

mocks:
	GO111MODULE=off go get -u github.com/vektra/mockery/.../
	$(GOPATH)/bin/mockery -dir app/repositories -all -output app/repositories/mocks -note 'Regenerate this file using `make mocks`.'
	$(GOPATH)/bin/mockery -dir app/services -all -output app/services/mocks -note 'Regenerate this file using `make mocks`.'

