package main

import (
	"encoding/json"
	"io"
	"strings"
)

func stringIsNilOrEmpty(s *string) bool {
	return s == nil || stringIsEmpty(*s)
}
func stringIsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func jsonEncode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func jsonDecode(w io.Reader, v any) error {
	return json.NewDecoder(w).Decode(v)
}
