package asm_vm_stack_base

import (
	"fmt"
	"go-compiler/utils"
)

const (
	BuiltinFuncNameLen  = "len"
	BuiltinFuncNamePush = "push"
	BuiltinFuncNameLog  = "log"
)

// 暂时用全局吧
var Builtins = []struct {
	Name    string
	Builtin *BuiltinObject
}{
	{
		BuiltinFuncNameLen,
		&BuiltinObject{Func: func(args ...Object) Object {
			if len(args) != 1 {
				utils.LogErrorFormat("wrong number of arguments. got=%d, want=1",
					len(args))
				return nil
			}

			switch arg := args[0].(type) {
			case *ArrayObject:
				return &NumberObject{Value: float64(len(arg.Values))}
			case *StringObject:
				return &NumberObject{Value: float64(len(arg.Value))}
			default:
				utils.LogErrorFormat("argument to %q not supported, got %s",
					BuiltinFuncNameLen, args[0].ValueType())
				return nil
			}
		},
		},
	},
	{
		BuiltinFuncNameLog,
		&BuiltinObject{Func: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return nil
		},
		},
	},

	{
		BuiltinFuncNamePush,
		&BuiltinObject{Func: func(args ...Object) Object {
			if len(args) != 2 {
				utils.LogErrorFormat("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			arrValueType := ArrayObject{}.ValueType()
			if args[0].ValueType() != arrValueType {
				utils.LogErrorFormat("argument to %q must be %s, got %s",
					BuiltinFuncNamePush, arrValueType, args[0].ValueType())
			}

			arr := args[0].(*ArrayObject)
			length := len(arr.Values)

			newElements := make([]Object, length+1, length+1)
			copy(newElements, arr.Values)
			newElements[length] = args[1]

			return &ArrayObject{Values: newElements}
		},
		},
	},
}

func GetBuiltinByName(name string) *BuiltinObject {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
