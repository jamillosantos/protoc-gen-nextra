package generator

// Config holds the options passed to the plugin via --nextra_opt.
type Config struct {
	// SplitServices generates one MDX page per service instead of one per package.
	// Enable with: --nextra_opt=split_services=true
	SplitServices bool

	// TypePreviews controls whether cross-package type references render an
	// inline hover preview card. Disable with: --nextra_opt=disable_type_previews=true
	TypePreviews bool
}
