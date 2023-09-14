server:
	cd cmd && go run main.go

proto-delete:
	rm internal/pb/*.go

proto:
	protoc -I internal/proto --go_out=./internal --go-grpc_out=./internal internal/proto/*.proto