.PHONY: default
default: all

.PHONY: all clean
all: hello

##### COMMAND-LIKE PHONY TARGETS #####
clean:
	make -C backend clean

##### SHORTCUT PHONY TARGETS #####
.PHONY: hello

hello: bin/hello

proto-go: proto/.sentinel

##### ACTUAL TARGETS #####

bin/hello:

PROTO_SOURCES := $(shell find proto -type f -name "*.proto")
proto/.sentinel: ${PROTO_SOURCES}
	protoc --go_out=proto/ --go_opt=paths=source_relative \
			--go-grpc_out=proto/ --go-grpc_opt=paths=source_relative \
			-I proto/ \
			${PROTO_SOURCES}
	cat $$(find proto -type f -name "*.pb.go") | md5sum > proto/.sentinel

# protoc --ts_proto_out=./@tagioalisi/proto --ts_proto_opt=env=browser,outputServices=nice-grpc,outputServices=generic-definitions,outputJsonMethods=false,useExactTypes=false -I ../proto ../proto/*



webapp/node_modules: webapp/package-lock.json webapp/package.json
