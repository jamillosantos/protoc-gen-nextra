package generator_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/jamillosantos/protoc-gen-nextra/internal/generator"
)

// TestGenerateFile runs the generator against a pre-compiled descriptor set.
// Generate the fixture with: make testdata
func TestGenerateFile(t *testing.T) {
	req := loadRequest(t, filepath.Join("..", "..", "testdata", "all.pb"))

	var genErr error

	opts := protogen.Options{}
	// protogen.Options.Run reads from os.Stdin; for tests we use NewPlugin directly.
	plugin, err := opts.New(req)
	if err != nil {
		t.Fatalf("create plugin: %v", err)
	}

	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}
		if err := generator.GenerateFile(plugin, f, generator.Config{TypePreviews: true, Examples: true}); err != nil {
			genErr = err
			break
		}
	}
	if genErr != nil {
		t.Fatalf("generate: %v", genErr)
	}

	resp := plugin.Response()
	if resp.GetError() != "" {
		t.Fatalf("plugin error: %s", resp.GetError())
	}

	update := os.Getenv("UPDATE_GOLDEN") == "1"

	for _, rf := range resp.File {
		golden := filepath.Join("..", "..", "testdata", "golden", rf.GetName())

		if update {
			if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
				t.Fatalf("mkdir: %v", err)
			}
			if err := os.WriteFile(golden, []byte(rf.GetContent()), 0o644); err != nil {
				t.Fatalf("write golden: %v", err)
			}
			t.Logf("updated golden: %s", golden)
			continue
		}

		want, err := os.ReadFile(golden)
		if err != nil {
			t.Fatalf("read golden %s (run UPDATE_GOLDEN=1 go test to create it): %v", golden, err)
		}
		if !bytes.Equal(want, []byte(rf.GetContent())) {
			t.Errorf("output mismatch for %s:\nwant:\n%s\ngot:\n%s", rf.GetName(), want, rf.GetContent())
		}
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
