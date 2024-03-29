package mconf

import (
	"reflect"
	"testing"
)

var jsonBytes = []byte(`{
  "name": "Tom"
}
`)
var jsonMap = map[string]interface{}{
	"name": "Tom",
}

func TestJSONParser_FromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&JSONParser{}).FromBytes(jsonBytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &BasicSingleDocument{Value: jsonMap}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("ValidMultiDocument", func(t *testing.T) {
		got, err := (&JSONParser{}).FromBytes(jsonBytesMulti)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &BasicMultiDocument{
			Values: jsonMapMulti,
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", jsonMap, got)
		}
	})
	t.Run("ValidMultiDocumentMixed", func(t *testing.T) {
		got, err := (&JSONParser{}).FromBytes(jsonBytesMultiMixed)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &BasicMultiDocument{
			Values: jsonMapMultiMixed,
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", jsonMap, got)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		got, err := (&JSONParser{}).FromBytes([]byte(``))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(nil, got) {
			t.Errorf("expected %v, got %v", nil, got)
		}
	})
}

func TestJSONParser_FromBytes_Error(t *testing.T) {
	_, err := (&JSONParser{}).FromBytes(yamlBytes)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
}

func TestJSONParser_ToBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&JSONParser{}).ToBytes(jsonMap)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytes) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytes), string(got))
		}
	})

	t.Run("ValidSingle", func(t *testing.T) {
		got, err := (&JSONParser{}).ToBytes(&BasicSingleDocument{Value: jsonMap})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytes) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytes), string(got))
		}
	})

	t.Run("ValidSingleNoPrettyPrint", func(t *testing.T) {
		res, err := (&JSONParser{}).ToBytes(&BasicSingleDocument{Value: jsonMap}, PrettyPrintOption(false))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		got := string(res)
		exp := `{"name":"Tom"}
`
		if exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})

	t.Run("ValidSingleColourise", func(t *testing.T) {
		got, err := (&JSONParser{}).ToBytes(&BasicSingleDocument{Value: jsonMap}, ColouriseOption(true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		expBuf, _ := Colourise(`{
  "name": "Tom"
}
`, "json")
		exp := expBuf.Bytes()
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})

	t.Run("ValidSingleCustomIndent", func(t *testing.T) {
		res, err := (&JSONParser{}).ToBytes(&BasicSingleDocument{Value: jsonMap}, IndentOption("   "))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		got := string(res)
		exp := `{
   "name": "Tom"
}
`
		if exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})

	t.Run("ValidMulti", func(t *testing.T) {
		got, err := (&JSONParser{}).ToBytes(&BasicMultiDocument{Values: jsonMapMulti})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytesMulti) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytesMulti), string(got))
		}
	})

	t.Run("ValidMultiMixed", func(t *testing.T) {
		got, err := (&JSONParser{}).ToBytes(&BasicMultiDocument{Values: jsonMapMultiMixed})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytesMultiMixed) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytesMultiMixed), string(got))
		}
	})
}

var jsonBytesMulti = []byte(`{
  "name": "Tom"
}
{
  "name": "Ellis"
}
`)

var jsonMapMulti = []interface{}{
	map[string]interface{}{"name": "Tom"},
	map[string]interface{}{"name": "Ellis"},
}

var jsonBytesMultiMixed = []byte(`{
  "name": "Tom",
  "other": true
}
{
  "name": "Ellis"
}
`)

var jsonMapMultiMixed = []interface{}{
	map[string]interface{}{"name": "Tom", "other": true},
	map[string]interface{}{"name": "Ellis"},
}
