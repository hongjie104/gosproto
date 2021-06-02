package main

import (
	"bytes"
	"fmt"
	"go/token"

	"github.com/davyxu/gosproto/meta"
)

const goCodeTemplate = `package {{.PackageName}}

import (
	{{if gt (.Enums|len) 0}}sproto "github.com/hongjie104/gosproto"{{end}}
	"github.com/hongjie104/leaf/network/sproto"
)

{{range $a, $enumobj := .Enums}}
type {{.Name}} int32
const (	{{range .StFields}}
	{{$enumobj.Name}}_{{.Name}} {{.GoEnumTypeName}} = {{.TagNumber}} {{end}}
)

var (
{{$enumobj.Name}}MapperValueByName = map[string]int32{ {{range .StFields}}
	"{{.Name}}": {{.TagNumber}}, {{end}}
}

{{$enumobj.Name}}MapperNameByValue = map[int32]string{ {{range .StFields}}
	{{.TagNumber}}: "{{.Name}}" , {{end}}
}
)

func (p {{$enumobj.Name}}) String() string {
	return sproto.EnumName({{$enumobj.Name}}MapperNameByValue, int32(p))
}
{{end}}

{{range $a, $stobj := .Structs}}
// {{.Name}} {{.Name}}
type {{.Name}} struct{
	{{range .StFields}}
		{{.GoFieldName}} {{.GoTypeName}} {{.GoTags}}
	{{end}}
}

{{range .StFields}}{{if .IsExtendType}}
func (p *{{$stobj.Name}}) {{.GoExtendFieldGetterName}}() {{.GoExtendFieldGetterType}} {
	{{.GoExtendFieldGetter}}
}
func (p *{{$stobj.Name}}) Set{{.GoExtendFieldGetterName}}(v {{.GoExtendFieldGetterType}})  {
	{{.GoExtendFieldSetter}}
}
{{end}} {{end}}
{{end}}

{{if .EnumValueGroup}}
// ResultToString ResultToString
func ResultToString(result int32) string {
	switch( result ) {
		case 0: return "OK";
	{{range $a, $enumObj := .Enums}} {{if .IsResultEnum}} {{range .Fields}} {{if ne .TagNumber 0}}
		case {{.TagNumber}}: return "{{$enumObj.Name}}.{{.Name}}"; {{end}} {{end}} {{end}} {{end}}
	}

	return fmt.Sprintf("result: %d", result)
}
{{end}}

// MD5 MD5
const MD5 = "{{.MD5}}"

// Processor Processor
var Processor = sproto.NewProcessor()

func init() {
	Processor.SetByteOrder(true)
	{{range .Structs}}
	Processor.Register(&{{.Name}}{}){{end}}
}

`

func (fm *fieldModel) GoEnumTypeName() string {

	if fm.st.EnumValueIgnoreType {
		return ""
	}

	return fm.st.Name
}

func (fm *fieldModel) GoExtendFieldGetterName() string {
	pname := publicFieldName(fm.Name)

	if token.Lookup(pname).IsKeyword() {
		return pname + "_"
	}

	return pname
}

func (fm *fieldModel) GoExtendFieldGetter() string {
	switch fm.Type {
	case meta.FieldType_Float32,
		meta.FieldType_Float64:
		return fmt.Sprintf("return %s(self.%s) * %f", fm.GoExtendFieldGetterType(), fm.GoFieldName(), 1.0/float32(fm.ExtendTypePrecision()))
	}
	return "unknown extend type:" + fm.Type.String()
}

func (fm *fieldModel) GoExtendFieldSetter() string {
	switch fm.Type {
	case meta.FieldType_Float32:
		return fmt.Sprintf("self.%s = int32(v* %d)", fm.GoFieldName(), fm.ExtendTypePrecision())
	case meta.FieldType_Float64:
		return fmt.Sprintf("self.%s = int64(v* %d)", fm.GoFieldName(), fm.ExtendTypePrecision())
	}
	return "unknown extend type:" + fm.Type.String()
}

func (fm *fieldModel) GoExtendFieldGetterType() string {
	var b bytes.Buffer
	if fm.Repeatd {
		b.WriteString("[]")
	}
	// 字段类型映射go的类型
	switch fm.Type {
	case meta.FieldType_Float32:
		b.WriteString("float32")
	case meta.FieldType_Float64:
		b.WriteString("float64")
	default:
		b.WriteString("unknown extend type:" + fm.Type.String())
	}
	return b.String()
}

func (fm *fieldModel) GoFieldName() string {
	var pname string
	// 扩展类型不能直接访问
	if fm.IsExtendType() {
		pname = "Extend_" + fm.Name
	} else {
		pname = publicFieldName(fm.Name)
	}

	// 碰到关键字在尾部加_
	if token.Lookup(pname).IsKeyword() {
		return pname + "_"
	}

	return pname
}

func (fm *fieldModel) GoTypeName() string {

	var b bytes.Buffer
	if fm.Repeatd {
		b.WriteString("[]")
	}

	if fm.Type == meta.FieldType_Struct {
		b.WriteString("*")
	}

	// 字段类型映射go的类型
	switch fm.Type {
	case meta.FieldType_Integer:
		b.WriteString("int")
	case meta.FieldType_Bool:
		b.WriteString("bool")
	case meta.FieldType_Struct,
		meta.FieldType_Enum:
		b.WriteString(fm.Complex.Name)
	case meta.FieldType_Float32:
		b.WriteString("int32")
	case meta.FieldType_Float64:
		b.WriteString("int64")
	case meta.FieldType_Bytes:
		b.WriteString("[]byte")
	default:
		b.WriteString(fm.Type.String())
	}

	return b.String()
}

func (fm *fieldModel) GoTags() string {

	var b bytes.Buffer

	b.WriteString("`sproto:\"")

	// 整形类型对解码层都视为整形
	switch fm.Type {
	case meta.FieldType_Int32,
		meta.FieldType_Int64,
		meta.FieldType_UInt32,
		meta.FieldType_UInt64,
		meta.FieldType_Float32,
		meta.FieldType_Float64,
		meta.FieldType_Enum:
		b.WriteString("integer")
	case meta.FieldType_Bytes:
		b.WriteString("string")
	default:
		b.WriteString(fm.Kind())
	}

	b.WriteString(",")

	b.WriteString(fmt.Sprintf("%d", fm.TagNumber()))
	b.WriteString(",")

	if fm.Repeatd {
		b.WriteString("array,")
	}

	b.WriteString(fmt.Sprintf("name=%s", fm.GoFieldName()))

	b.WriteString("\"`")

	return b.String()
}

func gen_go(fm *fileModel, filename string) {
	addData(fm, "go")
	generateCode("sp->go", goCodeTemplate, filename, fm, &generateOption{
		formatGoCode: true,
	})
}
