rpc:
	protoc -I ./rpc/ \
	--go_out ./rpc --go_opt=paths=source_relative \
	--gin_out ./rpc --gin_opt=paths=source_relative ./rpc/demo/v0/demo.proto

	protoc -I ./rpc/ \
	--go_out=./rpc  --go_opt=paths=source_relative  \
	--go-grpc_out=./rpc --go-grpc_opt=paths=source_relative ./rpc/grpc/v0/grpc_example.proto

	# 老版本grpc生成方式
	# protoc -I ./rpc/ \
	# --go_out=plugins=grpc:./rpc  --go_opt=paths=source_relative ./rpc/grpc/v0/grpc_example.proto