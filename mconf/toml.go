package mconf

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
)

const errStr = "Only a struct or map can be marshaled to TOML"

func init() {
	registerReadParser([]string{"toml"}, []string{".toml"}, &TOMLParser{})
	registerWriteParser([]string{"toml"}, []string{".toml"}, &TOMLParser{})
}

// TOMLParser is a Parser implementation to handle toml files.
type TOMLParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *TOMLParser) FromBytes(byteData []byte) (interface{}, error) {
	var data interface{}
	if err := toml.Unmarshal(byteData, &data); err != nil {
		return data, fmt.Errorf("could not unmarshal data: %w", err)
	}
	return &BasicSingleDocument{
		Value: data,
	}, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *TOMLParser) ToBytes(value interface{}, options ...ReadWriteOption) ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := toml.NewEncoder(buf)

	colourise := false

	for _, o := range options {
		switch o.Key {
		case OptionIndent:
			if indent, ok := o.Value.(string); ok {
				enc.Indent = indent
			}
		case OptionColourise:
			if value, ok := o.Value.(bool); ok {
				colourise = value
			}
		}
	}

	switch d := value.(type) {
	case SingleDocument:
		if err := enc.Encode(d.Document()); err != nil {
			if err.Error() == errStr {
				buf.Write([]byte(fmt.Sprintf("%v\n", d.Document())))
			} else {
				return nil, err
			}
		}
	case MultiDocument:
		for _, dd := range d.Documents() {
			if err := enc.Encode(dd); err != nil {
				if err.Error() == errStr {
					buf.Write([]byte(fmt.Sprintf("%v\n", dd)))
				} else {
					return nil, err
				}
			}
		}
	default:
		if err := enc.Encode(d); err != nil {
			if err.Error() == errStr {
				buf.Write([]byte(fmt.Sprintf("%v\n", d)))
			} else {
				return nil, err
			}
		}
	}

	if colourise {
		if err := ColouriseBuffer(buf, "toml"); err != nil {
			return nil, fmt.Errorf("could not colourise output: %w", err)
		}
	}

	return buf.Bytes(), nil
}
