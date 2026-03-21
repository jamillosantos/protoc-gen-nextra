package generator

import "google.golang.org/protobuf/compiler/protogen"

// Config holds the options passed to the plugin via --nextra_opt.
type Config struct {
	// SplitServices generates one MDX page per service instead of one per package.
	// Enable with: --nextra_opt=split_services=true
	SplitServices bool

	// TypePreviews controls whether cross-package type references render an
	// inline hover preview card. Disable with: --nextra_opt=disable_type_previews=true
	TypePreviews bool

	// typeIndex is an internal map from fully-qualified proto type name to MDX
	// page link (with anchor). Built once from all files in the compilation.
	typeIndex map[string]string

	// typeMessageIndex maps fully-qualified proto type names to their message
	// descriptors, used to build TypePreview data for error detail types.
	typeMessageIndex map[string]*protogen.Message
}
