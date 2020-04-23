package executor

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Legajo/structLegajo"
	"math"
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
)

func (executor *Executor) TotalImporteRemunerativo() float64 {

	return executor.calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "IMPORTE_REMUNERATIVO", false, false)
}

func (executor *Executor) TotalHaberesNoRemunerativosMensual() float64 {

	return executor.calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "IMPORTE_NO_REMUNERATIVO", false, false)
}

func (executor *Executor) TotalDescuentosMensual() float64 {

	return executor.calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "DESCUENTO", false, false)
}

func (executor *Executor) TotalRetencionesMensual() float64 {

	return executor.calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "RETENCION", false, false)
}

func (executor *Executor) TotalAportesPatronalesMensual() float64 {

	return executor.calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "APORTE_PATRONAL", false, false)
}

func (executor *Executor) Sueldo() float64 {
	var importe float64
	legajoid := executor.context.Currentliquidacion.Legajoid
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
	legajoid := executor.context.Currentliquidacion.Legajoid
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

	fechaLiquidacion := executor.FechadeLiquidacion()

	anioLiquidacion := fechaLiquidacion.Year()
	mesLiquidacion := fechaLiquidacion.Month()
	diaLiquidacion := fechaLiquidacion.Day()

	fechaAlta := executor.FechaDeIngreso()

	anioAlta := fechaAlta.Year()
	mesAlta := fechaAlta.Month()
	diaAlta := fechaAlta.Day()

	if anioLiquidacion == anioAlta && mesLiquidacion == mesAlta {
		respuesta = math.Max(float64(diaLiquidacion-diaAlta), 0)
	}

	if anioLiquidacion > anioAlta || (anioLiquidacion == anioAlta && mesLiquidacion > mesAlta) {
		respuesta = float64(diaLiquidacion)
	}

	if anioLiquidacion < anioAlta || (anioLiquidacion == anioAlta && mesLiquidacion < mesAlta) {
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

	fechaAlta := executor.FechaDeIngreso()

	anioAlta := fechaAlta.Year()
	mesAlta := fechaAlta.Month()
	diaAlta := fechaAlta.Day()

	if anioPeriodoLiquidacion == anioAlta && mesPeriodoLiquidacion == mesAlta {
		respuesta = math.Max(float64(maximoDiaMesPeriodoLiquidacion-diaAlta), 0)
	}

	if anioPeriodoLiquidacion > anioAlta || (anioPeriodoLiquidacion == anioAlta && mesPeriodoLiquidacion > mesAlta) {
		respuesta = float64(maximoDiaMesPeriodoLiquidacion)
	}

	if anioPeriodoLiquidacion < anioAlta || (anioPeriodoLiquidacion == anioAlta && mesPeriodoLiquidacion < mesAlta) {
		respuesta = 0
	}

	return respuesta
}

func (executor *Executor) CantidadMesesTrabajados() float64 {

	liquidacion := executor.context.Currentliquidacion

	fechaLiquidacion := liquidacion.Fecha //No puede ser null
	anioLiquidacion := fechaLiquidacion.Year()
	mesLiquidacion := fechaLiquidacion.Month()

	fechaAlta := executor.FechaDeIngreso()
	anioAlta := fechaAlta.Year()
	mesAlta := fechaAlta.Month()

	year, month, _, _, _, _ := diff(fechaLiquidacion, fechaAlta)

	mesesDiferencia := math.Max(float64(year*12+month), 0)

	if anioLiquidacion < anioAlta || (anioLiquidacion == anioAlta && mesLiquidacion < mesAlta) {
		return 0
	}

	return mesesDiferencia + 1
}

func (executor *Executor) MejorRemRemunerativaSemestre() float64 {
	return executor.calcularMejorRemunerativa(semestral, false, false)
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
				acumuladorMensual += executor.calculoNoRemunerativos(liquidacion, false)
			}
		}
		if mejorRemuneracion < acumuladorMensual {
			mejorRemuneracion = acumuladorMensual
		}

	}
	return mejorRemuneracion
}

