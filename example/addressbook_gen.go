// Generated by github.com/davyxu/gosproto/sprotogen
// DO NOT EDIT!

package example

import (
	"fmt"
	"github.com/davyxu/cellnet/codec/sproto"
	"github.com/davyxu/gosproto"
	"reflect"
)

type MyCar int32

const (
	MyCar_Monkey MyCar = 0
	MyCar_Monk   MyCar = 1
	MyCar_Pig    MyCar = 2
)

var (
	MyCarMapperValueByName = map[string]int32{
		"Monkey": 0,
		"Monk":   1,
		"Pig":    2,
	}

	MyCarMapperNameByValue = map[int32]string{
		0: "Monkey",
		1: "Monk",
		2: "Pig",
	}
)

func (self MyCar) String() string {
	return sproto.EnumName(MyCarMapperNameByValue, int32(self))
}

type PhoneNumber struct {
	Number string `sproto:"string,0,name=Number"`

	Type int32 `sproto:"integer,1,name=Type"`
}

func (self *PhoneNumber) String() string { return fmt.Sprintf("%+v", *self) }

type Person struct {
	Name string `sproto:"string,0,name=Name"`

	Id int32 `sproto:"integer,1,name=Id"`

	Email string `sproto:"string,2,name=Email"`

	Phone []*PhoneNumber `sproto:"struct,3,array,name=Phone"`
}

func (self *Person) String() string { return fmt.Sprintf("%+v", *self) }

type AddressBook struct {
	Person []*Person `sproto:"struct,0,array,name=Person"`
}

func (self *AddressBook) String() string { return fmt.Sprintf("%+v", *self) }

type MyData struct {
	Name string `sproto:"string,0,name=Name"`

	Type MyCar `sproto:"integer,1,name=Type"`

	Int32 int32 `sproto:"integer,2,name=Int32"`

	Uint32 uint32 `sproto:"integer,3,name=Uint32"`

	Int64 int64 `sproto:"integer,4,name=Int64"`

	Uint64 uint64 `sproto:"integer,5,name=Uint64"`

	Bool bool `sproto:"boolean,6,name=Bool"`

	Extend_Float32 int32 `sproto:"integer,7,name=Extend_Float32"`

	Extend_Float64 int64 `sproto:"integer,8,name=Extend_Float64"`

	Stream []byte `sproto:"string,9,name=Stream"`
}

func (self *MyData) String() string { return fmt.Sprintf("%+v", *self) }

func (self *MyData) Float32() float32 {
	return float32(self.Extend_Float32) * 0.001000
}
func (self *MyData) SetFloat32(v float32) {
	self.Extend_Float32 = int32(v * 1000)
}

func (self *MyData) Float64() float64 {
	return float64(self.Extend_Float64) * 0.001000
}
func (self *MyData) SetFloat64(v float64) {
	self.Extend_Float64 = int64(v * 1000)
}

type MyProfile struct {
	NameField *MyData `sproto:"struct,0,name=NameField"`

	NameArray []*MyData `sproto:"struct,1,array,name=NameArray"`

	NameMap []*MyData `sproto:"struct,2,array,name=NameMap"`
}

func (self *MyProfile) String() string { return fmt.Sprintf("%+v", *self) }

var SProtoStructs = []reflect.Type{

	reflect.TypeOf((*PhoneNumber)(nil)).Elem(), // 4271979557
	reflect.TypeOf((*Person)(nil)).Elem(),      // 1498745430
	reflect.TypeOf((*AddressBook)(nil)).Elem(), // 2618161298
	reflect.TypeOf((*MyData)(nil)).Elem(),      // 2244887298
	reflect.TypeOf((*MyProfile)(nil)).Elem(),   // 438153711
}

var SProtoEnumValue = map[string]map[int32]string{
	"MyCar": MyCarMapperNameByValue,
}

func init() {
	sprotocodec.AutoRegisterMessageMeta(SProtoStructs)
}
