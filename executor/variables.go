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

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "IMPORTE_REMUNERATIVO")
}

func (executor *Executor) TotalHaberesNoRemunerativosMensual() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "IMPORTE_NO_REMUNERATIVO")
}

func (executor *Executor) TotalDescuentosMensual() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "DESCUENTO")
}

func (executor *Executor) TotalRetencionesMensual() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "RETENCION")
}

func (executor *Executor) TotalAportesPatronalesMensual() float64 {

	return calcularImporteSegunTipoConcepto(executor.context.Currentliquidacion, "APORTE_PATRONAL")
}

func calcularImporteSegunTipoConcepto(liquidacion structLiquidacion.Liquidacion, tipoConceptoCodigo string) float64 {
	var importeCalculado float64
	var importeNil *float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		if liquidacionitem.DeletedAt == nil {
			if liquidacionitem.Concepto.Tipoconcepto.Codigo == tipoConceptoCodigo {
				if liquidacionitem.Importeunitario != importeNil {
					importeCalculado = importeCalculado + *liquidacionitem.Importeunitario
				}
			}
		}
	}

	return importeCalculado
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