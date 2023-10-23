.PHONY: default
default: all


##### COMMAND-LIKE PHONY TARGETS #####
.PHONY: all clean clean-bin clean-webapp clean-proto wire
all: hello id-generator

clean: clean-bin clean-webapp clean-proto

clean-bin:
	rm -rf bin

clean-webapp:
	rm -rf webapp/dist

clean-proto:
	rm -rf proto/*.pb.go proto/.sentinel

wire:
	wire ./...

##### SHORTCUT PHONY TARGETS #####
.PHONY: hello id-generator proto-go

hello: bin/hello

id-generator: bin/bounce-id-generator

proto-go: proto/.sentinel

##### DEPENDENCIES FOR REBUILDING GO STUFF ##### 
# Includes all Go files, stuff that generates Go files, etc
GO_BUILD_DEPS := proto-go go.mod go.sum $(shell find . -type f -name "*.go")

### ID generator
bin/bounce-id-generator: ${GO_BUILD_DEPS}
	go build -o bin/bounce-id-generator ./id-generator/main/

bin/hello: ${GO_BUILD_DEPS}
	go build -o bin/bounce-hello ./helloworld


PROTO_SOURCES := $(shell find proto -type f -name "*.proto")
proto/.sentinel: ${PROTO_SOURCES}
	protoc --go_out=proto/ --go_opt=paths=source_relative \
			--go-grpc_out=proto/ --go-grpc_opt=paths=source_relative \
			-I proto/ \
			${PROTO_SOURCES}
	cat $$(find proto -type f -name "*.pb.go") | md5sum > proto/.sentinel

# protoc --ts_proto_out=./@tagioalisi/proto --ts_proto_opt=env=browser,outputServices=nice-grpc,outputServices=generic-definitions,outputJsonMethods=false,useExactTypes=false -I ../proto ../proto/*


webapp/node_modules: webapp/package-lock.json webapp/package.json
