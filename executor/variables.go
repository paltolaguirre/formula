package executor

import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"math"
	"strconv"
	"time"
)

func (executor *Executor) TotalImporteRemunerativo() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "IMPORTE_REMUNERATIVO", false)
}

func (executor *Executor) TotalHaberesNoRemunerativosMensual() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "IMPORTE_NO_REMUNERATIVO", false)
}

func (executor *Executor) TotalDescuentosMensual() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "DESCUENTO", false)
}

func (executor *Executor) TotalRetencionesMensual() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "RETENCION", false)
}

func (executor *Executor) TotalAportesPatronalesMensual() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "APORTE_PATRONAL", false)
}

func (executor *Executor) Sueldo() float64 {
	var importe float64
	legajoid := executor.context.Currentliquidacion.Legajoid;
	if legajoid == nil {
		fmt.Println("Para realizar el calculo automatico de Sueldo, debe seleccionar primero un legajo")
		return 0
	}

	sql := "SELECT SUM( remuneracion ) from legajo where id = " + strconv.Itoa(*legajoid)
	err := executor.db.Raw(sql).Row().Scan(&importe)

	if err != nil {
		fmt.Println("Error al buscar la remuneracion para el legajo " + strconv.Itoa(*legajoid))
		return 0
	}

	return importe
}

func (executor *Executor) ValorDiasVacaciones() float64 {

	return executor.Sueldo() / 25
}

func (executor *Executor) HorasMensuales() float64 {
	var importe float64
	legajoid := executor.context.Currentliquidacion.Legajoid;
	if legajoid == nil {
		fmt.Println("Para realizar el calculo automatico de HorasMensuales, debe seleccionar primero un legajo")
		return 0
	}

	sql := "SELECT SUM( horasmensualesnormales::numeric(19,4) ) from legajo where id = " + strconv.Itoa(*legajoid)
	err := executor.db.Raw(sql).Row().Scan(&importe)

	if err != nil {
		fmt.Println("Error al buscar las HorasMensuales para el legajo " + strconv.Itoa(*legajoid))
		return 0
	}

	return importe
}

func (executor *Executor) DiasMesTrabajadosFechaLiquidacion() float64 {

	var respuesta float64 = 0

	liquidacion := executor.context.Currentliquidacion

	fechaLiquidacion := liquidacion.Fecha //No puede ser null

	anioLiquidacion := fechaLiquidacion.Year()
	mesLiquidacion := fechaLiquidacion.Month()
	diaLiquidacion := fechaLiquidacion.Day()

	legajoid := liquidacion.Legajoid;
	if legajoid == nil {
		println("Para realizar el calculo automatico de DiasMesTrabajadosFechaLiquidacion, debe seleccionar primero un legajo")
		return 0
	}

	var fechaAlta time.Time
	sql := "SELECT fechaalta from legajo where id = " + strconv.Itoa(*legajoid)
	err := executor.db.Raw(sql).Row().Scan(&fechaAlta)

	if err != nil {
		fmt.Println("Error al buscar la fechaalta para el legajo " + strconv.Itoa(*legajoid))
		return 0
	}

	anioAlta := fechaAlta.Year()
	mesAlta := fechaAlta.Month()
	diaAlta := fechaAlta.Day()

	if anioLiquidacion == anioAlta && mesLiquidacion == mesAlta {
		respuesta = math.Max(float64(diaLiquidacion - diaAlta) , 0)
	}

	if anioLiquidacion > anioAlta || (anioLiquidacion == anioAlta  && mesLiquidacion > mesAlta) {
		respuesta = float64(diaLiquidacion)
	}

	if anioLiquidacion < anioAlta || (anioLiquidacion == anioAlta  && mesLiquidacion < mesAlta) {
		respuesta = 0
	}

	return respuesta
}

