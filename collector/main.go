package main

import (
	"WeatherStation/collector/cpufreq"
	"WeatherStation/collector/handlers"
	"WeatherStation/utils"
	"fmt"
	"log"
	"net/http"
)

// Инициализация конфигурации
var Config utils.Config

func main() {

	// Загружаем конфигурацию
	Config.LoadConfig()

	// Определяем хранение данных с датчиков
	dataSensors := cpufreq.SensorsInfo{}

	// Собираем данные с датчика
	dataSensors.CollectSensorData()

	// Передаём в структуру для дальнейшей инициализации ответа
	// для Веб-сервера
	h := handlers.Handlers{Sensors: &dataSensors}
	h.Init()

	// Слушаем внутренний порт общения между
	// коллектором и веб-сервером
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Config.WebPort), nil))
}
