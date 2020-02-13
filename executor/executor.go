package executor

import (
	"errors"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Function/structFunction"
)

type Executor struct {
	Db     *gorm.DB
	Invoke structFunction.Invoke
}

func NewExecutor(db *gorm.DB, invoke structFunction.Invoke) *Executor {
	var executor *Executor = new(Executor)
	executor.Db = db
	executor.Invoke = invoke

	return executor
}

func (executor *Executor) GetValue() (*structFunction.Value, error) {
	var function structFunction.Function
	var valueResult *structFunction.Value

	//gorm:auto_preload se usa para que complete todos los struct con su informacion
	if err := executor.Db.Set("gorm:auto_preload", true).First(&function, "name = ?", executor.Invoke.Functionname).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if function.Origin == "primitive" {
		result, err := call(function, executor.Invoke.Args)
		if err != nil {
			return nil, err
		}

		if len(result) == 1 {
			valueResult.Valuenumber = result[0].Interface().(int64)
			return valueResult, nil
		}
	} else {
		err := errors.New("The function is not primitive.")
		return nil, err
	}

	return nil, nil
}

func call(function structFunction.Function, args []structFunction.Value) (result []reflect.Value, err error) {
	f := reflect.ValueOf(function.Name)
	if len(args) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(args))
	for k, arg := range args {
		value := arg.Valuenumber
		in[k] = reflect.ValueOf(value)
	}
	result = f.Call(in)
	return
}

func sum(val1 int64, val2 int64) int64 {
	return val1 + val2
}
