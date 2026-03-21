package generator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldsToJSON converts a list of ErrorDetailFieldData into an indented JSON
// string. Field names support dotted paths and [N] array indexing, e.g.:
//
//	field_violations[0].field       → "email"
//	field_violations[0].description → "must be valid"
//
// produces:
//
//	{
//	  "field_violations": [
//	    { "field": "email", "description": "must be valid" }
//	  ]
//	}
func FieldsToJSON(fields []ErrorDetailFieldData) string {
	if len(fields) == 0 {
		return ""
	}
	var root interface{} = make(map[string]interface{})
	for _, f := range fields {
		root = setNested(root, parsePath(f.Name), f.Value)
	}
	b, _ := json.MarshalIndent(root, "", "  ")
	return string(b)
}

// segment is one step in a field path.
type segment struct {
	key     string
	isIndex bool
	index   int
}

// parsePath splits a path like "field_violations[0].field" into segments.
func parsePath(path string) []segment {
	var segs []segment
	for _, part := range strings.Split(path, ".") {
		if i := strings.Index(part, "["); i >= 0 {
			if key := part[:i]; key != "" {
				segs = append(segs, segment{key: key})
			}
			var idx int
			fmt.Sscanf(part[i:], "[%d]", &idx)
			segs = append(segs, segment{isIndex: true, index: idx})
		} else if part != "" {
			segs = append(segs, segment{key: part})
		}
	}
	return segs
}

// setNested recursively navigates or creates the tree and sets the leaf value.
func setNested(node interface{}, segs []segment, value string) interface{} {
	if len(segs) == 0 {
		return value
	}
	seg := segs[0]
	rest := segs[1:]

	if seg.isIndex {
		var arr []interface{}
		if a, ok := node.([]interface{}); ok {
			arr = a
		}
		for len(arr) <= seg.index {
			arr = append(arr, nil)
		}
		arr[seg.index] = setNested(arr[seg.index], rest, value)
		return arr
	}

	m, ok := node.(map[string]interface{})
	if !ok || m == nil {
		m = make(map[string]interface{})
	}
	m[seg.key] = setNested(m[seg.key], rest, value)
	return m
}
