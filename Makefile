.PHONY: build install testdata generate test update-golden lint

BINARY := protoc-gen-nextra

build:
	go build -o bin/$(BINARY) ./cmd/$(BINARY)

install:
	go install ./cmd/$(BINARY)

# Compile the test proto into a FileDescriptorSet used by the Go tests.
# Requires: protoc (e.g. brew install protobuf)
testdata/greeter.pb: testdata/proto/greeter.proto
	protoc \
		--descriptor_set_out=testdata/greeter.pb \
		--include_source_info \
		-I testdata/proto \
		testdata/proto/greeter.proto

testdata: testdata/greeter.pb

# Run the plugin against the test proto and write MDX pages to testdata/content/.
generate: build testdata/greeter.pb
	mkdir -p testdata/content
	protoc \
		--plugin=protoc-gen-nextra=bin/$(BINARY) \
		--nextra_out=testdata/content \
		-I testdata/proto \
		testdata/proto/greeter.proto

test: generate
	go test ./...

# Re-generate golden files from current output.
update-golden: testdata/greeter.pb
	UPDATE_GOLDEN=1 go test ./...

lint:
	golangci-lint run ./...
