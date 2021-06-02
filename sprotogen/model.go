package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/davyxu/gosproto/meta"
)

type fieldModel struct {
	*meta.FieldDescriptor
	FieldIndex int
	st         *structModel
}

func (fm *fieldModel) UpperName() string {
	return strings.ToUpper(string(fm.Name[0])) + fm.Name[1:]
}

type structModel struct {
	*meta.Descriptor
	StFields []fieldModel
	f        *fileModel
}

func (fm *structModel) IsResultEnum() bool {
	return fm.IsEnum() && strings.HasSuffix(fm.Name, "Result")
}

func (fm *structModel) IsEnum() bool {
	return fm.Type == meta.DescriptorType_Enum
}

func (fm *structModel) IsStruct() bool {
	return fm.Type == meta.DescriptorType_Struct
}

func (fm *structModel) MsgID() uint32 {
	return StringHash(fm.MsgFullName())
}

func (fm *structModel) MsgFullName() string {
	return fm.f.PackageName + "." + fm.Name
}
func (fm *structModel) FieldCount() int {
	return len(fm.StFields)
}

type fileModel struct {
	*meta.FileDescriptorSet

	Structs []*structModel

	Enums []*structModel

	Objects []*structModel

	PackageName string

	CellnetReg bool

	forceAutoTag bool

	CSClassAttr string

	CSFieldAttr string

	EnumValueGroup bool

	MD5 string
}

func (fm *fileModel) Len() int {
	return len(fm.Structs)
}

func (fm *fileModel) Swap(i, j int) {
	fm.Structs[i], fm.Structs[j] = fm.Structs[j], fm.Structs[i]
}

func (fm *fileModel) Less(i, j int) bool {

	a := fm.Structs[i]
	b := fm.Structs[j]

	return a.Name < b.Name
}

func addStruct(fm *fileModel, fileD *meta.FileDescriptor, srcName string) {
	for _, st := range fileD.Objects {
		// 过滤, 只输出某个来源
		if srcName != "" && st.SrcName != srcName {
			continue
		}
		stModel := &structModel{
			Descriptor: st,
		}
		for index, fd := range st.Fields {
			fdModel := fieldModel{
				FieldDescriptor: fd,
				FieldIndex:      index,
				st:              stModel,
			}
			stModel.StFields = append(stModel.StFields, fdModel)
		}
		stModel.f = fm
		fm.Objects = append(fm.Objects, stModel)
		switch stModel.Type {
		case meta.DescriptorType_Enum:
			fm.Enums = append(fm.Enums, stModel)
		case meta.DescriptorType_Struct:
			fm.Structs = append(fm.Structs, stModel)
		}
	}
}

func addData(fm *fileModel, matchTag string) {
	var md5StrList []string
	for _, file := range fm.FileDescriptorSet.Files {
		md5StrList = append(md5StrList, hashFile(file.FileName))
		if file.MatchTag(matchTag) {
			addStruct(fm, file, "")
		}
	}
	sort.Slice(md5StrList, func(i, j int) bool {
		return md5StrList[i] < md5StrList[j]
	})
	fm.MD5 = fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(md5StrList, ""))))
}

func hashFile(filePath string) string {
	data, _ := ioutil.ReadFile(filePath)
	s := string(data)
	reg := regexp.MustCompile(`[a-zA-Z]`)
	result := reg.FindAllStringSubmatch(s, -1)
	s = ""
	for _, strList := range result {
		s += strings.Join(strList, "")
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
