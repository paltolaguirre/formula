package executor

import (
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

type ContextLiquidacion struct {
	Currentliquidacion structLiquidacion.Liquidacion `json:"currentliquidacion"`
}

func (executor *Executor) GetParamValue(paramName string) float64 {
	length := len(executor.stack)
	var result float64 = 0
	if length > 0 {
		argsResolved := executor.stack[length-1]
		for i := 0; i < len(argsResolved); i++ {
			if argsResolved[i].Name == paramName {
				result = argsResolved[i].Valuenumber
				break
			}
		}
	}

	return result
}

func (executor *Executor) Sum(val1 float64, val2 float64) float64 {
	return val1 + val2
}

func (executor *Executor) Diff(val1 float64, val2 float64) float64 {
	return val1 - val2
}

/* comparison operators */
func (executor *Executor) Greater(val1 float64, val2 float64) bool {
	return val1 > val2
}

func (executor *Executor) GreaterEqual(val1 float64, val2 float64) bool {
	return val1 >= val2
}

/* FORMULAS DE XUBIO */
func (executor *Executor) Jornal() float64 {
	return executor.Sueldo() / 30
}
