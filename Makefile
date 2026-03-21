.PHONY: build install testdata generate test update-golden lint clean

BINARY := protoc-gen-nextra

PROTO_FILES := \
	testdata/proto/rpc/greeter/v1/greeter.proto \
	testdata/proto/rpc/greeter/v2/greeter.proto \
	testdata/proto/rpc/notifier/v1/notifier.proto \
	testdata/proto/shared/notifications/v1/common.proto

build:
	go build -o bin/$(BINARY) ./cmd/$(BINARY)

install:
	go install ./cmd/$(BINARY)

# Compile all test protos into a single FileDescriptorSet used by the Go tests.
# Requires: protoc (e.g. brew install protobuf)
testdata/all.pb: $(PROTO_FILES)
	protoc \
		--descriptor_set_out=testdata/all.pb \
		--include_source_info \
		--include_imports \
		-I testdata/proto \
		$(PROTO_FILES)

testdata: testdata/all.pb

# Run the plugin against the test protos and write MDX pages to testdata/content/.
generate: build testdata/all.pb
	mkdir -p testdata/content
	protoc \
		--plugin=protoc-gen-nextra=bin/$(BINARY) \
		--nextra_out=testdata/content \
		-I testdata/proto \
		$(PROTO_FILES)

test: generate
	go test ./...

# Re-generate golden files from current output.
update-golden: testdata/all.pb
	UPDATE_GOLDEN=1 go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ testdata/all.pb
	find testdata/content -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} +
