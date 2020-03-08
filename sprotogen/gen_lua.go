package main

import (
	"github.com/davyxu/gosproto/meta"
)

const luaCodeTemplate = `-- Generated by github.com/davyxu/gosproto/sprotogen
-- DO NOT EDIT!
local sproto = require "3rd/sproto/sproto"
{{if .EnumValueGroup}}
ResultToString = function ( result )
	if result == 0 then
		return "OK"
	end

	local str = ResultByID[result]
	if str == nil then
		return string.format("unknown result: %d", result )
	end

	return str
end

ResultByID = {
	{{range $a, $enumObj := .Enums}} {{if .IsResultEnum}} {{range .Fields}} {{if ne .TagNumber 0}}
	[{{.TagNumber}}] = "{{$enumObj.Name}}.{{.Name}}", {{end}} {{end}} {{end}} {{end}}
}

{{end}}

Enum = {
{{range $a, $enumObj := .Enums}}
	{{$enumObj.Name}} = { {{range .Fields}}
		{{.Name}} = {{.TagNumber}}, {{end}}
	},
	{{end}}
}

local {{.PackageName}} = {
	Schema = sproto.parse [[
{{range .Structs}}
.{{.Name}} {	{{range .StFields}}	
	{{.Name}} {{.TagNumber}} : {{.CompatibleTypeString}} {{end}}
}
{{end}}
	]],

	Names = { {{range .Structs}}
		"{{.Name}}", {{end}}
	},
}

function {{.PackageName}}.GetID(msgName)
    for i, v in ipairs({{.PackageName}}.Names) do
        if v == msgName then
            return i - 1
        end
    end
end

function {{.PackageName}}.Encode(msgName, msgTable)
	logDebug('send msg:', msgName)
    for k, v in pairs(msgTable) do
        logDebug(k, v)
    end
	return {{.PackageName}}.Schema:encode(msgName, msgTable)
end

function {{.PackageName}}.Decode(msgName, msgBytes)
	local data = {{.PackageName}}.Schema:decode(msgName, msgBytes)
	if msgName ~= 'S2C_SystemTime' then
		logDebug('received msg:', msgName)
		for k, v in pairs(data) do
			logDebug(k, v)
		end
	end
	return data
end

return {{.PackageName}}

`

func (self *fieldModel) LuaDefaultValueString() string {

	if self.Repeatd {
		return "nil"
	}

	switch self.Type {
	case meta.FieldType_Bool:
		return "false"
	case meta.FieldType_Int32,
		meta.FieldType_Int64,
		meta.FieldType_UInt32,
		meta.FieldType_UInt64,
		meta.FieldType_Integer,
		meta.FieldType_Float32,
		meta.FieldType_Float64,
		meta.FieldType_Enum:
		return "0"
	case meta.FieldType_String:
		return "\"\""
	case meta.FieldType_Struct,
		meta.FieldType_Bytes:
		return "nil"
	}

	return "unknown type" + self.Type.String()
}

func gen_lua(fm *fileModel, filename string) {

	addData(fm, "lua")

	generateCode("sp->lua", luaCodeTemplate, filename, fm, nil)

}