func (executor *Executor) MejorRemNormalYHabitualSemestre() float64 {
	return executor.calcularMejorRemunerativa(semestral, true, false)
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

func (executor *Executor) DiasLicenciaMes() float64 {
	return executor.obtenerDiasLicencia(mensual)
}

func (executor *Executor) DiasLicenciaSemestre() float64 {
	return executor.obtenerDiasLicencia(semestral)
}

func (executor *Executor) DiasSemTrabajados() float64 {

	liquidacionActual := executor.context.Currentliquidacion
	fechaPeriodoliquidacion := liquidacionActual.Fechaperiodoliquidacion

	ultimoDiaDelMesPeriodoLiquidacion := now.New(fechaPeriodoliquidacion).EndOfMonth()
	primerDiaDelSemestreLiquidacion := now.New(fechaPeriodoliquidacion).BeginningOfYear()

	if fechaPeriodoliquidacion.Month() >= time.July {
		primerDiaDelSemestreLiquidacion = now.New(primerDiaDelSemestreLiquidacion).AddDate(0, 6, 0)
	}

	fechaAlta := executor.FechaDeIngreso()

	if fechaAlta.Before(primerDiaDelSemestreLiquidacion) {
		return diffDias(ultimoDiaDelMesPeriodoLiquidacion, primerDiaDelSemestreLiquidacion)
	} else {
		return diffDias(ultimoDiaDelMesPeriodoLiquidacion, fechaAlta)
	}
}

func (executor *Executor) DiasEfectivamenteTrabajadosSemestre() float64 {
	return executor.DiasSemTrabajados() - executor.obtenerDiasSegunConcepto(semestral, diasAccidente) - executor.obtenerDiasSegunConcepto(semestral, diasEnfermedad) - executor.obtenerDiasSegunConcepto(semestral, diasFaltasInjustificadas) - executor.obtenerDiasSegunConcepto(semestral, diasLicencia)
}

const (
	diasAccidente = -16
	diasEnfermedad = -15
	diasFaltasInjustificadas = -17
	diasLicencia = -34
)

func (executor *Executor) AntiguedadResto() float64 {

	liquidacion := executor.context.Currentliquidacion

	start := executor.FechaDeIngreso()
	end := now.New(liquidacion.Fechaperiodoliquidacion).EndOfMonth()

	period := end.Sub(start)
	days := int(period.Hours() / 24)
	yearsRemainder := float64(days%365) / float64(365)

	return math.Round(yearsRemainder*100)/100
}

func (executor *Executor) MejorRemRemunerativaBaseSACSemestre() float64 {
	return executor.calcularMejorRemunerativa(semestral, false, true)
}



//AUXILIARES

func (executor *Executor) obtenerDiasLicencia(tipo int) float64 {
	return executor.obtenerDiasSegunConcepto(tipo, diasLicencia)
}

func (executor *Executor) obtenerDiasSegunConcepto(tipo int, conceptoid int) float64 {
	liquidacionActual := executor.context.Currentliquidacion
	var resultado float64
	periodoLiquidacion := liquidacionActual.Fechaperiodoliquidacion //No puede ser null

	anioPeriodoLiquidacion := periodoLiquidacion.Year()
	mesPeriodoLiquidacion := periodoLiquidacion.Month()

	mesInicial := devolverMesInicial(tipo, mesPeriodoLiquidacion)

	sql := fmt.Sprintf("select coalesce(sum(coalesce(cantidad, 0)), 0) from liquidacionitem left join liquidacion on liquidacionitem.liquidacionid = liquidacion.id where extract(MONTH from fecha) < %d and extract(MONTH from fecha) >= %d and extract(YEAR from fecha) = %d and legajoid = %d and conceptoid in (%d) and liquidacionid != %d", mesPeriodoLiquidacion, mesInicial, anioPeriodoLiquidacion, *liquidacionActual.Legajoid, conceptoid, liquidacionActual.ID)
	err := executor.db.Raw(sql).Row().Scan(&resultado)

	if err != nil {
		fmt.Printf("No se pudo consultar la cantidad de %d", conceptoid)
		fmt.Println()
	}

	return resultado + executor.getDiasLiquidacionActualSegunConcepto(conceptoid)
}

func (executor *Executor) getDiasLiquidacionActualSegunConcepto(conceptoid int) float64 {
	cantidadTotal := 0.0
	for _, liquidacionItem := range executor.context.Currentliquidacion.Liquidacionitems {
		if *liquidacionItem.Conceptoid == conceptoid {
			cantidadTotal += liquidacionItem.Cantidad
		}
	}

	return cantidadTotal
}

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
				total += executor.obtenerRemuneracionesVariables(liquidacion)
			}
		}
	}
	return math.Round(total / cantidadMeses * 100) / 100
}

func (executor *Executor) obtenerRemuneracionesVariables(liquidacion structLiquidacion.Liquidacion) float64 {
	var importeCalculado float64
	var importeNil *float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		if liquidacionitem.DeletedAt == nil && *liquidacionitem.Conceptoid != executor.context.Currentconcepto.ID {
			if liquidacionitem.Concepto.Esremvariable {
				if liquidacionitem.Importeunitario != importeNil {
					importeCalculado = importeCalculado + *liquidacionitem.Importeunitario
				}
			}
		}
	}

	return importeCalculado
}

