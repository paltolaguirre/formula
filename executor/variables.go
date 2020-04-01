package executor

import "strconv"

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
	legajo := executor.context.Currentliquidacion.Legajo;
	if legajo != nil {
		return float64(legajo.Remuneracion)
	} else {
		println("Para realizar el calculo automatico de Sueldo, debe seleccionar primero un legajo")
	}

	return 0
}

func (executor *Executor) HorasMensuales() float64 {
	legajo := executor.context.Currentliquidacion.Legajo;
	if legajo != nil {
		respuesta, err := strconv.ParseFloat(legajo.Horasmensualesnormales, 64)
		if err != nil {
			println("Las horas mensuales normales de ese legajo no son validas")
		} else {
			return respuesta
		}

	} else {
		println("Para realizar el calculo automatico de Sueldo, debe seleccionar primero un legajo")
	}

	return 0
}