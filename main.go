package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

// UART0
// 9600 baud
// PA4 -> TX
// PA5 -> RX

func main() {
	// Эта горутина будет работать в фоновом режиме, пока основная функция
	// ожидает завершения, что полезно для серверных приложений.
	go func() {
		for {
			dataFromUART, err := GetDataSensors()
			if err != nil {
				// Используем log.Printf для автоматического добавления времени
				log.Printf("Критическая ошибка получения данных с UART: %v\n", err)
			} else if dataFromUART != "" {
				// Выводим данные, только если они не пустые
				fmt.Printf("Данные с UART: %s\n", dataFromUART)
			}
			// Небольшая задержка, чтобы не загружать процессор в цикле
			time.Sleep(time.Second)
		}
	}()

	// Функция main должна работать бесконечно, иначе программа завершится.
	// Можно использовать канал или select, но пока просто заставим ее ждать.
	fmt.Println("Программа запущена. Ожидание данных с UART...")
	// Замените это на более надежный механизм, если это часть большего приложения.
	select {} // Бесконечное ожидание
}

// GetDataSensors считывает данные с последовательного порта Repka Pi 3 и возвращает их.
func GetDataSensors() (string, error) {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	// Рекомендуется использовать алиас "/dev/serial0", который указывает на нужный GPIO UART
	// (ttyS0 на Repka Pi 3 с включенным enable_uart=1).
	portName := "/dev/serial0"
	port, err := serial.Open(portName, mode)
	if err != nil {
		// Используйте fmt.Errorf для оборачивания ошибки, чтобы caller мог ее обработать
		return "", fmt.Errorf("не удалось открыть порт %s: %w", portName, err)
	}
	defer port.Close() // Закрываем порт при выходе из функции

	// Устанавливаем таймаут чтения, чтобы функция не блокировалась навсегда
	port.SetReadTimeout(time.Second * 2)

	buffer := make([]byte, 200)
	n, err := port.Read(buffer)

	// go.bug.st/serial часто возвращает EOF при таймауте чтения
	if err != nil && err.Error() != "EOF" && !strings.Contains(err.Error(), "timeout") {
		// Оборачиваем ошибку чтения
		return "", fmt.Errorf("ошибка чтения из порта: %w", err)
	}

	if n == 0 {
		return "", nil // Ничего не прочитали
	}

	// Преобразуем прочитанные байты в строку и обрезаем лишние нулевые байты
	data := string(buffer[:n])
	// Опционально: если данные содержат мусор после полезной нагрузки, можно почистить
	data = strings.TrimSpace(data)

	// Для отладки можно вывести необработанные данные
	// fmt.Printf("Прочитано %d байт: '%s'\n", n, data)

	return data, nil
}
