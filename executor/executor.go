package executor

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Function/structFunction"
)

type Executor struct {
	db    *gorm.DB
	stack [][]structFunction.Value
}

func NewExecutor(db *gorm.DB) *Executor {
	var executor *Executor = new(Executor)
	executor.db = db

	return executor
}

func (executor *Executor) Sum(val1 int64, val2 int64) int64 {
	return val1 + val2
}

func (executor *Executor) Diff(val1 int64, val2 int64) int64 {
	return val1 - val2
}

func (executor *Executor) GetParam(paramName string) int64 {
	length := len(executor.stack)
	var result int64 = 0
	argsResolved := executor.stack[length-1]
	for i := 0; i < len(argsResolved); i++ {
		if argsResolved[i].Name == paramName {
			result = argsResolved[i].Valuenumber
			break
		}
	}

	return result
}

func (executor *Executor) GetValueFromInvoke(invoke *structFunction.Invoke) (*structFunction.Value, error) {

	var function structFunction.Function
	valueResult := new(structFunction.Value)

	//gorm:auto_preload se usa para que complete todos los struct con su informacion
	if err := executor.db.Set("gorm:auto_preload", true).First(&function, "name = ?", invoke.Functionname).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if function.Origin == "primitive" {
		results, err := call(function, invoke.Args)
		if err != nil {
			return nil, err
		}

		if len(results) == 1 {
			valueResult.Valuenumber = results[0].Int()
			return valueResult, nil
		}
	} else {

		err := errors.New("The function is not primitive.")
		return nil, err

		//formula de usurio
		argsResolved := make([]structFunction.Value, len(function.Params))
		for i := 0; i < len(function.Params); i++ {
			var valueResolved *structFunction.Value
			//var err error
			if invoke.Args[i].Valueinvoke == nil {
				valueResolved = &invoke.Args[i]
			} else {
				valueResolved, err = executor.GetValueFromInvoke(invoke.Args[i].Valueinvoke)
				if err != nil {
					return nil, err
				}
			}

			switch function.Params[i].Type {
			case "number":
				argsResolved[i].Name = function.Params[i].Name
				argsResolved[i].Valuenumber = valueResolved.Valuenumber
			}
		}

		// Push on Stack
		executor.stack = append(executor.stack, argsResolved)
		result, err := executor.GetValueFromInvoke(invoke)
		if err != nil {
			return nil, err
		}
		// Pop on Stack
		_, executor.stack = executor.stack[len(executor.stack)-1], executor.stack[:len(executor.stack)-1]

		return result, nil
	}

	return nil, nil
}

func call(function structFunction.Function, args []structFunction.Value) (result []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(&Executor{})
	m := myClassValue.MethodByName(function.Name)
	if !m.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("Method not found \"%s\"", function.Name)
	}
	//f := reflect.ValueOf(function.Name)
	if len(args) != m.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(args))
	for k, arg := range args {
		value := arg.Valuenumber
		in[k] = reflect.ValueOf(value)

	}
	result = m.Call(in)
	return
}
