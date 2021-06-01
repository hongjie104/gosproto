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

func (self *fieldModel) GoEnumTypeName() string {

	if self.st.EnumValueIgnoreType {
		return ""
	}

	return self.st.Name
}

func (self *fieldModel) GoExtendFieldGetterName() string {
	pname := publicFieldName(self.Name)

	if token.Lookup(pname).IsKeyword() {
		return pname + "_"
	}

	return pname
}

func (self *fieldModel) GoExtendFieldGetter() string {

	switch self.Type {
	case meta.FieldType_Float32,
		meta.FieldType_Float64:
		return fmt.Sprintf("return %s(self.%s) * %f", self.GoExtendFieldGetterType(), self.GoFieldName(), 1.0/float32(self.ExtendTypePrecision()))
	}

	return "unknown extend type:" + self.Type.String()
}

func (self *fieldModel) GoExtendFieldSetter() string {

	switch self.Type {
	case meta.FieldType_Float32:
		return fmt.Sprintf("self.%s = int32(v* %d)", self.GoFieldName(), self.ExtendTypePrecision())
	case meta.FieldType_Float64:
		return fmt.Sprintf("self.%s = int64(v* %d)", self.GoFieldName(), self.ExtendTypePrecision())
	}

	return "unknown extend type:" + self.Type.String()
}

func (self *fieldModel) GoExtendFieldGetterType() string {
	var b bytes.Buffer
	if self.Repeatd {
		b.WriteString("[]")
	}

	// 字段类型映射go的类型
	switch self.Type {
	case meta.FieldType_Float32:
		b.WriteString("float32")
	case meta.FieldType_Float64:
		b.WriteString("float64")
	default:
		b.WriteString("unknown extend type:" + self.Type.String())
	}

	return b.String()
}

func (self *fieldModel) GoFieldName() string {

	var pname string

	// 扩展类型不能直接访问
	if self.IsExtendType() {
		pname = "Extend_" + self.Name
	} else {
		pname = publicFieldName(self.Name)
	}

	// 碰到关键字在尾部加_
	if token.Lookup(pname).IsKeyword() {
		return pname + "_"
	}

	return pname
}

func (self *fieldModel) GoTypeName() string {

	var b bytes.Buffer
	if self.Repeatd {
		b.WriteString("[]")
	}

	if self.Type == meta.FieldType_Struct {
		b.WriteString("*")
	}

	// 字段类型映射go的类型
	switch self.Type {
	case meta.FieldType_Integer:
		b.WriteString("int")
	case meta.FieldType_Bool:
		b.WriteString("bool")
	case meta.FieldType_Struct,
		meta.FieldType_Enum:
		b.WriteString(self.Complex.Name)
	case meta.FieldType_Float32:
		b.WriteString("int32")
	case meta.FieldType_Float64:
		b.WriteString("int64")
	case meta.FieldType_Bytes:
		b.WriteString("[]byte")
	default:
		b.WriteString(self.Type.String())
	}

	return b.String()
}

func (self *fieldModel) GoTags() string {

	var b bytes.Buffer

	b.WriteString("`sproto:\"")

	// 整形类型对解码层都视为整形
	switch self.Type {
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
		b.WriteString(self.Kind())
	}

	b.WriteString(",")

	b.WriteString(fmt.Sprintf("%d", self.TagNumber()))
	b.WriteString(",")

	if self.Repeatd {
		b.WriteString("array,")
	}

	b.WriteString(fmt.Sprintf("name=%s", self.GoFieldName()))

	b.WriteString("\"`")

	return b.String()
}

func gen_go(fm *fileModel, filename string) {
	addData(fm, "go")
	generateCode("sp->go", goCodeTemplate, filename, fm, &generateOption{
		formatGoCode: true,
	})
}
