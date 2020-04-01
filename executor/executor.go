package executor

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Function/structFunction"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

type Executor struct {
	db      *gorm.DB
	stack   [][]structFunction.Value
	context *Context //[]byte
}

type Context struct {
	Currentliquidacion structLiquidacion.Liquidacion `json:"currentliquidacion"`
}

func NewExecutor(db *gorm.DB, context *Context) *Executor {
	var executor *Executor = new(Executor)
	executor.db = db
	executor.stack = [][]structFunction.Value{}
	executor.context = context

	return executor
}

func (executor *Executor) GetValueResolved(value *structFunction.Value) (*structFunction.Value, error) {
	if value.Valueinvoke == nil {
		return value, nil
	} else {
		valueResolved, err := executor.GetValueFromInvoke(value.Valueinvoke)
		if err != nil {
			return nil, err
		}

		return valueResolved, nil
	}
}

func (executor *Executor) GetValueFromInvoke(invoke *structFunction.Invoke) (*structFunction.Value, error) {

	var function structFunction.Function
	//valueResult := new(structFunction.Value)

	//gorm:auto_preload se usa para que complete todos los struct con su informacion
	if err := executor.db.Set("gorm:auto_preload", true).First(&function, "name = ?", invoke.Functionname).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if function.Origin == "primitive" {
		valueResult, err := executor.call(function, invoke.Args)
		if err != nil {
			return nil, err
		}

		return valueResult, nil
	} else {
		// Formula de usuario

		if len(invoke.Args) != len(function.Params) {
			err := errors.New("The number of params is not adapted.")
			return nil, err
		}

		argsResolved := make([]structFunction.Value, len(function.Params))
		for i := 0; i < len(function.Params); i++ {
			valueResolved, err := executor.GetValueResolved(&invoke.Args[i])
			if err != nil {
				return nil, err
			}
			valueResolved.Name = function.Params[i].Name

			argsResolved[i] = *valueResolved
		}

		// Push on Stack
		executor.stack = append(executor.stack, argsResolved)
		result, err := executor.GetValueFromInvoke(function.Value.Valueinvoke)
		if err != nil {
			return nil, err
		}
		// Pop on Stack
		_, executor.stack = executor.stack[len(executor.stack)-1], executor.stack[:len(executor.stack)-1]

		return result, nil
	}

	return nil, nil
}

func (executor *Executor) call(function structFunction.Function, args []structFunction.Value) (result *structFunction.Value, err error) {
	myClassValue := reflect.ValueOf(executor)
	m := myClassValue.MethodByName(function.Name)
	if !m.IsValid() {
		return nil, fmt.Errorf("Method not found \"%s\"", function.Name)
	}
	//f := reflect.ValueOf(function.Name)
	if len(args) != m.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return nil, err
	}
	in := make([]reflect.Value, len(args))
	for k, arg := range args {
		valueResolved, err := executor.GetValueResolved(&arg)
		if err != nil {
			return nil, err
		}
		paramType := m.Type().In(k).Name()
		var value interface{}
		switch paramType {
		case "float64":
			value = valueResolved.Valuenumber
		case "string":
			value = valueResolved.Valuestring
			/*default:
			jsonbody, err := json.Marshal(valueResolved.Valueobject)
			if err != nil {
				return nil, err
			}
			value = jsonbody*/
		}

		in[k] = reflect.ValueOf(value)
	}
	results := m.Call(in)

	result = new(structFunction.Value)
	if len(results) == 1 {
		paramType := m.Type().Out(0).Name()
		switch paramType {
		case "float64":
			val := results[0]
			result.Valuenumber = val.Float()
		case "string":
			result.Valuestring = results[0].String()
		case "bool":
			val := results[0]
			result.Valueboolean = val.Bool()
		}

		return result, nil
	}

	return result, nil
}
