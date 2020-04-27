package executor

import (
	"math"
	"testing"
)

func TestSac(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(129411.11*100)/100
	respuesta := executor.Sac()

	if respuesta != esperado {
		t.Errorf("La funcion Sac con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestSacNoRemunerativo(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(12666.670000*100)/100
	respuesta := executor.SacNoRemunerativo()

	if respuesta != esperado {
		t.Errorf("La funcion SacNoRemunerativo con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestIntegracionMesDespido(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(386666.570000*100)/100
	respuesta := executor.IntegracionMesDespido()

	if respuesta != esperado {
		t.Errorf("La funcion IntegracionMesDespido con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestSacSinImd(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(32222.210000*100)/100
	respuesta := executor.SacSinImd()

	if respuesta != esperado {
		t.Errorf("La funcion SacSinImd con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestHoraExtra50(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(6000*100)/100
	respuesta := executor.HoraExtra50()

	if respuesta != esperado {
		t.Errorf("La funcion HoraExtra50 con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestPreaviso(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaLiquidacion(&executor, getFechaLiquidacionEnero2019())

	esperado := executor.Sueldo()/2
	respuesta := executor.Preaviso()

	if respuesta != esperado {
		t.Errorf("La funcion Preaviso con getFechaLiquidacionEnero2019 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionSeptiembre2019())

	esperado = executor.Sueldo()
	respuesta = executor.Preaviso()

	if respuesta != esperado {
		t.Errorf("La funcion Preaviso con getFechaLiquidacionSeptiembre2019 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionEnero2021())

	esperado = executor.Sueldo()
	respuesta = executor.Preaviso()

	if respuesta != esperado {
		t.Errorf("La funcion Preaviso con getFechaLiquidacionEnero2021 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionSeptiembre2021())

	esperado = executor.Sueldo()
	respuesta = executor.Preaviso()

	if respuesta != esperado {
		t.Errorf("La funcion Preaviso con getFechaLiquidacionSeptiembre2021 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionSeptiembre2025())

	esperado = executor.Sueldo() * 2
	respuesta = executor.Preaviso()

	if respuesta != esperado {
		t.Errorf("La funcion Preaviso con getFechaLiquidacionSeptiembre2025 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionEnero2025())

	esperado = executor.Sueldo() * 2
	respuesta = executor.Preaviso()

	if respuesta != esperado {
		t.Errorf("La funcion Preaviso con getFechaLiquidacionEnero2025 devuelve %f y se esperaba %f", respuesta, esperado)
	}
}