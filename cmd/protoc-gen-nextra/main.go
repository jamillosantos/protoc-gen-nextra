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

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		cfg := generator.Config{SplitServices: *splitServices}
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			if err := generator.GenerateFile(gen, f, cfg); err != nil {
				return err
			}
		}
		return nil
	})
}
