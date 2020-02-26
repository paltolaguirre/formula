package executor

import (
	"encoding/json"
	"fmt"

	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

func (executor *Executor) GetParamValue(paramName string) int64 {
	length := len(executor.stack)
	var result int64 = 0
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

func (executor *Executor) Sum(val1 int64, val2 int64) int64 {
	return val1 + val2
}

func (executor *Executor) Diff(val1 int64, val2 int64) int64 {
	return val1 - val2
}

func (executor *Executor) TotalImporteRemunerativo(object []byte) int64 {
	liquidacion := structLiquidacion.Liquidacion{}
	if err := json.Unmarshal(object, &liquidacion); err != nil {
		// do error check
		fmt.Println(err)
		return 0
	}
	var total float64 = 0
	for _, item := range liquidacion.Liquidacionitems {
		if item.Concepto.Tipoconcepto.Codigo == "IMPORTE_REMUNERATIVO" {
			total += *item.Importeunitario
		}
	}

	return int64(total)
}
