package executor

import (
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Legajo/structLegajo"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"math"
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

func (executor *Executor) GetConceptValue(id float64) float64 {
	liquidacion := executor.context.Currentliquidacion
	var acumulador float64 = 0

	for _, item := range liquidacion.Liquidacionitems {
		if item.Concepto.ID == int(id) && item.Importeunitario != nil {
			acumulador += *item.Importeunitario
		}
	}

	return acumulador
}

func (executor *Executor) If(expression bool, valueTrue float64, valueFalse float64) float64 {
	if expression {
		return valueTrue
	}

	return valueFalse
}

/* Comparison operators */
func (executor *Executor) Equality(val1 float64, val2 float64) bool {
	return val1 == val2
}

func (executor *Executor) Inequality(val1 float64, val2 float64) bool {
	return val1 != val2
}

func (executor *Executor) BooleanInequality(val1 bool, val2 bool) bool {
	return val1 != val2
}

func (executor *Executor) Greater(val1 float64, val2 float64) bool {
	return val1 > val2
}

func (executor *Executor) GreaterEqual(val1 float64, val2 float64) bool {
	return val1 >= val2
}

func (executor *Executor) Less(val1 float64, val2 float64) bool {
	return val1 < val2
}

func (executor *Executor) LessEqual(val1 float64, val2 float64) bool {
	return val1 <= val2
}

/* Logical operators */
func (executor *Executor) Not(val bool) bool {
	return !val
}

func (executor *Executor) And(val1 bool, val2 bool) bool {
	return val1 && val2
}

func (executor *Executor) Or(val1 bool, val2 bool) bool {
	return val1 || val2
}

/* Arithmetic operators */
func (executor *Executor) Percent(val float64, percent float64) float64 {
	return val * (percent / 100)
}

func (executor *Executor) Sum(val1 float64, val2 float64) float64 {
	return val1 + val2
}

func (executor *Executor) Diff(val1 float64, val2 float64) float64 {
	return val1 - val2
}

func (executor *Executor) Div(val1 float64, val2 float64) float64 {
	return val1 / val2
}

func (executor *Executor) Multi(val1 float64, val2 float64) float64 {
	return val1 * val2
}

/* FORMULAS DE XUBIO */
func (executor *Executor) Jornal() float64 {
	return math.Round(executor.Sueldo() / 30 * 100) / 100
}

func (executor *Executor) ValorHora() float64 {
	return math.Round(executor.Sueldo() / executor.HorasMensuales() * 100) / 100
}

func (executor *Executor) HoraExtra50() float64 {
	liquidacion := executor.context.Currentliquidacion
	var totalCantidad float64 = 0
	for _, item := range liquidacion.Liquidacionitems {
		if *item.Conceptoid == -5 { // Horas Extras 50%
			totalCantidad += item.Cantidad
		}
	}

	return executor.ValorHora() * 1.5 * totalCantidad
}

func (executor *Executor) HoraExtra100() float64 {
	liquidacion := executor.context.Currentliquidacion
	var totalCantidad float64 = 0
	for _, item := range liquidacion.Liquidacionitems {
		if *item.Conceptoid == -6 { // Horas Extras 100%
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

	start := *legajo.Fechaalta
	end := liquidacion.Fecha

	period := end.Sub(start)
	days := int(period.Hours() / 24)
	years := int(days / 365)
	antiguedad := float64(years)

	return antiguedad
}

func (executor *Executor) DiasFaltasInjustificadas() float64 {
	liquidacion := executor.context.Currentliquidacion
	var totalCantidad float64 = 0
	for _, item := range liquidacion.Liquidacionitems {
		if *item.Conceptoid == -17 { // Días Faltas Injustificadas
			totalCantidad += item.Cantidad
		}
	}

	return totalCantidad
}

func (executor *Executor) Vacaciones() float64 {

	antiguedad := executor.Antiguedad()

	if antiguedad <= 5 {
		return executor.ValorDiasVacaciones() * 14
	}

	if antiguedad <= 10 {
		return executor.ValorDiasVacaciones() * 21
	}

	if antiguedad <= 15 {
		return executor.ValorDiasVacaciones() * 28
	}

	return executor.ValorDiasVacaciones() * 35
}

func (executor *Executor) Preaviso() float64 {
	liquidacion := executor.context.Currentliquidacion

	var legajo structLegajo.Legajo

	if err := executor.db.Set("gorm:auto_preload", true).First(&legajo, "id = ?", liquidacion.Legajo.ID).Error; gorm.IsRecordNotFoundError(err) {
		return 0
	}

	year, month, _, _, _, _ := diff(liquidacion.Fecha, *legajo.Fechaalta)

	if year == 0 && month <= 3 {
		return executor.Sueldo() / 2
	}

	if year < 5 || (year == 5 && month == 0){
		return executor.Sueldo()
	}

	return executor.Sueldo() * 2
}

func (executor *Executor) SacSinPreaviso() float64 {
	return math.Round(executor.Preaviso() / 12 * 100) / 100
}

func (executor *Executor) IntegracionMesDespido() float64 {
	diasMesTrabajados := executor.DiasMesTrabajadosFechaLiquidacion()
	jornal := executor.Jornal()

	return math.Round(math.Max((30 - diasMesTrabajados) * jornal, 0) * 100) / 100
}

func (executor *Executor) SacSinImd() float64 {
	return math.Round(math.Max(executor.IntegracionMesDespido() / 12, 0) * 100) / 100
}

func (executor *Executor) Sac() float64 {
	return math.Round((executor.MejorRemRemunerativaBaseSACSemestre() / 2) * executor.DiasSemTrabajados() / 180 * 100) / 100
}

func (executor *Executor) SacNoRemunerativo() float64 {
	return math.Round((executor.MejorRemNoRemunerativaSemestre() / 2) * executor.DiasSemTrabajados() / 180 * 100) / 100
}

/*
 * TODO: se podria hacer un refactor para calcular la antiguedad directamente aca. Seria mas optimo.
 */
func (executor *Executor) CantidadSueldos() float64 {
	if executor.AntiguedadResto() < 0.25 {
		return executor.Antiguedad()
	} else {
		return executor.Antiguedad() + 1
	}
}

func (executor *Executor) IndemnizacionPorDespido(importe float64) float64 {
	mejorRem := executor.MejorRemNormalYHabitualSemestre()

	if mejorRem <= importe*3 {
		return mejorRem * executor.CantidadSueldos()
	}

	if mejorRem*0.67 <= importe*3 {
		return mejorRem * 0.67
	}

	return importe * 3 * executor.CantidadSueldos()
}
