// TEMPORARY AUTOGENERATED FILE: easyjson stub code to make the package
// compilable during generation.

package  jsontmp

import (
  "github.com/mailru/easyjson/jwriter"
  "github.com/mailru/easyjson/jlexer"
)

func ( Item ) MarshalJSON() ([]byte, error) { return nil, nil }
func (* Item ) UnmarshalJSON([]byte) error { return nil }
func ( Item ) MarshalEasyJSON(w *jwriter.Writer) {}
func (* Item ) UnmarshalEasyJSON(l *jlexer.Lexer) {}

type EasyJSON_exporter_Item *Item

func ( ItemList ) MarshalJSON() ([]byte, error) { return nil, nil }
func (* ItemList ) UnmarshalJSON([]byte) error { return nil }
func ( ItemList ) MarshalEasyJSON(w *jwriter.Writer) {}
func (* ItemList ) UnmarshalEasyJSON(l *jlexer.Lexer) {}

type EasyJSON_exporter_ItemList *ItemList
