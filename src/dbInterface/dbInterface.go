//Модуль запросов к базе ProgreSQL

package dbInterface

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"cities/src/structs"
)

func DbEmptyCheck(db *sql.DB) (empty bool, dbErr error) {
	request := "SELECT * FROM citydata limit 5"
	resp, dbErr := db.Query(request)
	if dbErr != nil {
		return empty, dbErr
	}
	empty = !resp.Next()
	return empty, nil
}

func DbFillFromCsv(fileName string, db *sql.DB) error {
	inFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer inFile.Close()
	inFileInfo, _ := inFile.Stat()
	var buff = make([]byte, inFileInfo.Size())
	_, err = inFile.Read(buff)
	if err != nil {
		return err
	}
	r := csv.NewReader(strings.NewReader(string(buff)))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		newCityData := new(structs.CityInfo)
		newCityData.Id, err = strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		newCityData.Name = record[1]
		newCityData.Region = record[2]
		newCityData.District = record[3]
		newCityData.Population, err = strconv.Atoi(record[4])
		if err != nil {
			return err
		}
		newCityData.Foundation, err = strconv.Atoi(record[5])
		if err != nil {
			return err
		}
		err = DbAddCity(newCityData, db)
		if err != nil {
			return err
		}
	}
	return nil
}

func DbNextIndex(db *sql.DB) (int, error) { //Вспомогательная функция поиска свободного индекса. Не задействовано.
	request := "select max(cityid) from citydata"
	resp, dbErr := db.Query(request)
	if dbErr != nil {
		return 0, dbErr
	}
	defer resp.Close()
	var maxId string
	if !resp.Next() {
		return 0, errors.New("max index search error")
	}
	err := resp.Scan(&maxId)
	if err != nil {
		return 0, err
	}
	nextId, _ := strconv.Atoi(maxId)
	return nextId + 1, nil
}

func DbGetCityInfo(id int, db *sql.DB) ([]byte, error) {
	request := fmt.Sprint("SELECT cityName, region, district, population, foundation FROM cityData WHERE cityID=", id)
	resp, dbErr := db.Query(request)
	if dbErr != nil {
		return nil, dbErr
	}
	defer resp.Close()
	if !resp.Next() {
		return nil, errors.New("no city with such ID was found")
	}
	cityData := new([5]string)
	err := resp.Scan(&cityData[0], &cityData[1], &cityData[2], &cityData[3], &cityData[4])
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %s %s %s", cityData[0], cityData[1], cityData[2], cityData[3], cityData[4])), nil
}

