package codegen

import (
	"bytes"
	"text/template"
)

const factoryTemplate = `
{{ if .Append | not -}}
package factory

import (
    gofactory "github.com/vx416/gogo-factory"
	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/genutil"
)
{{ end -}}

{{ range .Structs}}
type {{.Name}}Factory struct {
	*gofactory.Factory
}

var {{.Name}} = &{{.Name}}Factory{gofactory.New(
    &{{.Package}}.{{.Name}}{},
    {{- range .Fields}}
    attr.{{ .Type | GetTypeName }}("{{.Name}}", {{ .Type | GetGenFunc }}),
    {{- end}}
)}
{{end}}
`

func GetTempalte(append bool, fileMeta FileMeta) (string, error) {
	var buf = bytes.NewBufferString("")
	t, err := template.New("").Funcs(map[string]interface{}{
		"GetTypeName": GetTypeName,
		"GetGenFunc":  GetGetFunc,
	}).Parse(factoryTemplate)
	if err != nil {
		return "", err
	}

	data := convertMetaToTemplateData(fileMeta)
	data["Append"] = append
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func convertMetaToTemplateData(fileMeta FileMeta) map[string]interface{} {
	res := make(map[string]interface{})
	structsData := make([]map[string]interface{}, 0, len(fileMeta.Structs))

	for _, st := range fileMeta.Structs {
		structData := map[string]interface{}{
			"Name":    st.Name,
			"Package": fileMeta.Package,
		}
		fields := make([]map[string]interface{}, 0, len(st.Fields))
		for _, field := range st.Fields {
			fields = append(fields, map[string]interface{}{
				"Name": field.Name,
				"Type": field.Type,
			})
		}
		structData["Fields"] = fields
		structsData = append(structsData, structData)
	}
	res["Structs"] = structsData
	return res
}

func GetTypeName(in string) string {
	switch in {
	case "string", "Decimal", "String", "NullString":
		return "Str"
	case "int", "int8", "int32", "int64", "Int", "NullInt64", "NullInt32":
		return "Int"
	case "uint", "uint8", "uint32", "uint64", "Uint":
		return "Uint"
	case "float", "float32", "float64", "Float", "NullFloat64":
		return "Float"
	case "Time", "NullTime":
		return "Time"
	case "byte":
		return "Bytes"
	case "bool", "Bool", "NullBool":
		return "Bool"
	default:
		return "Attr"
	}
}

func GetGetFunc(in string) string {
	switch in {
	case "string", "Decimal", "String", "NullString":
		return `genutil.RandName(3)`
	case "int", "int8", "int32", "int64", "Int", "NullInt64", "NullInt32":
		return `genutil.SeqInt(1, 1)`
	case "uint", "uint8", "uint32", "uint64", "Uint":
		return `genutil.SeqUint(1, 1)`
	case "Time", "NullTime":
		return `genutil.Now(time.UTC)`
	case "float", "float32", "float64", "Float", "NullFloat64":
		return `genutil.RandFloat(0, 10)`
	case "byte":
		return `genutil.FixBytes([]byte("test"))`
	case "bool", "Bool", "NullBool":
		return `genutil.RandBool(1)`
	default:
		return `genutil.FixInterface(nil)`
	}
}
