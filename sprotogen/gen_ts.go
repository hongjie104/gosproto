package main

import (
	"bytes"
	"strings"

	"github.com/davyxu/gosproto/meta"
)

const tsCodeTemplate = `
namespace socket {
    const _proto =&
{{range .Structs}}
.{{.Name}} { {{range .StFields}}
	{{.Name}} {{.TagNumber}} : {{.CompatibleTypeString}}{{end}}
}
{{end}}&;

    const _msg = [{{range .Structs}}
		"{{.Name}}",{{end}}
	];

    const _sp = new sproto.Sproto(_proto);

    let _socket: Socket;

    export const event = new egret.EventDispatcher();

    export const md5 = "{{.MD5}}";

    export function init() {
        if (!_socket) {
            _socket = new Socket(() => {
                const byte = new egret.ByteArray();
                byte.endian = egret.Endian.LITTLE_ENDIAN;
                _socket.readBytes(byte);
                const msgId = byte.readShort();
                const dataByteArray = new egret.ByteArray();
                dataByteArray.endian = egret.Endian.LITTLE_ENDIAN;
                byte.readBytes(dataByteArray);
                const msgName = _getMsgName(msgId);
                const data = _sp.decode(msgName, dataByteArray);
                event.dispatchEvent(new socket.MessageEvent(msgName, data));
                if (DEBUG) {
					if ('S2C_SystemTime' !== msgName) {
						console.log(&========received msg = ${msgName}==========&);
						for (const key in data) {
							console.log(&key = ${key}, val =&);
							console.log(data[key]);
						}
					}
				}
            });
        }
    }

    export function send(msg: string, params: any = {}) {
        if (DEBUG) {
			console.log(&========send msg = ${msg}==========&);
			for (const key in params) {
                console.log(&key = ${key}, val =&);
				console.log(params[key]);
			}
		}
        const buffer = _sp.encode(msg, params);
        const byteArray = new egret.ByteArray();
        byteArray.endian = egret.Endian.LITTLE_ENDIAN;
        const msgId = _getID(msg);
        if (msgId < 0) {
            console.error(&通讯错误，${msg}不存在&);
        } else {
            byteArray.writeShort(msgId);
            byteArray.writeBytes(buffer);
            _socket.writeBytes(byteArray, 0, byteArray.bytesAvailable);
        }
    }

    function _getID(msg: string): number {
        for (let i = 0; i < _msg.length; i++) {
            if (_msg[i] === msg) {
                return i;
            }
        }
        return -1;
    }

    function _getMsgName(msgId: number):string {
        return _msg[msgId];
    }
}
`

func (self *fieldModel) TSTypeName() string {
	var b bytes.Buffer
	// 字段类型映射go的类型
	switch self.Type {
	case meta.FieldType_Bool:
		b.WriteString("boolean")
	case meta.FieldType_Int64,
		meta.FieldType_Int32,
		meta.FieldType_Float64,
		meta.FieldType_Float32,
		meta.FieldType_Integer:
		b.WriteString("number")
	case meta.FieldType_Struct,
		meta.FieldType_Enum:
		b.WriteString(self.Complex.Name)
	default:
		b.WriteString(self.Type.String())
	}
	if self.Repeatd {
		b.WriteString("[]")
	}

	return b.String()
}

func gen_ts(fm *fileModel, filename string) {
	addData(fm, "ts")
	generateCode("sp->ts", strings.ReplaceAll(tsCodeTemplate, "&", "`"), filename, fm, nil)
}
