package web_server

import (
	"WeatherStation/resource"
	"WeatherStation/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Инициализация конфига для переменных окружения
// и базы данных
var Config utils.Config
var db *sql.DB

// Вставка данных в базу данных
func insertMeteoData(db *sql.DB, stationID int, data resource.SensorData) error {

	insertSQL := `
        INSERT INTO meteodata (
            meteostation_id, 
            humidity, 
            temperature, 
            pressure, 
            CO2, 
            TVOC
        ) VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id`

	var insertedID int

	// db.QueryRow выполняет запрос и ожидает одну строку результата.
	err := db.QueryRow(
		insertSQL,
		stationID,
		data.Hum,
		data.Temp,
		data.Pres,
		data.Co2,
		data.Tvoc,
	).Scan(&insertedID)

	if err != nil {
		return fmt.Errorf("ошибка вставки метеоданных в БД: %w", err)
	}

	log.Printf("Данные успешно вставлены в таблицу meteodata с ID: %d", insertedID)

	return nil
}

func generateDBUrl(config utils.Config) string {
	dbURL := fmt.Sprintf("user=%s "+
		"password=%s "+
		"dbname=%s "+
		"host=%s "+
		"sslmode=%s "+
		"port=%d",
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBHost,
		config.DBSSLMode,
		config.DBPort,
	)

	return dbURL
}

// Подключение к базе данных
func DBConnection(config utils.Config) {

	var err error
	db, err = sql.Open(config.DBName, generateDBUrl(config))
	if err != nil {
		log.Print("Ошибка при открытии базы данных:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Print("Ошибка при пинге базы данных:", err)
	}
}

// Обновляет базу данных каждые 10 секунд
func updateDB() {
	fmt.Println("Web_server: Запущена фоновая горутина обновления БД (каждые 10 секунд)...")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		func() {
			log.Println("Выполняю HTTP GET запрос к коллектору...")
			collectorURL := fmt.Sprintf("http://%s:%s/data", Config.CollectorHost, Config.WebPort)

			resp, err := http.Get(collectorURL)
			if err != nil {
				log.Printf("Ошибка при запросе к Collector: %s", err)
				return
			}

			defer resp.Body.Close()

			var receivedData resource.SensorData
			err = json.NewDecoder(resp.Body).Decode(&receivedData)
			if err != nil {
				log.Printf("Ошибка декодирования JSON: %v", err)
				return
			}

			err = insertMeteoData(db, 0, receivedData)
			if err != nil {
				log.Printf("Ошибка при вставке данных в БД: %v", err)
				return
			}

			log.Println("Запрос к коллектору успешно выполнен и обработан.")
		}()
	}
}

// Извлекаем данные из базы данных для клиента
func fetchMeteoDataFromDB(db *sql.DB, stationID int, hours int) ([]resource.MeteoDataPoint, error) {
	timeAgo := time.Now().Add(time.Duration(-hours) * time.Hour)

	query := `
        SELECT 
            -- Используем новое поле created_at
            created_at, 
            temperature,
            humidity,
            pressure,
            co2,
            tvoc
        FROM meteodata
        WHERE meteostation_id = $1 AND created_at >= $2
        ORDER BY created_at ASC;
    `

	rows, err := db.Query(query, stationID, timeAgo)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения SQL-запроса: %w", err)
	}
	defer rows.Close()

	var dataPoints []resource.MeteoDataPoint
	for rows.Next() {
		var dp resource.MeteoDataPoint
		if err := rows.Scan(
			&dp.Timestamp,
			&dp.Temperature,
			&dp.Humidity,
			&dp.Pressure,
			&dp.Co2,
			&dp.Tvoc,
		); err != nil {
			fmt.Printf("Ошибка сканирования строки БД: %v\n", err)
			continue
		}
		dataPoints = append(dataPoints, dp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по результатам БД: %w", err)
	}

	return dataPoints, nil
}

// Запрос на извлечение данных с базы данных
func getStationDataHandler(w http.ResponseWriter, r *http.Request) {

	stationIDStr := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/data"), "/")

	if stationIDStr == "" {
		http.Error(w, "Отсутствует core_id", http.StatusBadRequest)
		return
	}

	stationID, err := strconv.Atoi(stationIDStr)
	if err != nil {
		http.Error(w, "Неверный core_id", http.StatusBadRequest)
		return
	}

	hours, err := strconv.Atoi(r.URL.Query().Get("hours"))
	if err != nil || hours <= 0 {
		http.Error(w, "Неверный hours", http.StatusBadRequest)
		return
	}

	dataPoints, err := fetchMeteoDataFromDB(db, stationID, hours)
	if err != nil {
		log.Printf("Ошибка получения данных из БД: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dataPoints); err != nil {
		log.Printf("Ошибка при кодировании JSON: %v", err)
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./dist/index.html")
}

func StartWebServer(port string) {
	Config.LoadConfig()
	go updateDB()

	mux := http.NewServeMux()

	mux.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("./dist"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/data") {
			getStationDataHandler(w, r)
			return
		}
		pageHandler(w, r)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
