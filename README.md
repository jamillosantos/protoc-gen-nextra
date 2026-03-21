# protoc-gen-nextra

A `protoc` plugin that generates [Nextra](https://nextra.site) MDX documentation pages from gRPC service definitions.

![Screenshot of generated documentation](_assets/doc_screenshot.png)

Each `.proto` file with services produces a `.mdx` page with:

- **RPC type badges** — `UNARY`, `SERVER STREAM`, `CLIENT STREAM`, `BIDI STREAM` — color-coded and visually distinct
- **Request/response field tables** with type, `optional`/`repeated` pills, and inline doc comments
- **Cross-package type links** with hover preview cards showing the referenced type's fields
- **Error documentation** — per-method error codes with descriptions, detail types, and example JSON payloads
- **Usage examples** — grpcurl and Go code tabs per unary method (opt-in)
- Source comments from your `.proto` file carried through as descriptions

## Requirements

- Go 1.21+
- `protoc` — e.g. `brew install protobuf`

## Installation

```sh
go install github.com/jamillosantos/protoc-gen-nextra/cmd/protoc-gen-nextra@latest
```

## Usage

### With Buf (recommended)

Install the binary, then add the plugin to your `buf.gen.yaml`:

```yaml
version: v2
plugins:
  - local: protoc-gen-nextra
    out: docs/content
```

Then run:

```sh
buf generate
```

### With protoc

```sh
protoc \
  --nextra_out=./docs/content \
  -I ./proto \
  ./proto/**/*.proto
```

`--nextra_out` should point to the `content/` directory of your Nextra project. The plugin produces one `.mdx` file per proto directory, named after the directory (e.g. `greeter/v1/greeter.proto` → `greeter/v1.mdx`).

### Example

Given this proto at `greeter/v1/greeter.proto`:

```proto
syntax = "proto3";
package greeter.v1;

// Greeter provides greeting functionality.
service Greeter {
  // SayHello sends a greeting.
  rpc SayHello (HelloRequest) returns (HelloReply);

  // SayHelloStream streams greetings back to the client.
  rpc SayHelloStream (HelloRequest) returns (stream HelloReply);
}
```

The plugin generates `greeter/v1.mdx` with:

- A `UNARY` badge for `SayHello`
- A `SERVER STREAM` badge for `SayHelloStream`
- Tabbed request/response field tables for each method

## Options

Options are passed via `--nextra_opt` (protoc) or the `opt` key (buf).

| Option | Default | Description |
|---|---|---|
| `split_services` | `false` | Generate one page per service instead of one page per proto file. Each service is written to `<proto_dir>/<service-name>.mdx`. |
| `disable_type_previews` | `false` | Disable the hover preview cards for cross-package type references. |
| `examples` | `false` | Generate grpcurl and Go usage examples for every unary method. |
| `go_module` | `<YOUR-GO-MODULE>` | Go module path used to build import paths in Go usage examples. |
| `server_addr` | `localhost:50051` | Default gRPC server address used in usage examples. Can be overridden per service with the `(nextra.server_addr)` proto option. |

### With Buf

```yaml
version: v2
plugins:
  - local: protoc-gen-nextra
    out: docs/content
    opt:
      - examples=true
      - go_module=github.com/org/repo/gen
      - server_addr=api.example.com:443
```

### With protoc

```sh
protoc \
  --nextra_opt=examples=true \
  --nextra_opt=go_module=github.com/org/repo/gen \
  --nextra_opt=server_addr=api.example.com:443 \
  --nextra_out=./docs/content \
  -I ./proto \
  ./proto/**/*.proto
```

## Proto options

Import `nextra/options.proto` to annotate your services and methods.

### `(nextra.server_addr)` — per-service server address

Overrides the `server_addr` plugin option for a specific service:

```proto
import "nextra/options.proto";

service MyService {
  option (nextra.server_addr) = "grpc.example.com:443";

  rpc GetFoo (GetFooRequest) returns (GetFooResponse);
}
```

### `(nextra.method_errors)` — error documentation

Documents the errors a method may return. Each entry renders as a heading with description, an optional detail type, and a JSON example payload:

```proto
import "nextra/options.proto";

rpc CreateFoo (CreateFooRequest) returns (CreateFooResponse) {
  option (nextra.method_errors) = {
    errors: [
      {
        code: "NOT_FOUND",
        description: "the requested resource does not exist."
      },
      {
        code: "INVALID_ARGUMENT",
        description: "one or more fields failed validation.",
        detail_type: "google.rpc.BadRequest",
        fields: [
          { name: "field_violations[0].field",       value: "email" },
          { name: "field_violations[0].description", value: "must be a valid email address." }
        ]
      }
    ]
  };
}
```

Field names support dotted paths and `[N]` array indexing. The generator reconstructs the full nested JSON structure — the example above produces:

```json
{
  "field_violations": [
    {
      "description": "must be a valid email address.",
      "field": "email"
    }
  ]
}
```

If `detail_type` refers to a message defined in the same proto compilation, it is rendered as a hyperlink with a hover preview card.

## Nextra setup

Your Nextra project needs nextra v4+:

```sh
bun add nextra nextra-theme-docs next react react-dom
```

Your `next.config.mjs`:

```js
import nextra from 'nextra'

const withNextra = nextra({ contentDirBasePath: '/' })
export default withNextra()
```

See [`testdata/`](./testdata) for a working example — run `bun run dev` inside it to preview the output locally.

## Badge colours

| Streaming type | Badge |
|---|---|
| Unary | ![blue](https://img.shields.io/badge/UNARY-0070f3?style=flat) |
| Server streaming | ![purple](https://img.shields.io/badge/SERVER%20STREAM-7928ca?style=flat) |
| Client streaming | ![orange](https://img.shields.io/badge/CLIENT%20STREAM-f5a623?style=flat) |
| Bidirectional | ![red](https://img.shields.io/badge/BIDI%20STREAM-ee0000?style=flat) |

## Development

```sh
# Build the binary
make build

# Generate MDX from the test proto and run tests
make test

# Preview in Nextra
cd testdata && bun run dev

# Regenerate golden test files after template changes
make update-golden
```

## License

MIT
