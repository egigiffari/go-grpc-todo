gen:
	@protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--proto_path=proto proto/*.proto

build:
	@go build -o ./bin/ctl main.go

run: build
	@bin/ctl