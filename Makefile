.PHONY: docker
docker:
	@rm demo || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -o demo .
	@docker rmi -f wenbin/k8s-demo:v0.0.1
	@docker build -t wenbin/k8s-demo:v0.0.1 .
