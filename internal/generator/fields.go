package generator

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var scalarNames = map[protoreflect.Kind]string{
	protoreflect.BoolKind:    "bool",
	protoreflect.Int32Kind:   "int32",
	protoreflect.Sint32Kind:  "sint32",
	protoreflect.Uint32Kind:  "uint32",
	protoreflect.Int64Kind:   "int64",
	protoreflect.Sint64Kind:  "sint64",
	protoreflect.Uint64Kind:  "uint64",
	protoreflect.Sfixed32Kind: "sfixed32",
	protoreflect.Fixed32Kind:  "fixed32",
	protoreflect.FloatKind:    "float",
	protoreflect.Sfixed64Kind: "sfixed64",
	protoreflect.Fixed64Kind:  "fixed64",
	protoreflect.DoubleKind:   "double",
	protoreflect.StringKind:   "string",
	protoreflect.BytesKind:    "bytes",
}

func fieldTypeName(f *protogen.Field) string {
	if name, ok := scalarNames[f.Desc.Kind()]; ok {
		return name
	}
	switch f.Desc.Kind() {
	case protoreflect.EnumKind:
		return string(f.Enum.Desc.Name())
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return string(f.Message.Desc.Name())
	}
	return "unknown"
}
