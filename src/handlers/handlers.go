//Модуль обрабтчиков http-запросов

package handlers

import (
	"database/sql"
	"cities/src/dbInterface"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"cities/src/structs"
)

func getId(r *http.Request) (int, error) {
	query := strings.Split(r.URL.Path, "/")
	Id, err := strconv.Atoi(query[len(query)-1])
	if err != nil {
		outErr := errors.New("ID must be of int type")
		return Id, outErr
	}
	return Id, nil
}

func outError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
}

func GetCityInfo(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(r)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		outBuff, err := dbInterface.DbGetCityInfo(id, db)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(outBuff)
	}
}

func AddCityInfo(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		newCity := new(structs.CityInfo)
		err = json.Unmarshal(body, newCity)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		err = dbInterface.DbAddCity(newCity, db)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("City %s added\n", newCity.Name)))

	}
}

func DeleteCity(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(r)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		err = dbInterface.DbDeleteCity(id, db)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Deleted city with ID %d", id)))
	}
}

func UpdatePopulation(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getId(r)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		population := new(structs.NewPopulation)
		err = json.Unmarshal(body, population)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		err = dbInterface.DbPopulationUpdate(id, population, db)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "New population set")
	}
}

func ListByDistrict(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		district := new(structs.StringQuery)
		err = json.Unmarshal(body, district)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		outBuff, err := dbInterface.DbListByDistrict(district, db)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(outBuff)
	}

}

func ListByRegion(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		region := new(structs.StringQuery)
		err = json.Unmarshal(body, region)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		outBuff, err := dbInterface.DbListByRegion(region, db)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(outBuff)
	}
}

func ListByPopulation(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		populationRange := new(structs.Values)
		err = json.Unmarshal(body, populationRange)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		if populationRange.MinValue > populationRange.MaxValue && populationRange.MaxValue != 0 {
			outError(w, http.StatusBadRequest, errors.New("chech the population range data"))
			return
		}
		outBuff, err := dbInterface.DbListByPopulation(populationRange, db)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(outBuff)
	}
}

func ListByFoundation(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		foundationRange := new(structs.Values)
		err = json.Unmarshal(body, foundationRange)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		if foundationRange.MinValue > foundationRange.MaxValue && foundationRange.MaxValue != 0 {
			outError(w, http.StatusBadRequest, errors.New("chech the foundation range data"))
			return
		}
		outBuff, err := dbInterface.DbListByFoundation(foundationRange, db)
		if err != nil {
			outError(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(outBuff)
	}

}
