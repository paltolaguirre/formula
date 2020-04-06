package executor

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Legajo/structLegajo"
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

func (executor *Executor) ValorHora() float64 {
	return executor.Sueldo() / executor.HorasMensuales()
}

func (executor *Executor) HoraExtra50() float64 {
	liquidacion := executor.context.Currentliquidacion
	var totalCantidad float64 = 0
	for _, item := range liquidacion.Liquidacionitems {
		if item.Conceptoid == -5 { // Horas Extras 50%
			totalCantidad += item.Cantidad
		}
	}

	return executor.ValorHora() * 1.5 * totalCantidad
}

func (executor *Executor) HoraExtra100() float64 {
	liquidacion := executor.context.Currentliquidacion
	var totalCantidad float64 = 0
	for _, item := range liquidacion.Liquidacionitems {
		if item.Conceptoid == -6 { // Horas Extras 100%
			totalCantidad += item.Cantidad
		}
	}

	return executor.ValorHora() * 2 * totalCantidad
}

func (executor *Executor) Antiguedad() float64 {
	liquidacion := executor.context.Currentliquidacion

	var legajo structLegajo.Legajo

	if err := executor.db.Set("gorm:auto_preload", true).First(&legajo, "id = ?", liquidacion.Legajo.ID).Error; gorm.IsRecordNotFoundError(err) {
		return 0
	}

	currentDate := time.Now()
	antiguedad := currentDate.Year() - legajo.Fechaalta.Year()
	if (currentDate.Month() - legajo.Fechaalta.Month()) > 5 {
		antiguedad++
	}

	return antiguedad
}
