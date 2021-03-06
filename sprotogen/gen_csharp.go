package main

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/davyxu/gosproto/meta"
)

const csharpCodeTemplate = `// Generated by github.com/davyxu/gosproto/sprotogen
// DO NOT EDIT!
using System;
using Sproto;
using System.Collections.Generic;

namespace {{.PackageName}}
{
{{range $a, $enumobj := .Enums}}
	public enum {{.Name}} {
		{{range .StFields}}
		{{.Name}} = {{.TagNumber}},
		{{end}}
	}
{{end}}

{{range .Structs}}
	{{.CSClassAttr}}
	public class {{.Name}} : SprotoTypeBase {
		private static int max_field_count = {{.MaxFieldCount}};
		
		{{range .StFields}}
		[SprotoHasField]
		public bool Has{{.UpperName}}{
			get { return base.has_field.has_field({{.FieldIndex}}); }
		}
		{{.CSFieldAttr}}
		private {{.CSTypeString}} _{{.Name}}; // tag {{.TagNumber}}
		public {{.CSTypeString}} {{.Name}} {
			get{ return _{{.Name}}; }
			set{ base.has_field.set_field({{.FieldIndex}},true); _{{.Name}} = value; }
		}
		{{end}}
		
		public {{.Name}}() : base(max_field_count) {}
		
		public {{.Name}}(byte[] buffer) : base(max_field_count, buffer) {
			this.decode ();
		}
		
		protected override void decode () {
			int tag = -1;
			while (-1 != (tag = base.deserialize.read_tag ())) {
				switch (tag) {
				{{range .StFields}}
				case {{.TagNumber}}:
					this.{{.Name}} = base.deserialize.{{.CSReadFunc}}{{.CSTemplate}}({{.CSLamdaFunc}});
					break;
				{{end}}
				default:
					base.deserialize.read_unknow_data ();
					break;
				}
			}
		}
		
		public override int encode (SprotoStream stream) {
			base.serialize.open (stream);

			{{range .StFields}}
			if (base.has_field.has_field ({{.FieldIndex}})) {
				base.serialize.{{.CSWriteFunc}}(this.{{.Name}}, {{.TagNumber}});
			}
			{{end}}

			return base.serialize.close ();
		}
	}
{{end}}

    public class RegisterEntry
    {
        static readonly Type[] _types = new Type[]{ {{range .Structs}}
                typeof({{.Name}}), // {{.MsgID}}{{end}}
            };

        public static Type[] GetClassTypes()
        {
            return _types;
        }
{{if .EnumValueGroup}}
        public static string ResultToString( int result )
        {
            switch( result)
            {
            	case 0: return "OK";
			{{range $a, $enumObj := .Enums}} {{if .IsResultEnum}} {{range .Fields}} {{if ne .TagNumber 0}}
				case {{.TagNumber}}: return "{{$enumObj.Name}}.{{.Name}}"; {{end}} {{end}} {{end}} {{end}}
            }

            return "Result: " + result.ToString();
        }
{{end}}
    }
}
`

func (fm *fieldModel) CSTemplate() string {
	var buf bytes.Buffer
	var needTemplate bool
	switch fm.Type {
	case meta.FieldType_Struct,
		meta.FieldType_Enum:
		needTemplate = true
	}
	if needTemplate {
		buf.WriteString("<")
	}
	if fm.MainIndex != nil {
		buf.WriteString(csharpTypeName(fm.MainIndex))
		buf.WriteString(",")
	}
	if needTemplate {
		buf.WriteString(fm.Complex.Name)
		buf.WriteString(">")
	}
	return buf.String()
}

func (fm *fieldModel) CSLamdaFunc() string {
	if fm.MainIndex == nil {
		return ""
	}

	return fmt.Sprintf("v => v.%s", fm.MainIndex.Name)
}

func (fm *fieldModel) CSWriteFunc() string {

	return "write_" + fm.serializer()
}

func (fm *fieldModel) CSReadFunc() string {

	funcName := "read_"

	if fm.Repeatd {

		if fm.MainIndex != nil {
			return funcName + "map"
		} else {
			return funcName + fm.serializer() + "_list"
		}

	}

	return funcName + fm.serializer()
}

func (fm *fieldModel) serializer() string {

	var baseName string

	switch fm.Type {
	case meta.FieldType_Integer:
		baseName = "integer"
	case meta.FieldType_Int32:
		baseName = "int32"
	case meta.FieldType_Int64:
		baseName = "int64"
	case meta.FieldType_UInt32:
		baseName = "uint32"
	case meta.FieldType_UInt64:
		baseName = "uint64"
	case meta.FieldType_Float32:
		baseName = "float32"
	case meta.FieldType_Float64:
		baseName = "double"
	case meta.FieldType_String:
		baseName = "string"
	case meta.FieldType_Bool:
		baseName = "boolean"
	case meta.FieldType_Struct:
		baseName = "obj"
	case meta.FieldType_Enum:
		baseName = "enum"
	case meta.FieldType_Bytes:
		baseName = "bytes"
	default:
		baseName = "unknown"
	}

	return baseName
}

func (fm *fieldModel) CSTypeName() string {
	// 字段类型映射go的类型
	return csharpTypeName(fm.FieldDescriptor)
}

func csharpTypeName(fd *meta.FieldDescriptor) string {
	switch fd.Type {
	case meta.FieldType_Integer:
		return "Int64"
	case meta.FieldType_Int32:
		return "Int32"
	case meta.FieldType_Int64:
		return "Int64"
	case meta.FieldType_UInt32:
		return "UInt32"
	case meta.FieldType_UInt64:
		return "UInt64"
	case meta.FieldType_Float32:
		return "float"
	case meta.FieldType_Float64:
		return "double"
	case meta.FieldType_String:
		return "string"
	case meta.FieldType_Bool:
		return "bool"
	case meta.FieldType_Bytes:
		return "byte[]"
	case meta.FieldType_Struct,
		meta.FieldType_Enum:
		return fd.Complex.Name
	}
	return "unknown"
}

func (fm *fieldModel) CSTypeString() string {

	var b bytes.Buffer
	if fm.Repeatd {

		if fm.MainIndex != nil {
			b.WriteString("Dictionary<")

			b.WriteString(csharpTypeName(fm.MainIndex))

			b.WriteString(",")

		} else {
			b.WriteString("List<")
		}

	}

	b.WriteString(fm.CSTypeName())

	if fm.Repeatd {
		b.WriteString(">")
	}

	return b.String()
}

func (fm *fieldModel) CSFieldAttr() string {
	return fm.st.f.CSFieldAttr
}

func (fm *structModel) CSClassAttr() string {
	return fm.f.CSClassAttr
}

func gen_csharp(fm *fileModel, filename string) {

	addData(fm, "cs")

	sort.Sort(fm)

	generateCode("sp->cs", csharpCodeTemplate, filename, fm, nil)

}
