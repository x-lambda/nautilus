rpc:
	protoc -I ./api/ \
	--go_out ./api --go_opt=paths=source_relative \
	--gin_out ./api --gin_opt=paths=source_relative ./api/demo/v0/demo.proto

	protoc -I ./api/ \
	--go_out ./api --go_opt=paths=source_relative \
	--gin_out ./api --gin_opt=paths=source_relative ./api/pension/v0/service.proto

	# protoc -I ./api/ \
	# --go_out=./api  --go_opt=paths=source_relative  \
	# --go-grpc_out=./api --go-grpc_opt=paths=source_relative ./api/grpc_demo/v0/grpc_example.proto

	# 老版本grpc生成方式
	# protoc -I ./rpc/ \
	# --go_out=plugins=grpc:./rpc  --go_opt=paths=source_relative ./rpc/grpc/v0/grpc_example.proto