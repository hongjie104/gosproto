package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/davyxu/gosproto/meta"
)

// const spCodeTemplate = `
// {{range $a, $obj := .Objects}}
// {{.SpLeadingComment}}
// {{.TypeName}} {{.Name}} {
// 	{{range .StFields}}{{.SpLeadingComment}}{{if $obj.IsStruct}}
// 	{{.SpFieldString}}
// 	{{else}}
// 	{{.Name}} = {{.TagNumber}}{{.SpTrailingComment}}
// 	{{end}}	{{end}}
// }
// {{end}}

// `

func addCommentSignAtEachLine(sign, comment string) string {

	if comment == "" {
		return ""
	}
	var out bytes.Buffer

	scanner := bufio.NewScanner(strings.NewReader(comment))

	var index int
	for scanner.Scan() {

		if index > 0 {
			out.WriteString("\n")
		}

		out.WriteString(sign)
		out.WriteString(scanner.Text())

		index++
	}

	return out.String()
}

func (fm *fieldModel) SpFieldString() string {
	if fm.st.f.forceAutoTag {
		return fmt.Sprintf("%s %s%s", fm.Name, fm.TypeString(), fm.SpTrailingComment())
	} else {
		return fmt.Sprintf("%s %s %s%s", fm.Name, fm.TagString(), fm.TypeString(), fm.SpTrailingComment())
	}
}

func (fm *fieldModel) SpTrailingComment() string {

	return addCommentSignAtEachLine("//", fm.Trailing)
}

func (fm *fieldModel) TagString() string {
	if fm.AutoTag == -1 {
		return strconv.Itoa(fm.Tag)
	}

	return ""
}

func (fm *fieldModel) SpLeadingComment() string {

	return addCommentSignAtEachLine("//", fm.Leading)
}

func (fm *structModel) SpLeadingComment() string {

	return addCommentSignAtEachLine("//", fm.Leading)
}

func (fm *structModel) TypeName() string {
	switch fm.Type {
	case meta.DescriptorType_Enum:
		return "enum"
	case meta.DescriptorType_Struct:
		return "message"
	}

	return "none"
}

// func gen_sp(fileset *meta.FileDescriptorSet, forceAutoTag bool) {
// 	//for srcName := range fileD.ObjectsBySrcName {
// 	//	fm := &fileModel{
// 	//		FileDescriptorSet: fileset,
// 	//		forceAutoTag:      forceAutoTag,
// 	//	}
// 	//
// 	//	addStruct(fm, fileD, srcName)
// 	//
// 	//	generateCode("sp->sp", spCodeTemplate, srcName, fm, nil)
// 	//}

// }
