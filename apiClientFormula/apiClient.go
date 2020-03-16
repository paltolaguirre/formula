package apiClientFormula

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"git-codecommit.us-east-1.amazonaws.com/v1/repos/sueldos-formula/executor"
	"github.com/xubiosueldos/conexionBD/Function/structFunction"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/framework/configuracion"
	"io/ioutil"
	"net/http"
)

func ExecuteFormulaLiquidacion(authorization string, liquidacion *structLiquidacion.Liquidacion, formulaName string) (float64, error) {

	config := configuracion.GetInstance()
	url := configuracion.GetUrlMicroservicio(config.Puertomicroservicioformula) + "api/formula/execute"

	executeBody := executor.FormulaExecute{
		Context: executor.Context{Currentliquidacion: *liquidacion},
		Invoke: structFunction.Invoke{
			Functionname: formulaName,
		},
	}

	requestByte, err := json.Marshal(executeBody)

	if err != nil {
		return 0, errors.New("Error al convertir el body a string: " + err.Error())
	}

	requestReader := bytes.NewReader(requestByte)

	req, err := http.NewRequest("POST", url, requestReader)

	if err != nil {
		return 0, errors.New("Error al generar el request a " + url + ": " + err.Error())
	}

	req.Header.Add("Authorization", authorization)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, errors.New("Error al enviar el request a " + url + ": " + err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return 0, errors.New("Error al leer la respuesta de " + url +": " + err.Error())
	}

	if res.StatusCode != http.StatusCreated {
		return 0, errors.New("No se pudo resolver la formula")
	}

	var value structFunction.Value

	json.Unmarshal([]byte(string(body)), value)

	return value.Valuenumber, nil

}
