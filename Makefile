test:
	go test ./... -v -cover -covermode=count -coverprofile=count.out
	go tool cover -func=count.out

lint:
	golint
vet:
	go vet ./...

mockgen:
	mockgen -source tmux_wrapper.go -destination mocks/command_executor.go -package mocks  ICommandExecutor

prepare:
	go mod download
	go install golang.org/x/lint/golint@latest
