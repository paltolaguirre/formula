package executor

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/conexionBD/structGormModel"
	"os"
	"testing"
	"time"
)

var DB *gorm.DB

func TestMain(m *testing.M) {
	tenantPrueba := "tnt_143124"
	DB = conexionBD.ObtenerDB(tenantPrueba);
	defer conexionBD.CerrarDB(DB);
	os.Exit(m.Run())

}

func getExecutorTest() Executor {
	legajoid := 1
	liquidacion := structLiquidacion.Liquidacion{
		GormModel:                            structGormModel.GormModel{},
		Nombre:                               "",
		Codigo:                               "",
		Descripcion:                          "",
		Activo:                               0,
		Legajo:                               nil,
		Legajoid:                             &legajoid,
		Tipo:                                 nil,
		Tipoid:                               nil,
		Fecha:                                time.Time{},
		Fechaultimodepositoaportejubilatorio: time.Time{},
		Zonatrabajo:                          "",
		Condicionpago:                        nil,
		Condicionpagoid:                      nil,
		Cuentabancoid:                        nil,
		Cuentabanco:                          nil,
		Bancoaportejubilatorioid:             nil,
		Bancoaportejubilatorio:               nil,
		Fechaperiododepositado:               time.Time{},
		Fechaperiodoliquidacion:              time.Time{},
		Estacontabilizada:                    false,
		Asientomanualtransaccionid:           0,
		Asientomanualnombre:                  "",
		Liquidacionitems:                     nil,
	}
	contexto := Context{Currentliquidacion:liquidacion}

	return Executor{context:&contexto, db:DB}
}

func TestSueldo(t *testing.T) {

	executor := getExecutorTest()

	esperado := float64(400000)
	respuesta := executor.Sueldo()

	if respuesta != esperado {
		t.Errorf("La funcion Sueldo devuelve %f; want y se esperaba %f", respuesta, esperado)
	}
}

func TestHorasMensuales(t *testing.T) {

	executor := getExecutorTest()

	respuesta := executor.HorasMensuales()

	expected := float64(200)

	if respuesta != expected {
		t.Errorf("La funcion HorasMensuales devuelve %f; want y se esperaba %f", respuesta, expected)
	}
}

func TestDiasMesTrabajadosFechaLiquidacion(t *testing.T) {

	executor := getExecutorTest()

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaTest())

	respuesta := executor.DiasMesTrabajadosFechaLiquidacion()

	expected := float64(0)

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionAntesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionDespuesDeAltaTest())

	respuesta = executor.DiasMesTrabajadosFechaLiquidacion()

	expected = float64(getFechaLiquidacionDespuesDeAltaTest().Day())

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionDespuesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaMismoMesTest())

	respuesta = executor.DiasMesTrabajadosFechaLiquidacion()

	expected = float64(getFechaLiquidacionAntesDeAltaMismoMesTest().Day() - getFechaAltaTest().Day())

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionAntesDeAltaMismoMesTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}
}

func setFechaLiquidacion(executor *Executor, fechaLiquidacion time.Time)  {
	executor.context.Currentliquidacion.Fecha = fechaLiquidacion
}

func setFechaPeriodoLiquidacion(executor *Executor, fechaPeriodoLiquidacion time.Time)  {
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

/*
import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/xubiosueldos/conexionBD/Legajo/structLegajo"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/conexionBD/structGormModel"
	"strconv"
	"testing"
	"time"
)

func getRemuneracionTest() float32 {
	return 3000
}

func getHorasMensualesNormalesTest() string {
	return "200"
}




func TestSueldo(t *testing.T) {

	executor := getExecutorTest()

	respuesta := executor.Sueldo()

	if respuesta != float64(getRemuneracionTest()) {
		t.Errorf("La funcion Sueldo devuelve %f; want y se esperaba %f", respuesta, getRemuneracionTest())
	}
}

func TestHorasMensuales(t *testing.T) {

	executor := getExecutorTest()

	respuesta := executor.HorasMensuales()

	expected, err := strconv.ParseFloat(getHorasMensualesNormalesTest(), 64)

	if err != nil {
		t.Errorf("No se pudo convertir %s a float", getHorasMensualesNormalesTest())
	}

	if respuesta != expected{
		t.Errorf("La funcion HorasMensuales devuelve %f; want y se esperaba %f", respuesta, expected)
	}
}

func TestDiasMesTrabajadosFechaLiquidacion(t *testing.T) {

	executor := getExecutorTest()

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaTest())

	respuesta := executor.DiasMesTrabajadosFechaLiquidacion()

	expected := float64(0)

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionAntesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionDespuesDeAltaTest())

	respuesta = executor.DiasMesTrabajadosFechaLiquidacion()

	expected = float64(getFechaLiquidacionDespuesDeAltaTest().Day())

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionDespuesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaLiquidacion(&executor, getFechaLiquidacionAntesDeAltaMismoMesTest())

	respuesta = executor.DiasMesTrabajadosFechaLiquidacion()

	expected = float64(getFechaLiquidacionAntesDeAltaMismoMesTest().Day() - getFechaAltaTest().Day())

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaLiquidacion con getFechaLiquidacionAntesDeAltaMismoMesTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}
}


func TestDiasMesTrabajadosFechaPeriodo(t *testing.T) {

	executor := getExecutorTest()

	setFechaPeriodoLiquidacion(&executor, getFechaPeriodoLiquidacionAntesDeAltaTest())

	respuesta := executor.DiasMesTrabajadosFechaPeriodo()

	expected := float64(0)

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaPeriodo con getFechaPeriodoLiquidacionAntesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaPeriodoLiquidacion(&executor, getFechaPeriodoLiquidacionDespuesDeAltaTest())

	respuesta = executor.DiasMesTrabajadosFechaPeriodo()

	expected = float64(now.New(getFechaPeriodoLiquidacionDespuesDeAltaTest()).EndOfMonth().Day())

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaPeriodo con getFechaPeriodoLiquidacionDespuesDeAltaTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}

	setFechaPeriodoLiquidacion(&executor, getFechaPeriodoLiquidacionAntesDeAltaMismoMesTest())

	respuesta = executor.DiasMesTrabajadosFechaPeriodo()

	expected = float64(now.New(getFechaPeriodoLiquidacionAntesDeAltaMismoMesTest()).EndOfMonth().Day() - getFechaAltaTest().Day())

	if respuesta != expected{
		t.Errorf("La funcion DiasMesTrabajadosFechaPeriodo con getFechaPeriodoLiquidacionAntesDeAltaMismoMesTest devuelve %f; want y se esperaba %f", respuesta, expected)
	}
}

*/