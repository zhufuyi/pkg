package mconf

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp"
)

func TestTrimKey(t *testing.T) {
	k := "   --abc "
	fmt.Printf("(%s)\n", trimKey(k))
}

func TestGetMatchArgs(t *testing.T) {
	s := "--requirepass 123456"
	argsMap := map[string]string{
		"requirepass": "1q2w3e4r",
		"port":        "36379",
	}

	records, newKV := getMatchArgs(s, argsMap)
	pp.Println(newKV, records)
}

func TestAddOrReplaceArgs(t *testing.T) {
	args := [][]string{
		{
			"-c",
		},
		{
			"--requirepass=123456",
			"--port=6379",
			"-host=127.0.0.1",
		},
	}
	argsMap := map[string]string{
		"requirepass": "1q2w3e4r",
		"port":        "36379",
		"host":        "127.0.0.1",
		"flag":        "",
	}

	records := addOrReplaceArgs(args, argsMap)

	pp.Println(records)
	pp.Println(args)
}

func TestStr2Map(t *testing.T) {
	str := "--a=1, -b=2, -c"
	fmt.Println(str2Map(str))
}

// ---------------------------------------------------------------------------------------

var yamlData = []byte(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-dm
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      name: nginx-pod
      labels:
        app: nginx
        tier: frontend
    spec:
      containers:
        - name: nginx
          image: nginx:1.15.2
          imagePullPolicy: IfNotPresent
          ports:
          - name: nginx
            containerPort: 80
`)

var jsonData = []byte(`
{
  "annotations": {
    "list": [
      {
        "$$hashKey": "object:75",
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "gnetId": 13639,
  "id": 5,
  "iteration": 1644922431014,
  "panels": [
    {
      "datasource": "Loki",
      "fieldConfig": {
        "defaults": {
          "custom": {}
        },
        "overrides": []
      },
      "gridPos": {
        "h": 25,
        "w": 24,
        "x": 0,
        "y": 3
      },
      "id": 2,
      "maxDataPoints": "",
      "options": {
        "showLabels": false,
        "showTime": true,
        "sortOrder": "Descending",
        "wrapLogMessage": false
      },
      "targets": [
        {
          "expr": "{job=\"$app\"} |= \"$search\" | logfmt",
          "hide": false,
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "timeFrom": null,
      "timeShift": null,
      "title": "",
      "transparent": true,
      "type": "logs"
    }
  ],
  "title": "Logs / App",
  "uid": "liz0yRCZz",
  "version": 1
}
`)
