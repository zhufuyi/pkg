package mconf

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var yamlBytes = []byte(`name: Tom
numbers:
- 1
- 2
`)
var yamlMap = map[string]interface{}{
	"name": "Tom",
	"numbers": []interface{}{
		1,
		2,
	},
}

var yamlBytesMulti = []byte(`name: Tom
---
name: Jim
`)
var yamlMapMulti = []interface{}{
	map[string]interface{}{
		"name": "Tom",
	},
	map[string]interface{}{
		"name": "Jim",
	},
}

func TestYAMLParser_FromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&YAMLParser{}).FromBytes(yamlBytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &BasicSingleDocument{Value: yamlMap}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("ValidMultiDocument", func(t *testing.T) {
		got, err := (&YAMLParser{}).FromBytes(yamlBytesMulti)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &BasicMultiDocument{Values: yamlMapMulti}

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		_, err := (&YAMLParser{}).FromBytes([]byte(`{1:asd`))
		if err == nil || !strings.Contains(err.Error(), "could not unmarshal data") {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
	t.Run("Empty", func(t *testing.T) {
		got, err := (&YAMLParser{}).FromBytes([]byte(``))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(nil, got) {
			t.Errorf("expected %v, got %v", nil, got)
		}
	})
}

func TestYAMLParser_ToBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&YAMLParser{}).ToBytes(yamlMap)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytes) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytes, got)
		}
	})
	t.Run("ValidSingle", func(t *testing.T) {
		got, err := (&YAMLParser{}).ToBytes(&BasicSingleDocument{Value: yamlMap})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytes) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytes, got)
		}
	})
	t.Run("ValidSingleColourise", func(t *testing.T) {
		got, err := (&YAMLParser{}).ToBytes(&BasicSingleDocument{Value: yamlMap}, ColouriseOption(true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		expBuf, _ := Colourise(string(yamlBytes), "yaml")
		exp := expBuf.Bytes()
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("ValidMulti", func(t *testing.T) {
		got, err := (&YAMLParser{}).ToBytes(&BasicMultiDocument{Values: yamlMapMulti})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytesMulti) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytesMulti, got)
		}
	})
}

var yamlFileData = []byte(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-dm
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
        tier: frontend
      name: nginx-pod
    spec:
      containers:
      - image: nginx:1.11.2
        imagePullPolicy: IfNotPresent
        name: nginx
        ports:
        - containerPort: 80
          name: nginx
`)

func TestYAMLParser_ToBytes2(t *testing.T) {
	values, err := (&YAMLParser{}).FromBytes(yamlFileData)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%+v\n\n", values)

	val2, err := (&YAMLParser{}).ToBytes(values)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%s\n\n", string(val2))
}
