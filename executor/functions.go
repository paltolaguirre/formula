package executor

import (
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

	start := legajo.Fechaalta
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
		if item.Conceptoid == -17 { // Días Faltas Injustificadas
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

	liquidacionDate := liquidacion.Fecha

	if (liquidacionDate.Month() - legajo.Fechaalta.Month()) <= 3 {
		return executor.Sueldo() / 2
	}

	if (liquidacionDate.Year() - legajo.Fechaalta.Year()) <= 5 {
		return executor.Sueldo()
	}

	return executor.Sueldo() * 2
}

func (executor *Executor) SacSinPreaviso() float64 {
	return executor.Preaviso() / 12
}

func (executor *Executor) IntegracionMesDespido() float64 {
	diasMesTrabajados := executor.DiasMesTrabajadosFechaLiquidacion()
	jornal := executor.Jornal()

	return (30 - diasMesTrabajados) * jornal
}

func (executor *Executor) SacSinImd() float64 {
	return executor.IntegracionMesDespido() / 12
}

func (executor *Executor) Sac() float64 {
	// Falta la formula DiasSemTrabajados()
	return 0 // (executor.MejorRemRemunerativaSemestre() / 2) * executor.DiasSemTrabajados() / 180
}

func (executor *Executor) SacNoRemunerativo() float64 {
	// Falta la formula DiasSemTrabajados()
	return 0 //(executor.MejorRemNoRemunerativaSemestre() / 2) * executor.DiasSemTrabajados() / 180
}

func (executor *Executor) CantidadSueldos() float64 {
	// Falta la formula Resto()
	if 0 /*executor.Resto()*/ < 0.25 {
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