func (executor *Executor) DiasMesTrabajadosFechaPeriodo() float64 {

	var respuesta float64 = 0

	liquidacion := executor.context.Currentliquidacion

	periodoLiquidacion := liquidacion.Fechaperiodoliquidacion //No puede ser null

	anioPeriodoLiquidacion := periodoLiquidacion.Year()
	mesPeriodoLiquidacion := periodoLiquidacion.Month()

	maximoDiaMesPeriodoLiquidacion := now.New(periodoLiquidacion).EndOfMonth().Day()

	legajoid := liquidacion.Legajoid;
	if legajoid == nil {
		fmt.Println("Para realizar el calculo automatico de Sueldo, debe seleccionar primero un legajo")
		return 0
	}

	var fechaAlta time.Time
	sql := "SELECT fechaalta from legajo where id = " + strconv.Itoa(*legajoid)
	err := executor.db.Raw(sql).Row().Scan(&fechaAlta)

	if err != nil {
		fmt.Println("Error al buscar la fechaalta para el legajo " + strconv.Itoa(*legajoid))
		return 0
	}

	anioAlta := fechaAlta.Year()
	mesAlta := fechaAlta.Month()
	diaAlta := fechaAlta.Day()

	if anioPeriodoLiquidacion == anioAlta && mesPeriodoLiquidacion == mesAlta {
		respuesta = math.Max(float64(maximoDiaMesPeriodoLiquidacion - diaAlta) , 0)
	}

	if anioPeriodoLiquidacion > anioAlta || (anioPeriodoLiquidacion == anioAlta  && mesPeriodoLiquidacion > mesAlta) {
		respuesta = float64(maximoDiaMesPeriodoLiquidacion)
	}

	if anioPeriodoLiquidacion < anioAlta || (anioPeriodoLiquidacion == anioAlta  && mesPeriodoLiquidacion < mesAlta) {
		respuesta = 0
	}

	return respuesta
}

func (executor *Executor) CantidadMesesTrabajados() float64 {

	liquidacion := executor.context.Currentliquidacion

	fechaLiquidacion := liquidacion.Fecha //No puede ser null
	anioLiquidacion := fechaLiquidacion.Year()
	mesLiquidacion := fechaLiquidacion.Month()

	legajoid := liquidacion.Legajoid;
	if legajoid == nil {
		fmt.Println("Para realizar el calculo automatico de Sueldo, debe seleccionar primero un legajo")
		return 0
	}

	var fechaAlta time.Time
	sql := "SELECT fechaalta from legajo where id = " + strconv.Itoa(*legajoid)
	err := executor.db.Raw(sql).Row().Scan(&fechaAlta)

	anioAlta := fechaAlta.Year()
	mesAlta := fechaAlta.Month()

	if err != nil {
		fmt.Println("Error al buscar la fechaalta para el legajo " + strconv.Itoa(*legajoid))
		return 0
	}

	year, month, _, _, _, _ := diff(fechaLiquidacion, fechaAlta)

	mesesDiferencia := math.Max(float64(year * 12 + month), 0)

	if anioLiquidacion < anioAlta || (anioLiquidacion == anioAlta  && mesLiquidacion < mesAlta) {
		return 0
	}

	return mesesDiferencia
}

func (executor *Executor) MejorRemRemunerativaSemestre() float64 {
	return executor.calcularMejorRemunerativa(semestral, false)
}

func (executor *Executor) MejorRemNoRemunerativaSemestre() float64 {

	liquidacionActual := executor.context.Currentliquidacion
	mesliquidacion := liquidacionActual.Fechaperiodoliquidacion.Month()
	liquidaciones := executor.obtenerLiquidacionesIgualAnioLegajoMenorOIgualMesMensualesOQuincenales()
	var mejorRemuneracion float64 = 0
	mesInicial := devolverMesInicial(semestral, mesliquidacion)
	for mes := mesInicial; mes <= mesliquidacion; mes++ {
		var acumuladorMensual float64 = 0
		for _, liquidacion := range *liquidaciones {
			if liquidacion.Fechaperiodoliquidacion.Month() == mes {
				acumuladorMensual += calculoNoRemunerativos(liquidacion, false)
			}
		}
		if mejorRemuneracion < acumuladorMensual {
			mejorRemuneracion = acumuladorMensual
		}

	}
	return mejorRemuneracion
}

