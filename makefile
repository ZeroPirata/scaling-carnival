.PHONY: proto server client clean

# Define onde estão os arquivos proto e onde gerar o código Go
GOOGLEAPIS_PROTOS = googleapis_protos
PROTO_DIR=proto
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)
GENERATED_DIR=server/pb # Ou onde você preferir colocar o código gerado

# Gera o código Go a partir dos arquivos .proto
proto-win:
	protoc -I proto -I $(GOOGLEAPIS_PROTOS) --go_out=$(GENERATED_DIR) \
	  --go_opt=paths=source_relative --go-grpc_out=$(GENERATED_DIR) \
	  --go-grpc_opt=paths=source_relative $(PROTO_FILES)
	@echo "Generated Go code from proto files."


# Compila e executa o servidor
server: proto
	go run ./server/main.go

# Compila e executa o cliente (exemplo simples)
client: proto
	go run ./client/main.go

# Limpa os arquivos gerados
clean:
	rm -rf $(GENERATED_DIR)/*