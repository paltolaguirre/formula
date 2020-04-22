package executor

import (
	"math"
	"testing"
)

func TestSac(t *testing.T) {

	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(129411.11*100)/100
	respuesta := executor.Sac()

	if respuesta != esperado {
		t.Errorf("La funcion Sac con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestSacNoRemunerativo(t *testing.T) {

	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(12666.670000*100)/100
	respuesta := executor.SacNoRemunerativo()

	if respuesta != esperado {
		t.Errorf("La funcion SacNoRemunerativo con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestIntegracionMesDespido(t *testing.T) {

	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(386666.670000*100)/100
	respuesta := executor.IntegracionMesDespido()

	if respuesta != esperado {
		t.Errorf("La funcion IntegracionMesDespido con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}
}