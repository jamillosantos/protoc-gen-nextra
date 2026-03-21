# protoc-gen-nextra

A `protoc` plugin that generates [Nextra](https://nextra.site) MDX documentation pages from gRPC `.proto` files.

![Screenshot of generated documentation](_assets/doc_screenshot.png)

## What it does

`protoc-gen-nextra` reads your `.proto` files and writes `.mdx` files into a directory of your choice — one file per proto package. That's it.

Each generated page includes:

- **RPC type badges** — `UNARY`, `SERVER STREAM`, `CLIENT STREAM`, `BIDI STREAM`
- **Request/response field tables** with types, `optional`/`repeated` pills, and inline doc comments from your proto
- **Cross-package type links** — fields that reference types from other packages link to their page and show a hover preview card
- **Error documentation** — per-method error codes with descriptions, detail types, and example JSON payloads (via proto annotations)
- **Usage examples** — grpcurl and Go code tabs per unary method (opt-in)

## What it does not do

`protoc-gen-nextra` **does not** set up Nextra, create a Next.js project, configure routing, or manage your documentation site. It only generates `.mdx` files.

You are responsible for:

- Creating and configuring your Nextra project
- Pointing the plugin output at your `content/` directory
- Adding the `TypePreview` component to your project if you want hover previews (see [TypePreview setup](#typepreview-setup))

Refer to the [Nextra documentation](https://nextra.site) for how to set up a Nextra project.

## Getting started

### 1. Install

```sh
go install github.com/jamillosantos/protoc-gen-nextra/cmd/protoc-gen-nextra@latest
```

Requirements: Go 1.21+, `protoc` (e.g. `brew install protobuf`).

### 2. Point the plugin at your content directory

#### With Buf (recommended)

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

#### With protoc

```sh
protoc \
  --nextra_out=./docs/content \
  -I ./proto \
  ./proto/**/*.proto
```

### 3. Run your Nextra dev server

```sh
cd docs && bun run dev   # or npm run dev
```

The generated files are picked up automatically — no manual page registration needed as long as your Nextra app uses a catch-all MDX route.

### Output format

The plugin produces one `.mdx` file per proto directory, named after the directory:

```
proto/greeter/v1/greeter.proto    →  content/greeter/v1.gen.mdx
proto/notifier/v1/notifier.proto  →  content/notifier/v1.gen.mdx
```

Generated files use the `.gen.mdx` extension so they are easily distinguishable from hand-written MDX files in your content directory.

## TypePreview setup

When `disable_type_previews` is not set, the generated MDX uses a `<TypePreview>` component to render hover cards on cross-package type links. You need to provide this component in your project and register it in `mdx-components.tsx`.

Copy [`docs/components/TypePreview.tsx`](./docs/components/TypePreview.tsx) into your project and register it:

```tsx
// mdx-components.tsx
import { TypePreview } from './components/TypePreview'
import type { MDXComponents } from 'mdx/types'

export function useMDXComponents(components: MDXComponents): MDXComponents {
  return { TypePreview, ...components }
}
```

See [`docs/`](./docs) for a complete working example.

## Options

Options are passed via `--nextra_opt` (protoc) or the `opt` key (buf).

| Option | Default | Description |
|---|---|---|
| `split_services` | `false` | Generate one page per service instead of one per proto file. Each service is written to `<proto_dir>/<service-name>.mdx`. |
| `disable_type_previews` | `false` | Disable hover preview cards for cross-package type references. |

### With Buf

```yaml
version: v2
plugins:
  - local: protoc-gen-nextra
    out: docs/content
    opt:
      - disable_type_previews=true
```

### With protoc

```sh
protoc \
  --nextra_opt=disable_type_previews=true \
  --nextra_out=./docs/content \
  -I ./proto \
  ./proto/**/*.proto
```

## Proto options

Import `nextra/options.proto` to annotate your services and methods with extra documentation.

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

Field names support dotted paths and `[N]` array indexing. The generator reconstructs the full nested JSON — the example above produces:

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

## Badge colours

| Streaming type | Badge |
|---|---|
| Unary | ![blue](https://img.shields.io/badge/UNARY-0070f3?style=flat) |
| Server streaming | ![purple](https://img.shields.io/badge/SERVER%20STREAM-7928ca?style=flat) |
| Client streaming | ![orange](https://img.shields.io/badge/CLIENT%20STREAM-f5a623?style=flat) |
| Bidirectional | ![red](https://img.shields.io/badge/BIDI%20STREAM-ee0000?style=flat) |

## Development

The [`docs/`](./docs) directory is a working Nextra app used to preview plugin output. Proto sources live in [`testdata/proto/`](./testdata/proto).

```sh
make generate        # generate MDX from testdata/proto/ into docs/content/
cd docs && bun run dev  # preview in Nextra

make test            # run tests
make update-golden   # update golden files after intentional template changes
make build           # compile binary to bin/protoc-gen-nextra
```

Directory layout:

```
testdata/
└── proto/                  # source .proto files (input to the plugin)

docs/                       # Nextra preview app
├── content/                # ← plugin writes .mdx files here
│   └── index.mdx           # hand-written landing page
├── app/
│   ├── layout.tsx          # Nextra Layout wrapper
│   └── [[...mdxPath]]/     # catch-all route serving MDX pages
├── components/
│   └── TypePreview.tsx     # hover preview component
├── mdx-components.tsx
└── next.config.mjs
```

## License

MIT
