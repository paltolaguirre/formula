package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/conexionBD/Formula/structFormula"
	"github.com/xubiosueldos/framework"
)

type IdsAEliminar struct {
	Ids []int `json:"ids"`
}

var nombreMicroservicio string = "formula"

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func FormulaList(w http.ResponseWriter, r *http.Request) {

	var tipo = r.URL.Query()["type"]

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		var formulas []structFormula.Formula

		if tipo != nil {
			db.Set("gorm:auto_preload", true).Where("type = ?", tipo).Find(&formulas)
		} else {
			db.Set("gorm:auto_preload", true).Find(&formulas)
		}

		framework.RespondJSON(w, http.StatusOK, formulas)
	}

}

func FormulaShow(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		formulaId := params["id"]

		var formulaData structFormula.FormulaWrapper
		var formulaPersistence structFormula.Formula //Con &var --> lo que devuelve el metodo se le asigna a la var
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&formulaPersistence, "id = ?", formulaId).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		formulaData.GormModel = formulaPersistence.GormModel
		json.Unmarshal([]byte(formulaPersistence.Value), &formulaData.FormulaPrime)

		framework.RespondJSON(w, http.StatusOK, formulaData)
	}

}

func FormulaAdd(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		var formulaData structFormula.FormulaWrapper
		var formulaPersistence structFormula.Formula
		var formulaPrime structFormula.FormulaPrime
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&formulaPrime); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		jsonValue, _ := json.Marshal(formulaPrime)
		formulaPersistence.Value = string(jsonValue)
		if err := db.Create(&formulaPersistence).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		formulaData.GormModel = formulaPersistence.GormModel
		json.Unmarshal([]byte(formulaPersistence.Value), &formulaData.FormulaPrime)

		framework.RespondJSON(w, http.StatusCreated, formulaData)
	}
}

func FormulaUpdate(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		//se convirtió el string en int para poder comparar
		paramFormulaid, _ := strconv.ParseInt(params["id"], 10, 64)
		p_formulaid := int(paramFormulaid)

		if p_formulaid == 0 {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroVacio)
			return
		}

		decoder := json.NewDecoder(r.Body)

		var formulaData structFormula.Formula

		if err := decoder.Decode(&formulaData); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		formulaid := formulaData.ID

		if p_formulaid == formulaid || formulaid == 0 {

			formulaData.ID = p_formulaid

			tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
			db := conexionBD.ObtenerDB(tenant)

			defer conexionBD.CerrarDB(db)

			if err := db.Save(&formulaData).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			framework.RespondJSON(w, http.StatusOK, formulaData)

		} else {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroDistintoStruct)
			return
		}
	}

}

func FormulaRemove(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		formulaId := params["id"]

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db)

		//--Borrado Fisico
		if err := db.Unscoped().Where("id = ?", formulaId).Delete(structFormula.Formula{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, framework.Formula+formulaId+framework.MicroservicioEliminado)
	}

}

func FormulaRemoveMasivo(w http.ResponseWriter, r *http.Request) {
	var resultadoDeEliminacion = make(map[int]string)
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		var idsEliminar IdsAEliminar
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&idsEliminar); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		if len(idsEliminar.Ids) > 0 {
			for i := 0; i < len(idsEliminar.Ids); i++ {
				formulaId := idsEliminar.Ids[i]
				if err := db.Unscoped().Where("id = ?", formulaId).Delete(structFormula.Formula{}).Error; err != nil {
					//framework.RespondError(w, http.StatusInternalServerError, err.Error())
					resultadoDeEliminacion[formulaId] = string(err.Error())

				} else {
					resultadoDeEliminacion[formulaId] = "Fue eliminado con exito"
				}
			}
		} else {
			framework.RespondError(w, http.StatusInternalServerError, "Seleccione por lo menos un registro")
		}

		framework.RespondJSON(w, http.StatusOK, resultadoDeEliminacion)
	}

}
