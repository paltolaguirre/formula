package executor

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Function/structFunction"
)

type Executor struct {
	db     *gorm.DB
	invoke *structFunction.Invoke
}

func NewExecutor(db *gorm.DB, invoke *structFunction.Invoke) *Executor {
	var executor *Executor = new(Executor)
	executor.db = db
	executor.invoke = invoke

	return executor
}

func (executor *Executor) Sum(val1 int64, val2 int64) int64 {
	return val1 + val2
}

func (executor *Executor) GetValue() (*structFunction.Value, error) {
	var function structFunction.Function
	valueResult := new(structFunction.Value)

	//gorm:auto_preload se usa para que complete todos los struct con su informacion
	if err := executor.db.Set("gorm:auto_preload", true).First(&function, "name = ?", executor.invoke.Functionname).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	if function.Origin == "primitive" {
		results, err := call(function, executor.invoke.Args)
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
