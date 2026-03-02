package handlers

import (
	"WeatherStation/collector/cpufreq"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Определяем структуру для хранения информации с датчиков
type Handlers struct {
	Sensors *cpufreq.SensorsInfo
}

// Определяем хэндлер для обработки запроса по
// порту с инпоинтом data
func (h *Handlers) Init() {
	http.HandleFunc("/data", h.DataHandler)
}

// Собираем информацию с датчиков
// и отправляем Веб-серверу
func (h *Handlers) DataHandler(w http.ResponseWriter, r *http.Request) {
	h.Sensors.DynamicDataMutex.RLock()
	dataToSend := h.Sensors.DynamicData
	h.Sensors.DynamicDataMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(dataToSend)

	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка кодирования JSON: %v", err), http.StatusInternalServerError)
		log.Printf("Ошибка записи ответа: %v", err)
	}
}
