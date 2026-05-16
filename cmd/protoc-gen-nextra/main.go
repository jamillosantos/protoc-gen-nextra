package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/jamillosantos/protoc-gen-nextra/internal/generator"
)

func main() {
	var flags flag.FlagSet
	splitServices := flags.Bool("split_services", false, "generate one MDX page per service instead of one per package")
	disableTypePreviews := flags.Bool("disable_type_previews", false, "disable hover preview cards for cross-package type references")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		cfg := generator.Config{
			SplitServices: *splitServices,
			TypePreviews:  !*disableTypePreviews,
		}
		order, byDir := generator.GroupByDir(gen.Files)
		for _, dir := range order {
			if err := generator.GeneratePackage(gen, byDir[dir], cfg); err != nil {
				return err
			}
		}
		return nil
	})
}
