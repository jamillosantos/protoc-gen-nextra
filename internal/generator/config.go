package generator

// Config holds the options passed to the plugin via --nextra_opt.
type Config struct {
	// SplitServices generates one MDX page per service instead of one per package.
	// Enable with: --nextra_opt=split_services=true
	SplitServices bool

	// TypePreviews controls whether cross-package type references render an
	// inline hover preview card. Disable with: --nextra_opt=disable_type_previews=true
	TypePreviews bool

	// Examples enables generation of grpcurl and Go usage examples for unary
	// methods. Disabled by default.
	// Enable with: --nextra_opt=examples=true
	Examples bool

	// GoModule is the Go module path prefix used to build import paths in Go
	// usage examples. Default: "<YOUR-GO-MODULE>".
	// Set with: --nextra_opt=go_module=github.com/org/repo/gen
	GoModule string

	// ServerAddr is the gRPC server address used in usage examples.
	// Default: "localhost:50051".
	// Set with: --nextra_opt=server_addr=api.example.com:443
	ServerAddr string
}
