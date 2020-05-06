package executor

import (
	"fmt"
	"os"
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/framework/configuracion"
	"math"
)

var DB *gorm.DB
var configuracionS configuracion.Configuracion

func TestMain(m *testing.M) {
	configuracionS = configuracion.GetInstance()

	os.Exit(m.Run())

}

func setupTest() {
	tenantPrueba := configuracionS.Tenanttest
	DB = conexionBD.ObtenerDB(tenantPrueba)
}

func afterTest() {
	conexionBD.CerrarDB(DB)
}

func getExecutorTest() Executor {

	var liquidacion structLiquidacion.Liquidacion
	//gorm:auto_preload se usa para que complete todos los struct con su informacion
	if err := DB.Set("gorm:auto_preload", true).First(&liquidacion, "id = ?", 1).Error; gorm.IsRecordNotFoundError(err) {
		fmt.Println("No se pudo cargar la liquidacion 1")
		return Executor{}
	}
	contexto := Context{Currentliquidacion: liquidacion}

	return Executor{context: &contexto, db: DB}
}

func TestSueldo(t *testing.T) {

	setupTest()
	defer afterTest()
	defer afterTest()
	executor := getExecutorTest()

	esperado := float64(400000)
	respuesta := executor.Sueldo()

	if respuesta != esperado {
		t.Errorf("La funcion Sueldo devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestHorasMensuales(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	respuesta := executor.HorasMensuales()

	expected := float64(200)

	if respuesta != expected {
		t.Errorf("La funcion HorasMensuales devuelve %f y se esperaba %f", respuesta, expected)
	}
}

func TestDiasMesTrabajadosFechaLiquidacion(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaTest())

	respuesta := executor.DiasMesTrabajadosFechaLiquidacion()

	expected := float64(0)

	if respuesta != expected {
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionAntesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionDespuesDeAltaTest())

	respuesta = executor.DiasMesTrabajadosFechaLiquidacion()

	expected = float64(getFechaLiquidacionDespuesDeAltaTest().Day())

	if respuesta != expected {
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionDespuesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaMismoMesTest())

	respuesta = executor.DiasMesTrabajadosFechaLiquidacion()

	expected = float64(getFechaLiquidacionAntesDeAltaMismoMesTest().Day() - getFechaAltaTest().Day())

	if respuesta != expected {
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionAntesDeAltaMismoMesTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}
}

func setFechaLiquidacion(executor *Executor, fechaLiquidacion time.Time) {
	executor.context.Currentliquidacion.Fecha = fechaLiquidacion
}

func setFechaPeriodoLiquidacion(executor *Executor, fechaPeriodoLiquidacion time.Time) {
	executor.context.Currentliquidacion.Fechaperiodoliquidacion = fechaPeriodoLiquidacion
}

func getFechaPeriodoLiquidacionDespuesDeAltaTest() time.Time {
	return getFechaLiquidacionDespuesDeAltaTest()
}

func getFechaPeriodoLiquidacionAntesDeAltaTest() time.Time {
	return getFechaLiquidacionAntesDeAltaTest()
}

func getFechaPeriodoLiquidacionAntesDeAltaMismoMesTest() time.Time {
	return getFechaLiquidacionAntesDeAltaMismoMesTest()
}

func getFechaAltaTest() *time.Time {

	fecha, err := time.Parse("2006-01-02", "2019-01-14")

	if err != nil {
		fmt.Println("getFechaAltaTest mal creado ", err)
	}

	return &fecha
}

func getFechaLiquidacionDespuesDeAltaTest() time.Time {

	fecha, err := time.Parse("2006-01-02", "2020-01-30")

	if err != nil {
		fmt.Println("getFechaLiquidacionDespuesDeAltaTest mal creado ", err)
	}

	return fecha
}

func getFechaLiquidacionAntesDeAltaTest() time.Time {

	fecha, err := time.Parse("2006-01-02", "2018-01-01")

	if err != nil {
		fmt.Println("getFechaLiquidacionAntesDeAltaTest mal creado ", err)
	}

	return fecha
}

func getFechaLiquidacionAntesDeAltaMismoMesTest() time.Time {
	fecha, err := time.Parse("2006-01-02", "2019-01-21")

	if err != nil {
		fmt.Println("getFechaLiquidacionAntesDeAltaMismoMesTest mal creado ", err)
	}

	return fecha
}

func TestDiasMesTrabajadosFechaPeriodo(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getFechaPeriodoLiquidacionAntesDeAltaTest())

	respuesta := executor.DiasMesTrabajadosFechaPeriodo()

	expected := float64(0)

	if respuesta != expected {
		t.Errorf("La funcion DiasMesTrabajadosFechaPeriodo con getFechaPeriodoLiquidacionAntesDeAltaTest devuelve %f y se esperaba %f", respuesta, expected)
	}

	setFechaPeriodoLiquidacion(&executor, getFechaPeriodoLiquidacionDespuesDeAltaTest())

	respuesta = executor.DiasMesTrabajadosFechaPeriodo()

	expected = float64(now.New(getFechaPeriodoLiquidacionDespuesDeAltaTest()).EndOfMonth().Day())

	if respuesta != expected {
		t.Errorf("La funcion DiasMesTrabajadosFechaPeriodo con getFechaPeriodoLiquidacionDespuesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaPeriodoLiquidacion(&executor, getFechaPeriodoLiquidacionAntesDeAltaMismoMesTest())

	respuesta = executor.DiasMesTrabajadosFechaPeriodo()

	expected = float64(now.New(getFechaPeriodoLiquidacionAntesDeAltaMismoMesTest()).EndOfMonth().Day() - getFechaAltaTest().Day())

	if respuesta != expected {
		t.Errorf("La funcion DiasMesTrabajadosFechaPeriodo con getFechaPeriodoLiquidacionAntesDeAltaMismoMesTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}
}

func TestTotalHaberesNoRemunerativosMensual(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	esperado := float64(30000)
	respuesta := executor.TotalHaberesNoRemunerativosMensual()

	if respuesta != esperado {
		t.Errorf("La funcion Sueldo devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestTotalImporteRemunerativo(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	esperado := float64(120000)
	respuesta := executor.TotalImporteRemunerativo()

	if respuesta != esperado {
		t.Errorf("La funcion TotalImporteRemunerativo devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestTotalDescuentosMensual(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	esperado := float64(7500)
	respuesta := executor.TotalDescuentosMensual()

	if respuesta != esperado {
		t.Errorf("La funcion TotalDescuentosMensual devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestTotalRetencionesMensual(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	esperado := float64(14625)
	respuesta := executor.TotalRetencionesMensual()

	if respuesta != esperado {
		t.Errorf("La funcion TotalRetencionesMensual devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestTotalAportesPatronalesMensual(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	esperado := float64(3000)
	respuesta := executor.TotalAportesPatronalesMensual()

	if respuesta != esperado {
		t.Errorf("La funcion TotalAportesPatronalesMensual devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestValorDiasVacaciones(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	esperado := float64(16000)
	respuesta := executor.ValorDiasVacaciones()

	if respuesta != esperado {
		t.Errorf("La funcion ValorDiasVacaciones devuelve %f y se esperaba %f", respuesta, esperado)
	}
}

func TestCantidadMesesTrabajados(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaMismoMesTest())

	esperado := float64(1)
	respuesta := executor.CantidadMesesTrabajados()

	if respuesta != esperado {
		t.Errorf("La funcion CantidadMesesTrabajados con getFechaLiquidacionAntesDeAltaMismoMesTest devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaTest())

	esperado = float64(0)
	respuesta = executor.CantidadMesesTrabajados()

	if respuesta != esperado {
		t.Errorf("La funcion CantidadMesesTrabajados con getFechaLiquidacionAntesDeAltaTest devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionDespuesDeAltaTest())

	esperado = float64(13)
	respuesta = executor.CantidadMesesTrabajados()

	if respuesta != esperado {
		t.Errorf("La funcion CantidadMesesTrabajados con getFechaLiquidacionAntesDeAltaTest devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func getPeriodoLiquidacionMayo2020() time.Time {

	fecha, err := time.Parse("2006-01-02", "2020-05-01")

	if err != nil {
		fmt.Println("getPeriodoLiquidacionMayo2020 mal creado ", err)
	}

	return fecha
}

func getPeriodoLiquidacionAbril2020() time.Time {
	fecha, err := time.Parse("2006-01-02", "2020-04-01")

	if err != nil {
		fmt.Println("getPeriodoLiquidacionMayo2020 mal creado ", err)
	}

	return fecha
}


func getPeriodoLiquidacionMarzo2020() time.Time {
	fecha, err := time.Parse("2006-01-02", "2020-03-01")

	if err != nil {
		fmt.Println("getPeriodoLiquidacionMayo2020 mal creado ", err)
	}

	return fecha
}

func getPeriodoLiquidacionFebrero2020() time.Time {
	fecha, err := time.Parse("2006-01-02", "2020-02-01")

	if err != nil {
		fmt.Println("getPeriodoLiquidacionMayo2020 mal creado ", err)
	}

	return fecha
}

func getPeriodoLiquidacionEnero2020() time.Time {
	fecha, err := time.Parse("2006-01-02", "2020-01-01")

	if err != nil {
		fmt.Println("getPeriodoLiquidacionMayo2020 mal creado ", err)
	}

	return fecha
}


func getFechaLiquidacionEnero2019() time.Time {
	fecha, err := time.Parse("2006-01-02", "2019-01-01")

	if err != nil {
		fmt.Println("getFechaLiquidacionEnero2019 mal creado ", err)
	}

	return fecha
}

func getFechaLiquidacionSeptiembre2019() time.Time {
	fecha, err := time.Parse("2006-01-02", "2019-09-01")

	if err != nil {
		fmt.Println("getPeriodoLiquidacionSeptiembre2019 mal creado ", err)
	}

	return fecha
}

func getFechaLiquidacionEnero2021() time.Time {
	fecha, err := time.Parse("2006-01-02", "2021-01-01")

	if err != nil {
		fmt.Println("getPeriodoLiquidacionEnero2021 mal creado ", err)
	}

	return fecha
}

func getFechaLiquidacionSeptiembre2021() time.Time {
	fecha, err := time.Parse("2006-01-02", "2021-09-01")

	if err != nil {
		fmt.Println("getFechaLiquidacionSeptiembre2021 mal creado ", err)
	}

	return fecha
}

func getFechaLiquidacionEnero2025() time.Time {
	fecha, err := time.Parse("2006-01-02", "2025-01-01")

	if err != nil {
		fmt.Println("getFechaLiquidacionEnero2025 mal creado ", err)
	}

	return fecha
}

func getFechaLiquidacionSeptiembre2025() time.Time {
	fecha, err := time.Parse("2006-01-02", "2025-09-01")

	if err != nil {
		fmt.Println("getFechaLiquidacionSeptiembre2025 mal creado ", err)
	}

	return fecha
}

func getPeriodoLiquidacionJulio2020() time.Time {
	fecha, err := time.Parse("2006-01-02", "2020-07-01")

	if err != nil {
		fmt.Println("getPeriodoLiquidacionJulio2020 mal creado ", err)
	}

	return fecha
}

func TestMejorRemRemunerativaSemestre(t *testing.T) {
	setupTest()
	defer afterTest()

	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(306500)
	respuesta := executor.MejorRemRemunerativaSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion MejorRemRemunerativaSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestMejorRemNoRemunerativaSemestre(t *testing.T) {
	setupTest()
	defer afterTest()

	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(30000)
	respuesta := executor.MejorRemNoRemunerativaSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion MejorRemNoRemunerativaSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestMejorRemNormalYHabitualSemestre(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(296500)
	respuesta := executor.MejorRemNormalYHabitualSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion MejorRemNormalYHabitualSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestMejorRemTotalSinRemVarSemestre(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(296500)
	respuesta := executor.MejorRemTotalSinRemVarSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion MejorRemTotalSinRemVarSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestMejorRemTotalSinRemVarAnual(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(296500)
	respuesta := executor.MejorRemTotalSinRemVarAnual()

	if respuesta != esperado {
		t.Errorf("La funcion MejorRemTotalSinRemVarAnual con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestMejorRemTotalSemestre(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(306500)
	respuesta := executor.MejorRemTotalSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion MejorRemTotalSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestMejorRemTotalAnual(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(306500)
	respuesta := executor.MejorRemTotalAnual()

	if respuesta != esperado {
		t.Errorf("La funcion MejorRemTotalAnual con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestPromRemVariablesSemestre(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(6000)
	respuesta := executor.PromRemVariablesSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion PromRemVariablesSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionFebrero2020())

	esperado = float64(15000)
	respuesta = executor.PromRemVariablesSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion PromRemVariablesSemestre con getPeriodoLiquidacionFebrero2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionEnero2020())

	esperado = float64(20000)
	respuesta = executor.PromRemVariablesSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion PromRemVariablesSemestre con getPeriodoLiquidacionEnero2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}


	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMarzo2020())

	esperado = float64(10000)
	respuesta = executor.PromRemVariablesSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion PromRemVariablesSemestre con getPeriodoLiquidacionMarzo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionAbril2020())

	esperado = float64(7500)
	respuesta = executor.PromRemVariablesSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion PromRemVariablesSemestre con getPeriodoLiquidacionMarzo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestPromRemVariablesAnual(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(6000)
	respuesta := executor.PromRemVariablesAnual()

	if respuesta != esperado {
		t.Errorf("La funcion PromRemVariablesAnual con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestDiasSemTrabajados(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(152)
	respuesta := executor.DiasSemTrabajados()

	if respuesta != esperado {
		t.Errorf("La funcion DiasSemTrabajados con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}


	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionJulio2020())

	esperado = float64(31)
	respuesta = executor.DiasSemTrabajados()

	if respuesta != esperado {
		t.Errorf("La funcion DiasSemTrabajados con getPeriodoLiquidacionJulio2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestDiasEfectivamenteTrabajadosSemestre(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(140)
	respuesta := executor.DiasEfectivamenteTrabajadosSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion DiasEfectivamenteTrabajadosSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestAntiguedadResto(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := math.Round(0.37808*100)/100
	respuesta := executor.AntiguedadResto()

	if respuesta != esperado {
		t.Errorf("La funcion AntiguedadResto con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestFechaDeIngreso(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	esperado := time.Date(2019, 1, 14, 3, 0, 0, 0, time.UTC)
	respuesta := executor.FechaDeIngreso()

	if respuesta != esperado {
		t.Errorf("La funcion FechaDeIngreso devuelve %s y se esperaba %s", respuesta, esperado)
	}

}

func TestFechadeLiquidacion(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaMismoMesTest() )

	esperado := getFechaLiquidacionAntesDeAltaMismoMesTest()
	respuesta := executor.FechadeLiquidacion()

	if respuesta != esperado {
		t.Errorf("La funcion FechaDeIngreso devuelve %s y se esperaba %s", respuesta, esperado)
	}

}

func TestFecIngHASTAFecLiq(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaLiquidacion(&executor, getFechaLiquidacionDespuesDeAltaTest() )

	esperado := float64(1.04)
	respuesta := executor.FecIngHASTAFecLiq()

	if respuesta != esperado {
		t.Errorf("La funcion FecIngHASTAFecLiq devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestMejorRemRemunerativaBaseSACSemestre(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(306500)
	respuesta := executor.MejorRemRemunerativaBaseSACSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion MejorRemRemunerativaSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}


func TestDiasLicenciaMes(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(2)
	respuesta := executor.DiasLicenciaMes()

	if respuesta != esperado {
		t.Errorf("La funcion DiasLicenciaMes con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestDiasLicenciaSemestre(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getPeriodoLiquidacionMayo2020())

	esperado := float64(2)
	respuesta := executor.DiasLicenciaSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion DiasLicenciaSemestre con getPeriodoLiquidacionMayo2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func TestDiasDelSemestre(t *testing.T) {
	setupTest()
	defer afterTest()
	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getsetFechaPeriodoLiquidacionJunio2020())

	esperado := float64(182)
	respuesta := executor.DiasDelSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion DiasLicenciaSemestre con getsetFechaPeriodoLiquidacionJunio2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

	setFechaPeriodoLiquidacion(&executor, getsetFechaPeriodoLiquidacionDiciembre2020())

	esperado = float64(184)
	respuesta = executor.DiasDelSemestre()

	if respuesta != esperado {
		t.Errorf("La funcion DiasLicenciaSemestre con getsetFechaPeriodoLiquidacionDiciembre2020 devuelve %f y se esperaba %f", respuesta, esperado)
	}

}

func getsetFechaPeriodoLiquidacionJunio2020() time.Time {
	fecha, err := time.Parse("2006-01-02", "2020-06-01")

	if err != nil {
		fmt.Println("getFechaLiquidacionJunio2020 mal creado ", err)
	}

	return fecha
}

func getsetFechaPeriodoLiquidacionDiciembre2020() time.Time {
	fecha, err := time.Parse("2006-01-02", "2020-12-01")

	if err != nil {
		fmt.Println("getFechaLiquidacionDiciembre2020 mal creado ", err)
	}

	return fecha
}