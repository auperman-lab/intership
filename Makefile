proto-run:
	protoc -I internal/proto --go_out=./internal --go-grpc_out=./internal internal/proto/*.proto

proto-delete:
	rm pb/*.go

server-run:
	cd cmd && go run main.go