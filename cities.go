//Основная программа обработки запросов к базе данных городов

//Работает с базой данных PostgreSQL, скрипт создания рабочей таблицы в файле postgresql.sql
//При необходимости поправить имя пользователя, пароль и адрес в ini-файле settings.ini (формат json)
//Для первоначального заполнения базы можно использовать перечень городов, размещаемый в файле cities.csv
//Для остановки программы используйте ctrl+c
//при выходе из программы все данные базы сохраняются в файл cities.csv
//
//Требуемые для работы модули:
//1. github.com/lib/pq
//2. Модуль обработчиков запросов handlers
//3. Модуль с описанием структур (для передачи данных в формате json) structs
//4. Модуль для взаимодействия с базой данных dbInterface
//
//Примеры запросов:
//получение информации о городе по его id: GET-запрос по адреу вида http://server_adress:server_port/cities/xxx, где xxx- уникальный ID города
//добавление новой записи в список городов: POST-запрос по адресу вида http://server_adress:server_port/cities с json-структурой вида:
//	{	"id": int,
//		"name": string,
//		"region": string,
//		"district": string,
//		"population": int,
//		"foundation": int }
//удаление информации о городе по указанному id: запрос DELETE по адресу  http://server_adress:server_port/cities/xxx,
// где xxx- уникальный ID города
//обновление информации о численности населения города по указанному id: запрос PUT по адресу http://server_adress:server_port/cities/xxx
// где xxx- уникальный ID города, с json-структурой вида: {"value": int}
//получение списка городов по указанному региону: запрос POST http://server_adress:server_port/cities/region
// с json-структурой вида: {"request": string}
//получение списка городов по указанному округу: запрос POST http://server_adress:server_port/cities/district
// с json-структурой вида: {"request": string}
//получения списка городов по указанному диапазону численности населения: запрос POST http://server_adress:server_port/cities/population
// с json-структурой вида: {"min_value": int, "max_value": int}. Допускается указание одной границы диапазона.
//получения списка городов по указанному диапазону года основания: запрос POST http://server_adress:server_port/cities/foundation
// с json-структурой вида: {"min_value": int, "max_value": int}. Допускается указание одной границы диапазона.

package main

import (
	"context"
	"database/sql"
	"cities/src/dbInterface"
	"encoding/json"
	"fmt"
	"cities/src/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

type parameters struct {
	ServerAdress string
	ServerPort   string
	DbAdress     string
	DbName       string
	DbPort       string
	DbUserName   string
	DbPassword   string
}

func initRead(initParams *parameters) error {
	iniFile, err := os.Open("settings.ini")
	if err != nil {
		return err
	}
	defer iniFile.Close()
	fileInfo, _ := iniFile.Stat()
	buff := make([]byte, fileInfo.Size())
	_, err = iniFile.Read(buff)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, initParams)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	var initParams parameters
	err := initRead(&initParams)
	if err != nil {
		log.Fatal(err.Error())
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	connectAttributes := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		initParams.DbAdress, initParams.DbPort, initParams.DbUserName, initParams.DbPassword, initParams.DbName)
		
	db, err := sql.Open("postgres", connectAttributes)
	defer db.Close()
	if err != nil {
		log.Fatal("can't connect to db", err)
	}

	empty, err := dbInterface.DbEmptyCheck(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	if empty {
		log.Println("Db probably is empty, reading data from cities.csv")
		err = dbInterface.DbFillFromCsv("cities.csv", db)
		if err != nil {
			log.Println(err.Error())
			log.Println("Starting with empty database")
		}
	}
	fmt.Printf("db connected, waiting for command on adress %s port %s\n", initParams.ServerAdress, initParams.ServerPort)

	r := chi.NewRouter()
	if r != nil {
		r.Use(middleware.Logger)
	}
	r.Route("/cities", func(r chi.Router) {
		r.Get("/{city_Id}", handlers.GetCityInfo(db))
		r.Post("/", handlers.AddCityInfo(db))
		r.Delete("/{city_Id}", handlers.DeleteCity(db))
		r.Put("/{city_Id}", handlers.UpdatePopulation(db))
	})

	r.Route("/info", func(r chi.Router) {
		r.Post("/region", handlers.ListByRegion(db))
		r.Post("/district", handlers.ListByDistrict(db))
		r.Post("/population", handlers.ListByPopulation(db))
		r.Post("/foundation", handlers.ListByFoundation(db))
	})

	srv := &http.Server{
		Addr:    initParams.ServerAdress + ":" + initParams.ServerPort,
		Handler: r,
	}
	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Print(err.Error())
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed:%+v", err)
	}

	err = dbInterface.DbBackup(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Print("Node stopped")
}
