gen:
	@protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--proto_path=proto proto/*.proto

run_client:
	@go run cmd/client/main.go

run_server:
	@go run cmd/server/main.go