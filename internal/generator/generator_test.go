package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/jamillosantos/protoc-gen-nextra/internal/generator"
)

// fileChecks maps each expected output file to a list of strings that must
// appear somewhere in the generated content. This avoids brittle exact-match
// comparisons while still asserting that the key structural elements are present.
var fileChecks = map[string][]string{
	"rpc/greeter/v1.gen.mdx": {
		// Services
		"## Greeter",
		"## Farewell",
		// Methods and RPC badges
		"SayHello",
		"SayHelloStream",
		"SayGoodbye",
		">UNARY<",
		">SERVER STREAM<",
		// Request/response types
		"HelloRequest",
		"HelloReply",
		"GoodbyeRequest",
		"GoodbyeReply",
		// Error codes
		"`NOT_FOUND`",
		"`INVALID_ARGUMENT`",
		// Error detail type (external — plain code, not linked)
		"google.rpc.BadRequest",
		// Error JSON payload
		`"field_violations"`,
	},
	"rpc/greeter/v2.gen.mdx": {
		"SayHello",
		"SayHelloStream",
		"SayHelloMany",
		">UNARY<",
		">SERVER STREAM<",
		">CLIENT STREAM<",
		"HelloRequest",
		"HelloReply",
	},
	"rpc/notifier/v1.gen.mdx": {
		// Methods and all four RPC types
		"Send",
		"Subscribe",
		"Acknowledge",
		"Chat",
		">UNARY<",
		">SERVER STREAM<",
		">CLIENT STREAM<",
		">BIDI STREAM<",
		// Error codes
		"`NOT_FOUND`",
		"`INVALID_ARGUMENT`",
		"`UNAUTHENTICATED`",
		// Local detail type link (TypePreview)
		"shared.notifications.v1.ValidationFailure",
	},
	"shared/notifications/v1.gen.mdx": {
		"ValidationFailure",
		"DeliveryStatus",
		"Recipient",
	},
}

// TestGenerateFile runs the generator against a pre-compiled descriptor set
// and checks that each output file contains the expected key strings.
// Generate the fixture with: make testdata
func TestGenerateFile(t *testing.T) {
	req := loadRequest(t, filepath.Join("..", "..", "testdata", "all.pb"))

	opts := protogen.Options{}
	plugin, err := opts.New(req)
	if err != nil {
		t.Fatalf("create plugin: %v", err)
	}

	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}
		if err := generator.GenerateFile(plugin, f, generator.Config{TypePreviews: true}); err != nil {
			t.Fatalf("generate %s: %v", f.Desc.Path(), err)
		}
	}

	resp := plugin.Response()
	if resp.GetError() != "" {
		t.Fatalf("plugin error: %s", resp.GetError())
	}

	for _, rf := range resp.File {
		checks, ok := fileChecks[rf.GetName()]
		if !ok {
			continue
		}
		t.Run(rf.GetName(), func(t *testing.T) {
			content := rf.GetContent()
			for _, want := range checks {
				if !strings.Contains(content, want) {
					t.Errorf("expected output to contain %q", want)
				}
			}
		})
	}
}

// loadRequest reads a serialised FileDescriptorSet produced by
//
//	protoc --descriptor_set_out=testdata/all.pb --include_source_info --include_imports ...
//
// and converts it to a CodeGeneratorRequest that asks to generate all files.
func loadRequest(t *testing.T, pbPath string) *pluginpb.CodeGeneratorRequest {
	t.Helper()

	b, err := os.ReadFile(pbPath)
	if err != nil {
		t.Skipf("descriptor set not found at %s — run 'make testdata' first: %v", pbPath, err)
	}

	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(b, &fds); err != nil {
		t.Fatalf("unmarshal FileDescriptorSet: %v", err)
	}

	var toGenerate []string
	for _, fd := range fds.File {
		toGenerate = append(toGenerate, fd.GetName())
	}

	return &pluginpb.CodeGeneratorRequest{
		FileToGenerate: toGenerate,
		ProtoFile:      fds.File,
	}
}