func DbAddCity(newCity *structs.CityInfo, db *sql.DB) error {
	request := fmt.Sprintf("INSERT INTO citydata VALUES ('%d', '%s', '%s', '%s', '%d', '%d')",
		newCity.Id, newCity.Name, newCity.Region, newCity.District, newCity.Population, newCity.Foundation)
	_, dbErr := db.Exec(request)
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func DbDeleteCity(id int, db *sql.DB) error {
	request := fmt.Sprintf("DELETE FROM citydata WHERE cityid=%d", id)
	_, dbErr := db.Exec(request)
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func DbPopulationUpdate(id int, population *structs.NewPopulation, db *sql.DB) error {
	request := fmt.Sprintf("update citydata set population = %d where cityid = %d", population.Value, id)
	_, dbErr := db.Exec(request)
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func DbListByRegion(region *structs.StringQuery, db *sql.DB) ([]byte, error) {
	request := fmt.Sprintf("SELECT * FROM cityData WHERE region='%s'", region.Request)
	resp, dbErr := db.Query(request)
	if dbErr != nil {
		return nil, dbErr
	}

	defer resp.Close()
	cityList := ""
	for resp.Next() {
		cityData := new([6]string)
		err := resp.Scan(&cityData[0], &cityData[1], &cityData[2], &cityData[3], &cityData[4], &cityData[5])
		if err != nil {
			return nil, err
		}
		cityList += fmt.Sprintf("%s %s %s %s %s %s\n", cityData[0], cityData[1], cityData[2], cityData[3], cityData[4], cityData[5])
	}
	if cityList == "" {
		return []byte("No cities were found in region " + region.Request), nil
	}
	return []byte(cityList), nil
}

func DbListByDistrict(district *structs.StringQuery, db *sql.DB) ([]byte, error) {
	request := fmt.Sprintf("SELECT * FROM cityData WHERE district='%s'", district.Request)
	resp, dbErr := db.Query(request)
	if dbErr != nil {
		return nil, dbErr
	}

	defer resp.Close()
	cityList := ""
	for resp.Next() {
		cityData := new([6]string)
		err := resp.Scan(&cityData[0], &cityData[1], &cityData[2], &cityData[3], &cityData[4], &cityData[5])
		if err != nil {
			return nil, err
		}
		cityList += fmt.Sprintf("%s %s %s %s %s %s\n", cityData[0], cityData[1], cityData[2], cityData[3], cityData[4], cityData[5])
	}
	if cityList == "" {
		return []byte("No cities were found in district " + district.Request), nil
	}
	return []byte(cityList), nil
}

func DbListByPopulation(populationRange *structs.Values, db *sql.DB) ([]byte, error) {
	request := ""
	if populationRange.MaxValue != 0 {
		request = fmt.Sprintf("SELECT * FROM citydata WHERE population>=%d AND population<=%d", populationRange.MinValue, populationRange.MaxValue)
	} else {
		request = fmt.Sprintf("SELECT * FROM citydata WHERE population>=%d", populationRange.MinValue)
	}
	resp, dbErr := db.Query(request)
	if dbErr != nil {
		return nil, dbErr
	}
	defer resp.Close()
	cityList := ""
	for resp.Next() {
		cityData := new([6]string)
		err := resp.Scan(&cityData[0], &cityData[1], &cityData[2], &cityData[3], &cityData[4], &cityData[5])
		if err != nil {
			return nil, err
		}
		cityList += fmt.Sprintf("%s %s %s %s %s %s\n", cityData[0], cityData[1], cityData[2], cityData[3], cityData[4], cityData[5])
	}
	if cityList == "" {
		if populationRange.MaxValue != 0 {
			return []byte(fmt.Sprintf("No cities were found with population range from %d to %d", populationRange.MinValue, populationRange.MaxValue)), nil
		} else {
			return []byte(fmt.Sprintf("No cities were found with population starting from %d", populationRange.MinValue)), nil
		}
	}
	return []byte(cityList), nil
}

func DbListByFoundation(foundationRange *structs.Values, db *sql.DB) ([]byte, error) {
	request := ""
	if foundationRange.MaxValue != 0 {
		request = fmt.Sprintf("SELECT * FROM citydata WHERE foundation>=%d AND foundation<=%d", foundationRange.MinValue, foundationRange.MaxValue)
	} else {
		request = fmt.Sprintf("SELECT * FROM citydata WHERE foundation>=%d", foundationRange.MinValue)
	}
	resp, dbErr := db.Query(request)
	if dbErr != nil {
		return nil, dbErr
	}
	defer resp.Close()
	cityList := ""
	for resp.Next() {
		cityData := new([6]string)
		err := resp.Scan(&cityData[0], &cityData[1], &cityData[2], &cityData[3], &cityData[4], &cityData[5])
		if err != nil {
			return nil, err
		}
		cityList += fmt.Sprintf("%s %s %s %s %s %s\n", cityData[0], cityData[1], cityData[2], cityData[3], cityData[4], cityData[5])
	}
	if cityList == "" {
		if foundationRange.MaxValue != 0 {
			return []byte(fmt.Sprintf("No cities were found with foundation range from %d to %d", foundationRange.MinValue, foundationRange.MaxValue)), nil
		} else {
			return []byte(fmt.Sprintf("No cities were found with foundation starting from %d", foundationRange.MinValue)), nil
		}
	}
	return []byte(cityList), nil
}

func DbBackup(db *sql.DB) error {
	request := "SELECT * FROM cityData"
	resp, dbErr := db.Query(request)
	if dbErr != nil {
		return dbErr
	}
	defer resp.Close()
	outList := ""
	for resp.Next() {
		cityData := new([6]string)
		err := resp.Scan(&cityData[0], &cityData[1], &cityData[2], &cityData[3], &cityData[4], &cityData[5])
		if err != nil {
			return err
		}
		outList += fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", cityData[0], cityData[1], cityData[2], cityData[3], cityData[4], cityData[5])
	}
	err := os.WriteFile("cities.csv", []byte(outList), 0666)
	return err
}
