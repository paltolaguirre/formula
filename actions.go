package main

import (
	"encoding/json"
	"net/http"

	"git-codecommit.us-east-1.amazonaws.com/v1/repos/sueldos-formula/executor"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/conexionBD/Function/structFunction"
	"github.com/xubiosueldos/framework"
)

type IdsAEliminar struct {
	Ids []int `json:"ids"`
}

type FormulaExecute struct {
	structFunction.Invoke
	Context executor.Context `json:"context"`
}

var nombreMicroservicio string = "formula"

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func FunctionList(w http.ResponseWriter, r *http.Request) {

	var tipo = r.URL.Query()["type"]

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		var functions []structFunction.Function

		if tipo != nil {
			db.Set("gorm:auto_preload", true).Where("type = ?", tipo).Find(&functions)
		} else {
			db.Set("gorm:auto_preload", true).Find(&functions)
		}

		framework.RespondJSON(w, http.StatusOK, functions)
	}

}

func FunctionShow(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {
		params := mux.Vars(r)
		functionName := params["id"]

		var function structFunction.Function //Con &var --> lo que devuelve el metodo se le asigna a la var

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db)

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&function, "name = ?", functionName).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, function)
	}

}

func FunctionAdd(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var functionData structFunction.Function

		if err := decoder.Decode(&functionData); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		/*res2B, _ := json.Marshal(functionData)
		fmt.Fprintf(w, string(res2B))
		return*/

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)
		//abro una transacción para que si hay un error no persista en la DB
		tx := db.Begin()
		defer tx.Rollback()
		
		/*err := createValue(functionData.Value, tx)
		if err != nil {
			tx.Rollback()
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}*/

		if err := tx.Create(&functionData).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		tx.Commit()
		framework.RespondJSON(w, http.StatusCreated, functionData)
	}

}

func FunctionUpdate(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		//se convirtió el string en int para poder comparar
		pFunctionName := params["id"]

		if pFunctionName == "" {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroVacio)
			return
		}

		decoder := json.NewDecoder(r.Body)

		var formulaData structFunction.Function

		if err := decoder.Decode(&formulaData); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		//abro una transacción para que si hay un error no persista en la DB
		tx := db.Begin()

		var formula structFunction.Function
		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := tx.Set("gorm:auto_preload", true).First(&formula, "name = ?", formulaData.Name).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		err := createValue(formulaData.Value, tx)
		if err != nil {
			tx.Rollback()
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

			//abro una transacción para que si hay un error no persista en la DB
			tx := db.Begin()
			defer tx.Rollback()

			var formula structFunction.Function
			//gorm:auto_preload se usa para que complete todos los struct con su informacion
			if err := tx.Set("gorm:auto_preload", true).First(&formula, "name = ?", formulaData.Name).Error; gorm.IsRecordNotFoundError(err) {
				framework.RespondError(w, http.StatusNotFound, err.Error())
				return
			}

			err := createValue(formulaData.Value, tx)
			if err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			if err := tx.Save(&formulaData).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			if err := deleteValue(formula.Value, tx); err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			tx.Commit()
			framework.RespondJSON(w, http.StatusOK, formulaData)

		} else {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroDistintoStruct)
			return
		}

		tx.Commit()
		framework.RespondJSON(w, http.StatusOK, formulaData)
	}

}

func FunctionRemove(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		functionName := params["id"]

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db)

		//abro una transacción para que si hay un error no persista en la DB
		tx := db.Begin()
		defer tx.Rollback()

		var function structFunction.Function
		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := tx.Set("gorm:auto_preload", true).First(&function, "name = ?", functionName).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		//--Borrado Fisico
		if err := tx.Unscoped().Where("name = ?", functionName).Delete(structFunction.Function{}).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err := deleteValue(function.Value, tx)
		if err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		tx.Commit()
		framework.RespondJSON(w, http.StatusOK, framework.Function+functionName+framework.MicroservicioEliminado)
	}

}

func FunctionRemoveMasivo(w http.ResponseWriter, r *http.Request) {
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
				functionName := idsEliminar.Ids[i]
				if err := db.Unscoped().Where("name = ?", functionName).Delete(structFunction.Function{}).Error; err != nil {
					//framework.RespondError(w, http.StatusInternalServerError, err.Error())
					resultadoDeEliminacion[i] = string(err.Error())
				} else {
					resultadoDeEliminacion[i] = "Fue eliminado con exito"
				}
			}
		} else {
			framework.RespondError(w, http.StatusInternalServerError, "Seleccione por lo menos un registro")
		}

		framework.RespondJSON(w, http.StatusOK, resultadoDeEliminacion)
	}

}

func createValue(value *structFunction.Value, db *gorm.DB) error {
	if value.Valueinvoke != nil {
		value.Valueinvoke.ID = 0
		for i := 0; i < len(value.Valueinvoke.Args); i++ {
			err := createValue(&value.Valueinvoke.Args[i], db)
			if err != nil {
				return err
			}
		}
	}

	// para asegurarse que siempre se cree un nuevo Value
	value.ID = 0
	if err := db.Create(&value).Error; err != nil {
		return err
	}

	return nil
}

func deleteValue(value *structFunction.Value, db *gorm.DB) error {
	if value.Valueinvoke != nil {
		for i := 0; i < len(value.Valueinvoke.Args); i++ {
			err := deleteValue(&value.Valueinvoke.Args[i], db)
			if err != nil {
				return err
			}
		}

		if err := db.Unscoped().Where("id = ?", value.Valueinvokeid).Delete(structFunction.Invoke{}).Error; err != nil {
			return err
		}
	}

	if err := db.Unscoped().Where("id = ?", value.ID).Delete(structFunction.Value{}).Error; err != nil {
		return err
	}

	return nil
}

func FunctionExecute(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var formulaExecuteData FormulaExecute

		if err := decoder.Decode(&formulaExecuteData); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		myExecutor := executor.NewExecutor(db, &formulaExecuteData.Context)
		value, err := myExecutor.GetValueFromInvoke(&formulaExecuteData.Invoke)

		if err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusCreated, value)
	}

}
