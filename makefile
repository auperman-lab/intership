proto-authsvc:
	protoc -I pkg/authsvc/proto --go_out=./pkg/authsvc --go-grpc_out=./pkg/authsvc pkg/authsvc/proto/user.proto\
 	--grpc-gateway_out=pkg/authsvc/pb --grpc-gateway_opt=paths=source_relative

proto-usersvc:
	protoc -I pkg/usersvc/proto --go_out=./pkg/usersvc --go-grpc_out=./pkg/usersvc pkg/usersvc/proto/usermenu.proto\
 	--grpc-gateway_out=pkg/usersvc/pb --grpc-gateway_opt=paths=source_relative

proto-adminsvc:
	protoc -I pkg/adminsvc/proto --go_out=./pkg/adminsvc --go-grpc_out=./pkg/adminsvc pkg/adminsvc/proto/admin.proto\
 	--grpc-gateway_out=pkg/adminsvc/pb --grpc-gateway_opt=paths=source_relative

proto-authsvc-delete:
	rm pkg/authsvc/pb/*.go

proto-usersvc-delete:
	rm pkg/usersvc/pb/*.go

proto-adminsvc-delete:
	rm pkg.adminsvc/pb/*.go

server:
	go run cmd/main.go