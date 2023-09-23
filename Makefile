proto-deepthought:
	protoc --go_out=. --go-grpc_out=. deepthought.proto

proto-deepthought-doc:
	protoc --doc_out=html,index.html:./docs deepthought.proto

server:
	go run ./go/server/*.go

client:
	go run ./go/client/main.go 127.0.0.1:13333

