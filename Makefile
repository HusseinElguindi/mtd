mtd_proto:
	protoc -I ./protos/mtd --go-grpc_out=./protos/mtd --go_out=./protos/mtd ./protos/mtd/mtd.proto --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative