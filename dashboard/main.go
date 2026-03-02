package main

import (
	"WeatherStation/dashboard/web_server"
	"WeatherStation/utils"
	"database/sql"

	_ "github.com/lib/pq"
)

// Инициализация конфига для переменных окружения
// и базы данных
var Config utils.Config
var db *sql.DB

func main() {

	// Загрузка конфига из .env
	Config.LoadConfig()

	// Подключение базы данных
	web_server.DBConnection(Config)

	// Откладывание закрытия базы до выхода из main
	defer db.Close()

	// Запуск Веб-сервера
	web_server.StartWebServer(Config.WebPort)

}
