rpc:
	protoc -I ./rpc/ \
	--go_out ./rpc --go_opt=paths=source_relative \
	--gin_out ./rpc --gin_opt=paths=source_relative ./rpc/demo/v0/demo.proto
