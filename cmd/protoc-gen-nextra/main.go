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
	examples := flags.Bool("examples", false, "generate grpcurl and Go usage examples for unary methods")
	goModule := flags.String("go_module", "<YOUR-GO-MODULE>", "Go module path prefix for generated import paths in usage examples")
	serverAddr := flags.String("server_addr", "localhost:50051", "gRPC server address used in usage examples")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		cfg := generator.Config{
			SplitServices: *splitServices,
			TypePreviews:  !*disableTypePreviews,
			Examples:      *examples,
			GoModule:      *goModule,
			ServerAddr:    *serverAddr,
		}
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
