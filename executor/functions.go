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

func (executor *Executor) TotalImporteRemunerativo() float64 {
	//context := ContextLiquidacion{}
	/*context, ok := (*executor.context).(ContextLiquidacion)

	if !ok {
		fmt.Println("Error")
		return 0
	}*/

	/*if err := json.Unmarshal(executor.context, &context); err != nil {
		// do error check
		fmt.Println(err)
		return 0
	}*/

	liquidacion := executor.context.Currentliquidacion
	var total float64 = 0
	for _, item := range liquidacion.Liquidacionitems {
		if item.Concepto.Tipoconcepto.Codigo == "IMPORTE_REMUNERATIVO" {
			total += *item.Importeunitario
		}
	}

	return float64(total)
}

/* comparison operators */
func (executor *Executor) Greater(val1 float64, val2 float64) bool {
	return val1 > val2
}

func (executor *Executor) GreaterEqual(val1 float64, val2 float64) bool {
	return val1 >= val2
}