func (executor *Executor) MejorRemNormalYHabitualSemestre() float64 {
	return executor.calcularMejorRemunerativa(semestral, true)
}

func (executor *Executor) MejorRemTotalSinRemVarSemestre() float64 {
	return executor.calcularMejorTotal(semestral, true)
}

func (executor *Executor) MejorRemTotalSinRemVarAnual() float64 {
	return executor.calcularMejorTotal(anual, true)
}

func (executor *Executor) MejorRemTotalSemestre() float64 {
	return executor.calcularMejorTotal(semestral, false)
}

func (executor *Executor) MejorRemTotalAnual() float64 {
	return executor.calcularMejorTotal(anual, false)
}

func (executor *Executor) PromRemVariablesSemestre() float64 {
	return executor.calcularPromedioRemuneracionVariable(semestral)
}

func (executor *Executor) PromRemVariablesAnual() float64 {
	return executor.calcularPromedioRemuneracionVariable(anual)
}


//AUXILIARES

func (executor *Executor) calcularPromedioRemuneracionVariable(tipo int) float64 {

	liquidacionActual := executor.context.Currentliquidacion
	mesliquidacion := liquidacionActual.Fechaperiodoliquidacion.Month()
	liquidaciones := executor.obtenerLiquidacionesIgualAnioLegajoMenorOIgualMesMensualesOQuincenales()
	var total float64 = 0
	mesInicial := devolverMesInicial(tipo, mesliquidacion)
	var cantidadMeses float64 = 0
	for mes := mesInicial; mes <= mesliquidacion; mes++ {
		cantidadMeses++
		for _, liquidacion := range *liquidaciones {
			if liquidacion.Fechaperiodoliquidacion.Month() == mes {
				total += obtenerRemuneracionesVariables(liquidacion)
			}
		}
	}
	return total / cantidadMeses
}

func obtenerRemuneracionesVariables(liquidacion structLiquidacion.Liquidacion) float64 {
	var importeCalculado float64
	var importeNil *float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		if liquidacionitem.DeletedAt == nil {
			if liquidacionitem.Concepto.Esremvariable {
				if liquidacionitem.Importeunitario != importeNil {
					importeCalculado = importeCalculado + *liquidacionitem.Importeunitario
				}
			}
		}
	}

	return importeCalculado
}

func (executor *Executor) calcularMejorRemunerativa(tipo int, ignoraRemVariable bool) float64 {

	liquidacionActual := executor.context.Currentliquidacion
	mesliquidacion := liquidacionActual.Fechaperiodoliquidacion.Month()
	liquidaciones := executor.obtenerLiquidacionesIgualAnioLegajoMenorOIgualMesMensualesOQuincenales()
	var mejorRemuneracion float64 = 0
	mesInicial := devolverMesInicial(tipo, mesliquidacion)
	for mes := mesInicial; mes <= mesliquidacion; mes++ {
		var acumuladorMensual float64 = 0
		for _, liquidacion := range *liquidaciones {
			if liquidacion.Fechaperiodoliquidacion.Month() == mes {
				acumuladorMensual += calculoRemunerativos(liquidacion, ignoraRemVariable) - calculoDescuentos(liquidacion, ignoraRemVariable)
			}
		}
		if mejorRemuneracion < acumuladorMensual {
			mejorRemuneracion = acumuladorMensual
		}

	}
	return mejorRemuneracion
}

