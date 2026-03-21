package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// goPackageAlias derives a short Go package alias from a proto file path.
// e.g. "rpc/greeter/v1/greeter.proto" → "greeterv1"
func goPackageAlias(protoFilePath string) string {
	dir := filepath.Dir(protoFilePath)
	parts := strings.Split(dir, "/")
	skip := map[string]bool{"rpc": true, "proto": true, "protos": true, "api": true, "shared": true}
	var keep []string
	for _, p := range parts {
		if !skip[p] && p != "." {
			keep = append(keep, p)
		}
	}
	if len(keep) >= 2 {
		return keep[len(keep)-2] + keep[len(keep)-1]
	}
	if len(keep) == 1 {
		return keep[0]
	}
	return strings.ReplaceAll(dir, "/", "")
}

// toCamelCase converts a snake_case proto field name to Go CamelCase.
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return b.String()
}

// jsonExampleValue returns a JSON placeholder value for a proto field.
func jsonExampleValue(f *protogen.Field) string {
	if f.Desc.IsList() {
		return "[]"
	}
	switch f.Desc.Kind() {
	case protoreflect.StringKind:
		return `"example"`
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return "0.0"
	case protoreflect.BytesKind:
		return `""`
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return "{}"
	default:
		return "0"
	}
}

// goExampleValue returns a Go placeholder value for a proto field.
func goExampleValue(f *protogen.Field) string {
	if f.Desc.IsList() {
		return "nil"
	}
	switch f.Desc.Kind() {
	case protoreflect.StringKind:
		return `"example"`
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.FloatKind:
		return "float32(0)"
	case protoreflect.DoubleKind:
		return "float64(0)"
	case protoreflect.BytesKind, protoreflect.MessageKind, protoreflect.GroupKind:
		return "nil"
	default:
		return "0"
	}
}

// buildGrpcurlExample generates a grpcurl command for a unary RPC method.
func buildGrpcurlExample(pkg, service, method string, m *protogen.Method, serverAddr string) string {
	if serverAddr == "" {
		serverAddr = "localhost:50051"
	}
	var parts []string
	for _, f := range m.Input.Fields {
		parts = append(parts, fmt.Sprintf("%q: %s", string(f.Desc.Name()), jsonExampleValue(f)))
	}

	var jsonBody string
	if len(parts) == 0 {
		jsonBody = "{}"
	} else {
		jsonBody = "{" + strings.Join(parts, ", ") + "}"
	}

	fullMethod := fmt.Sprintf("%s.%s/%s", pkg, service, method)
	return fmt.Sprintf("grpcurl -plaintext \\\n  -d '%s' \\\n  %s \\\n  %s", jsonBody, serverAddr, fullMethod)
}

// buildGoExample generates a Go client snippet for a unary RPC method.
func buildGoExample(f *protogen.File, svc *protogen.Service, m *protogen.Method, cfg Config) string {
	goModule := cfg.GoModule
	if goModule == "" {
		goModule = "<YOUR-GO-MODULE>"
	}
	serverAddr := cfg.ServerAddr
	if serverAddr == "" {
		serverAddr = "localhost:50051"
	}
	alias := goPackageAlias(string(f.Desc.Path()))
	importPath := goModule + "/" + filepath.Dir(string(f.Desc.Path()))
	serviceName := string(svc.Desc.Name())
	methodName := string(m.Desc.Name())
	inputName := string(m.Input.Desc.Name())

	var fieldLines []string
	for _, field := range m.Input.Fields {
		name := toCamelCase(string(field.Desc.Name()))
		val := goExampleValue(field)
		fieldLines = append(fieldLines, fmt.Sprintf("\t\t%s: %s,", name, val))
	}

	var reqLiteral string
	if len(fieldLines) == 0 {
		reqLiteral = fmt.Sprintf("&%s.%s{}", alias, inputName)
	} else {
		reqLiteral = fmt.Sprintf("&%s.%s{\n%s\n\t}", alias, inputName, strings.Join(fieldLines, "\n"))
	}

	return fmt.Sprintf(`import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	%s %q
)

func main() {
	conn, err := grpc.NewClient(%q,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := %s.New%sClient(conn)
	resp, err := client.%s(context.Background(), %s)
	if err != nil {
		log.Fatal(err)
	}
	_ = resp
}`,
		alias, importPath,
		serverAddr,
		alias, serviceName,
		methodName, reqLiteral,
	)
}
