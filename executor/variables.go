package executor

import (
	"fmt"
	"github.com/jinzhu/now"
	"math"
	"strconv"
	"time"
)

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

	legajo := liquidacion.Legajo;
	if legajo == nil {
		println("Para realizar el calculo automatico de Sueldo, debe seleccionar primero un legajo")
		return 0
	}

	fechaAlta := legajo.Fechaalta //No puede ser null

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