func (executor *Executor) calcularMejorRemunerativa(tipo int, ignoraRemVariable bool, soloBaseSac bool) float64 {

	liquidacionActual := executor.context.Currentliquidacion
	mesliquidacion := liquidacionActual.Fechaperiodoliquidacion.Month()
	liquidaciones := executor.obtenerLiquidacionesIgualAnioLegajoMenorOIgualMesMensualesOQuincenales()
	var mejorRemuneracion float64 = 0
	mesInicial := devolverMesInicial(tipo, mesliquidacion)
	for mes := mesInicial; mes <= mesliquidacion; mes++ {
		var acumuladorMensual float64 = 0
		for _, liquidacion := range *liquidaciones {
			if liquidacion.Fechaperiodoliquidacion.Month() == mes {
				acumuladorMensual += executor.calculoRemunerativos(liquidacion, ignoraRemVariable, soloBaseSac) - executor.calculoDescuentos(liquidacion, ignoraRemVariable, soloBaseSac)
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
				acumuladorMensual += executor.calculoRemunerativos(liquidacion, ignoraRemVariable, false) - executor.calculoDescuentos(liquidacion, ignoraRemVariable, false) + executor.calculoNoRemunerativos(liquidacion, ignoraRemVariable)
			}
		}
		if mejorRemuneracion < acumuladorMensual {
			mejorRemuneracion = acumuladorMensual
		}

	}
	return mejorRemuneracion
}

func devolverMesInicial(tipo int, mesliquidacion time.Month) time.Month {
	if tipo == mensual {
		return mesliquidacion
	}
	if tipo == semestral {
		if mesliquidacion < time.July {
			return time.January
		}
		return time.July
	}
	return time.January
}

func (executor *Executor) calculoRemunerativos(liquidacion structLiquidacion.Liquidacion, ignoravariables bool, soloBaseSac bool) float64 {
	importeCalculado := executor.calcularImporteSegunTipoConcepto(liquidacion, "IMPORTE_REMUNERATIVO", ignoravariables, soloBaseSac)

	return importeCalculado
}

func (executor *Executor) calculoDescuentos(liquidacion structLiquidacion.Liquidacion, ignoravariables bool, soloBaseSac bool) float64 {
	importeCalculado := executor.calcularImporteSegunTipoConcepto(liquidacion, "DESCUENTO", ignoravariables, soloBaseSac)

	return importeCalculado
}

func (executor *Executor) calculoNoRemunerativos(liquidacion structLiquidacion.Liquidacion, ignoravariables bool) float64 {

	importeCalculado := executor.calcularImporteSegunTipoConcepto(liquidacion, "IMPORTE_NO_REMUNERATIVO", ignoravariables, false)
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

func diffDias(a time.Time, b time.Time) float64 {
	days := math.Round(a.Sub(b).Hours() / 24)
	return days
}

const (
	mensual   = 1
	semestral = 2
	anual     = 3
)

func (executor *Executor) calcularImporteSegunTipoConcepto(liquidacion structLiquidacion.Liquidacion, tipoConceptoCodigo string, ignoravariables bool, soloBaseSac bool) float64 {
	var importeCalculado float64
	var importeNil *float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		if liquidacionitem.DeletedAt == nil && *liquidacionitem.Conceptoid != executor.context.Currentconcepto.ID{
			if liquidacionitem.Concepto.Tipoconcepto.Codigo == tipoConceptoCodigo {
				if !soloBaseSac || liquidacionitem.Concepto.Basesac {
					if !ignoravariables || !liquidacionitem.Concepto.Esremvariable {
						if liquidacionitem.Importeunitario != importeNil {
							importeCalculado = importeCalculado + *liquidacionitem.Importeunitario
						}
					}
				}
			}
		}
	}

	return importeCalculado
}


func (executor *Executor) FechaDeIngreso() time.Time {
	liquidacion := executor.context.Currentliquidacion

	var legajo structLegajo.Legajo

	if err := executor.db.Set("gorm:auto_preload", true).First(&legajo, "id = ?", liquidacion.Legajoid).Error; gorm.IsRecordNotFoundError(err) {
		t := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
		fmt.Println("Error al buscar la fechaalta para el legajo " + strconv.Itoa(*liquidacion.Legajoid))
		return t
	}

	return *legajo.Fechaalta
}

func (executor *Executor) FechadeLiquidacion() time.Time {
	liquidacion := executor.context.Currentliquidacion

	return liquidacion.Fecha
}


func (executor *Executor) FecIngHASTAFecLiq() float64 {
	return math.Round(((executor.FechadeLiquidacion().Sub(executor.FechaDeIngreso()).Hours() / 24) / 365) * 100 ) / 100
}