func (executor *Executor) calcularMejorTotal(tipo int, ignoraRemVariable bool) float64 {

	liquidacionActual := executor.context.Currentliquidacion
	mesliquidacion := liquidacionActual.Fechaperiodoliquidacion.Month()
	liquidaciones := executor.obtenerLiquidacionesIgualAnioLegajoMenorOIgualMesMensualesOQuincenales()
	var mejorRemuneracion float64 = 0
	mesInicial := devolverMesInicial(tipo, mesliquidacion)
	for mes := mesInicial; mes <= mesliquidacion; mes++ {
		var acumuladorMensual float64 = 0
		for _, liquidacion := range *liquidaciones {
			if liquidacion.Fechaperiodoliquidacion.Month() == mes {
				acumuladorMensual += calculoRemunerativos(liquidacion, ignoraRemVariable) - calculoDescuentos(liquidacion, ignoraRemVariable) + calculoNoRemunerativos(liquidacion, ignoraRemVariable)
			}
		}
		if mejorRemuneracion < acumuladorMensual {
			mejorRemuneracion = acumuladorMensual
		}

	}
	return mejorRemuneracion
}

func devolverMesInicial(tipo int, mesliquidacion time.Month) time.Month {
	if tipo == semestral {
		if mesliquidacion < time.July {
			return time.January
		}
		return time.July
	}
	return time.January
}

func calculoRemunerativos(liquidacion structLiquidacion.Liquidacion, ignoravariables bool) float64 {
	importeCalculado := calcularImporteSegunTipoConcepto(liquidacion, "IMPORTE_REMUNERATIVO", ignoravariables)

	return importeCalculado
}

func calculoDescuentos(liquidacion structLiquidacion.Liquidacion, ignoravariables bool) float64 {
	importeCalculado := calcularImporteSegunTipoConcepto(liquidacion, "DESCUENTO", ignoravariables)

	return importeCalculado
}

func calculoNoRemunerativos(liquidacion structLiquidacion.Liquidacion, ignoravariables bool) float64 {

	importeCalculado := calcularImporteSegunTipoConcepto(liquidacion, "IMPORTE_NO_REMUNERATIVO", ignoravariables)
	return importeCalculado
}


func (executor *Executor) obtenerLiquidacionesIgualAnioLegajoMenorOIgualMesMensualesOQuincenales() *[]structLiquidacion.Liquidacion {
	var liquidaciones []structLiquidacion.Liquidacion
	liquidacion := executor.context.Currentliquidacion
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesliquidacion := liquidacion.Fechaperiodoliquidacion.Month()
	executor.db.Order("to_number(to_char(fechaperiodoliquidacion, 'MM'),'99') desc").Set("gorm:auto_preload", true).Find(&liquidaciones, "to_number(to_char(fechaperiodoliquidacion, 'MM'),'99') < ? AND to_char(fechaperiodoliquidacion, 'YYYY') = ? AND legajoid = ? AND tipoid in (-1, -2, -3)", mesliquidacion, strconv.Itoa(anioperiodoliquidacion), *liquidacion.Legajoid)
	if *liquidacion.Tipoid == -1 || *liquidacion.Tipoid == -2 || *liquidacion.Tipoid == -3 {
		liquidaciones = append(liquidaciones, liquidacion)
	}
	return &liquidaciones
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

const (
	semestral = 1
	anual = 2
)

func calcularImporteSegunTipoConcepto(liquidacion structLiquidacion.Liquidacion, tipoConceptoCodigo string, ignoravariables bool) float64 {
	var importeCalculado float64
	var importeNil *float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		if liquidacionitem.DeletedAt == nil {
			if liquidacionitem.Concepto.Tipoconcepto.Codigo == tipoConceptoCodigo {
				if !ignoravariables || !liquidacionitem.Concepto.Esremvariable {
					if liquidacionitem.Importeunitario != importeNil {
						importeCalculado = importeCalculado + *liquidacionitem.Importeunitario
					}
				}
			}
		}
	}

	return importeCalculado
}