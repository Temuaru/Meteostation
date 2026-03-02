package cpufreq

import (
	"WeatherStation/resource"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"go.bug.st/serial"
)

// Определяем структуру динамического хранения
// информации с датчиков
type SensorsInfo struct {
	DynamicDataMutex sync.RWMutex
	DynamicData      resource.SensorData
}

func (s *SensorsInfo) CollectSensorData() { // <- Переименовано
	go func() {
		for {
			data, err := GetSensorsData()
			if err != nil {
				log.Printf("Ошибка получения данных: %v", err)
			} else {
				s.DynamicDataMutex.Lock()

				s.DynamicData = data

				// fmt.Printf("Данные успешно прочитаны и записаны: %+v\n", data)
				s.DynamicDataMutex.Unlock()
			}
			time.Sleep(time.Second)
		}
	}()
}

func GetSensorsData() (resource.SensorData, error) {
	// 1. Настройка порта
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	portName := "/dev/serial0"
	port, err := serial.Open(portName, mode)
	if err != nil {
		return resource.SensorData{}, fmt.Errorf("не удалось открыть порт %s: %w", portName, err)
	}
	defer port.Close()

	scanner := bufio.NewScanner(port)

	if scanner.Scan() {
		rawData := scanner.Bytes()
		fmt.Printf("Прочитана полная строка (%d байт): %s\n", len(rawData), string(rawData))

		var sensorData resource.SensorData
		err = json.Unmarshal(rawData, &sensorData)
		if err != nil {
			return resource.SensorData{}, fmt.Errorf("Ошибка парсинга JSON: %w. Получено: %s", err, string(rawData))
		}
		return sensorData, nil
	}

	if err := scanner.Err(); err != nil {
		return resource.SensorData{}, fmt.Errorf("ошибка при сканировании порта: %w", err)

	}

	return resource.SensorData{}, fmt.Errorf("данные не получены (пустой поток)")
}
