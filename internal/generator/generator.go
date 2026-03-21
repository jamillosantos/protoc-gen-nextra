package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

//go:generate go run github.com/jamillosantos/protoc-gen-nextra/internal/embedgen

var tmpl *template.Template

func init() {
	// Use [[ ]] delimiters so JSX style={{...}} doesn't conflict with Go templates.
	tmpl = template.Must(template.New("").Delims("[[", "]]").Funcs(funcMap).ParseFS(templateFS, "templates/*.tmpl"))
}

// GenerateFile generates Nextra MDX documentation for a single .proto file.
// When cfg.SplitServices is true, each service gets its own page under <proto_dir>/<service>.mdx.
// Otherwise all services are combined into a single <proto_dir>.mdx page.
func GenerateFile(gen *protogen.Plugin, f *protogen.File, cfg Config) error {
	if len(f.Services) == 0 {
		return nil
	}

	if cfg.SplitServices {
		return generateSplit(gen, f)
	}
	return generateCombined(gen, f)
}

// generateCombined writes one MDX file containing all services in the proto file.
func generateCombined(gen *protogen.Plugin, f *protogen.File) error {
	// Output path: <proto_directory>.mdx (e.g. greeter/v1/greeter.proto → greeter/v1.mdx)
	outPath := filepath.Dir(f.Desc.Path()) + ".mdx"
	return renderPage(gen, f, outPath, buildPackageData(f, outPath))
}

// generateSplit writes one MDX file per service inside <proto_directory>/<service>.mdx.
func generateSplit(gen *protogen.Plugin, f *protogen.File) error {
	dir := filepath.Dir(f.Desc.Path())
	for _, svc := range f.Services {
		outPath := filepath.Join(dir, snakeit(string(svc.Desc.Name()))+".mdx")
		data := buildSingleServiceData(f, svc, outPath)
		if err := renderPage(gen, f, outPath, data); err != nil {
			return err
		}
	}
	return nil
}

func renderPage(gen *protogen.Plugin, f *protogen.File, outPath string, data PackageData) error {
	g := gen.NewGeneratedFile(outPath, f.GoImportPath)
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "service.tmpl", data); err != nil {
		return fmt.Errorf("executing template for %s: %w", outPath, err)
	}
	g.P(buf.String())
	return nil
}

// PackageData holds all template data for a package page.
type PackageData struct {
	Title                string
	PackageName          string
	ShowServiceHeadings  bool
	Services             []ServiceData
}

// ServiceData holds all template data for a single service.
type ServiceData struct {
	ServiceName string
	Description string
	Methods     []MethodData
}

// MethodData holds template data for a single RPC method.
type MethodData struct {
	Name            string
	Description     string
	RequestType     string
	ResponseType    string
	ClientStreaming  bool
	ServerStreaming  bool
	RequestFields   []FieldData
	ResponseFields  []FieldData
}

// StreamingType returns a short label for the RPC streaming mode.
func (m MethodData) StreamingType() string {
	switch {
	case m.ClientStreaming && m.ServerStreaming:
		return "BIDI STREAM"
	case m.ClientStreaming:
		return "CLIENT STREAM"
	case m.ServerStreaming:
		return "SERVER STREAM"
	default:
		return "UNARY"
	}
}

// BadgeColor returns a hex color for the streaming type badge.
func (m MethodData) BadgeColor() string {
	switch {
	case m.ClientStreaming && m.ServerStreaming:
		return "#e00"
	case m.ClientStreaming:
		return "#f5a623"
	case m.ServerStreaming:
		return "#7928ca"
	default:
		return "#0070f3"
	}
}

// FieldData holds template data for a message field.
type FieldData struct {
	Name        string
	Type        string
	Description string
	Optional    bool
	Repeated    bool
}

func buildPackageData(f *protogen.File, outPath string) PackageData {
	data := PackageData{
		Title:               filepath.Base(strings.TrimSuffix(outPath, ".mdx")),
		PackageName:         string(f.Desc.Package()),
		ShowServiceHeadings: len(f.Services) > 1,
	}
	for _, svc := range f.Services {
		data.Services = append(data.Services, buildServiceData(svc))
	}
	return data
}

func buildSingleServiceData(f *protogen.File, svc *protogen.Service, outPath string) PackageData {
	return PackageData{
		Title:               filepath.Base(strings.TrimSuffix(outPath, ".mdx")),
		PackageName:         string(f.Desc.Package()),
		ShowServiceHeadings: false,
		Services:            []ServiceData{buildServiceData(svc)},
	}
}

func buildServiceData(svc *protogen.Service) ServiceData {
	sd := ServiceData{
		ServiceName: string(svc.Desc.Name()),
		Description: commentString(svc.Comments),
	}
	for _, m := range svc.Methods {
		sd.Methods = append(sd.Methods, MethodData{
			Name:           string(m.Desc.Name()),
			Description:    commentString(m.Comments),
			RequestType:    string(m.Input.Desc.Name()),
			ResponseType:   string(m.Output.Desc.Name()),
			ClientStreaming: m.Desc.IsStreamingClient(),
			ServerStreaming: m.Desc.IsStreamingServer(),
			RequestFields:  buildFields(m.Input),
			ResponseFields: buildFields(m.Output),
		})
	}
	return sd
}

func buildFields(msg *protogen.Message) []FieldData {
	var fields []FieldData
	for _, f := range msg.Fields {
		fields = append(fields, FieldData{
			Name:        string(f.Desc.Name()),
			Type:        fieldTypeName(f),
			Description: commentString(f.Comments),
			Repeated:    f.Desc.IsList(),
			Optional:    f.Desc.HasOptionalKeyword(),
		})
	}
	return fields
}

func commentString(loc protogen.CommentSet) string {
	s := strings.TrimSpace(loc.Leading.String())
	if s == "" {
		s = strings.TrimSpace(loc.Trailing.String())
	}
	// protogen returns comments as "// text\n" per line. Strip the "// " prefix.
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		l = strings.TrimPrefix(l, "// ")
		l = strings.TrimPrefix(l, "//")
		lines = append(lines, l)
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func snakeit(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
