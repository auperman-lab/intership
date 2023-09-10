proto:
	protoc -I pkg/authsvc/proto --go_out=./pkg/authsvc --go-grpc_out=./pkg/authsvc pkg/authsvc/proto/user.proto\
 	--grpc-gateway_out=pkg/authsvc/pb --grpc-gateway_opt=paths=source_relative

proto-delete:
	rm pkg/authsvc/pb/*.go

server:
	go run cmd/main.go