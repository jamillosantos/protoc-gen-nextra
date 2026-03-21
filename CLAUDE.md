# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build          # compile binary to bin/protoc-gen-nextra
make install        # install binary locally
make test           # generate testdata then run go test ./...
make lint           # golangci-lint run ./...
make testdata       # regenerate testdata/greeter.pb from proto (requires protoc)
make generate       # run the plugin against test protos → testdata/content/
make update-golden  # UPDATE_GOLDEN=1 go test ./... — update golden files after intentional changes
```

Run a single test:
```bash
go test ./internal/generator/... -run TestGenerateFile
```

## Architecture

This is a `protoc` compiler plugin. The binary is invoked by protoc/buf — it reads a `CodeGeneratorRequest` from stdin and writes a `CodeGeneratorResponse` to stdout.

**Flow:**
1. `cmd/protoc-gen-nextra/main.go` — bootstraps `protogen.Options{}.Run()`, iterates files with services, delegates to `generator.GenerateFile()`
2. `internal/generator/generator.go` — core logic: extracts proto metadata into `ServiceData`/`MethodData`/`FieldData` structs, then renders via template
3. `internal/generator/templates/service.tmpl` — Go template producing `.mdx` output; uses `[[ ]]` delimiters (not `{{ }}`) to avoid JSX conflicts
4. `internal/generator/fields.go` — maps proto field kinds to human-readable type strings
5. `internal/generator/funcmap.go` + `embed.go` — template helpers and `//go:embed` for single-binary distribution

**Output:** One `.mdx` file per proto service, written to `<proto_package_path>/<service-name-kebab>.mdx`. The MDX uses Nextra's `<Tabs>` component with request/response field tables and a streaming-type badge per method.

## Testing

Tests use golden file comparison. The test fixture is `testdata/proto/greeter.proto` compiled to `testdata/greeter.pb`. Test output is compared against `testdata/golden/`.

When changing template or generator logic, run `make update-golden` to accept new output, then verify the diff is intentional before committing.
