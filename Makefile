binary_name=statuscentral
environment=local

build: clean
	@CGO_ENABLED=0 go build -v -a -o $$GOPATH/bin/${binary_name}/${binary_name} main.go
	@cp ./${binary_name}.${environment}.yaml $$GOPATH/bin/${binary_name}/${binary_name}.yaml
	@cp -r ./templates $$GOPATH/bin/${binary_name}/
	@cp -R ./static $$GOPATH/bin/${binary_name}/

clean:
	@docker-compose down
	@rm -rf $$GOPATH/bin/${binary_name}

docker:
	@docker-compose build

run: build
	@cd $$GOPATH/bin/${binary_name}/ && ./${binary_name}

start:
	@docker-compose up -d --build