package main

import (
	"encoding/json"
	"io"
)

func jsonEncode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func jsonDecode(w io.Reader, v any) error {
	return json.NewDecoder(w).Decode(v)
}
