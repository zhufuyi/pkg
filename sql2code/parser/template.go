package parser

import (
	"sync"
	"text/template"
)

var (
	modelStructTmpl    *template.Template
	modelStructTmplRaw = `
{{- if .Comment -}}
// {{.TableName}} {{.Comment}}
{{end -}}
type {{.TableName}} struct {
{{- range .Fields}}
	{{.Name}} {{.GoType}} {{if .Tag}}` + "`{{.Tag}}`" + `{{end}}{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
{{if .NameFunc}}
// TableName table name
func (m *{{.TableName}}) TableName() string {
	return "{{.RawTableName}}"
}
{{end}}
`

	modelTmpl    *template.Template
	modelTmplRaw = `package {{.Package}}
{{if .ImportPath}}
import (
	{{- range .ImportPath}}
	"{{.}}"
	{{- end}}
)
{{- end}}
{{range .StructCode}}
{{.}}
{{end}}`

	updateFieldTmpl    *template.Template
	updateFieldTmplRaw = `
{{- range .Fields}}
	if table.{{.Name}} {{.ConditionZero}} {
		update["{{.ColName}}"] = table.{{.Name}}
	}
{{- end}}`

	handlerCreateStructTmpl    *template.Template
	handlerCreateStructTmplRaw = `
// Create{{.TableName}}Request create params
type Create{{.TableName}}Request struct {
// todo fill in the binding rules https://github.com/go-playground/validator
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerUpdateStructTmpl    *template.Template
	handlerUpdateStructTmplRaw = `
// Update{{.TableName}}ByIDRequest update params
type Update{{.TableName}}ByIDRequest struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}" binding:""` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}
`

	handlerDetailStructTmpl    *template.Template
	handlerDetailStructTmplRaw = `
// Get{{.TableName}}ByIDRespond respond detail
type Get{{.TableName}}ByIDRespond struct {
{{- range .Fields}}
	{{.Name}}  {{.GoType}} ` + "`" + `json:"{{.ColName}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{- end}}
}`

	modelJSONTmpl    *template.Template
	modelJSONTmplRaw = `{
{{- range .Fields}}
	"{{.ColName}}" {{.GoZero}}
{{- end}}
}
`

	tmplParseOnce sync.Once
)

func initTemplate() {
	tmplParseOnce.Do(func() {
		var err error
		modelStructTmpl, err = template.New("goStruct").Parse(modelStructTmplRaw)
		if err != nil {
			panic(err)
		}
		modelTmpl, err = template.New("goFile").Parse(modelTmplRaw)
		if err != nil {
			panic(err)
		}
		updateFieldTmpl, err = template.New("goUpdateField").Parse(updateFieldTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerCreateStructTmpl, err = template.New("goPostStruct").Parse(handlerCreateStructTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerUpdateStructTmpl, err = template.New("goPutStruct").Parse(handlerUpdateStructTmplRaw)
		if err != nil {
			panic(err)
		}
		handlerDetailStructTmpl, err = template.New("goGetStruct").Parse(handlerDetailStructTmplRaw)
		if err != nil {
			panic(err)
		}
		modelJSONTmpl, err = template.New("modelJSON").Parse(modelJSONTmplRaw)
		if err != nil {
			panic(err)
		}
	})
}
