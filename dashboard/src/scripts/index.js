import '../styles/style.css'
import { 
    Chart, 
    LineController, 
    LineElement, 
    PointElement, 
    LinearScale, 
    CategoryScale, 
    Tooltip, 
    Legend 
} from 'chart.js';

// 1. Регистрируем компоненты Chart.js
Chart.register(
    LineController, 
    LineElement, 
    PointElement, 
    LinearScale, 
    CategoryScale, 
    Tooltip, 
    Legend
);

/**
 * Универсальная функция загрузки данных
 */
async function fetchData(url) {
    try {
        const response = await fetch(url);
        if (response.ok) {
            const json = await response.json();
            console.log("Данные успешно загружены:", json);
            return json; 
        } else {
            console.error(`Ошибка сервера: ${response.status}`);
        }
    } catch (error) {
        console.error("Ошибка сети:", error);
    }
    return null;
}

/**
 * Отрисовка графика
 */
async function renderChart(canvasId, stationID, hours, color, dataKey, label) {
    const canvasElement = document.getElementById(canvasId);
    if (!canvasElement) return;

    // Формируем URL. Например: /0/data?hours=5
    const apiUrl = `/${stationID}/data?hours=${hours}`; 
    const apiDataArray = await fetchData(apiUrl);

    if (!apiDataArray || !Array.isArray(apiDataArray) || apiDataArray.length === 0) {
        console.warn(`Нет данных для отрисовки ${label} (ID: ${stationID})`);
        return;
    }

    // Подготовка меток времени (Ось X)
    const labels = apiDataArray.map(item => {
        const date = new Date(item.timestamp);
        return date.toLocaleTimeString('ru-RU', { 
            hour: '2-digit', 
            minute: '2-digit'
        });
    });

    // Подготовка значений (Ось Y)
    const values = apiDataArray.map(item => item[dataKey]);
   
    const ctx = canvasElement.getContext('2d');
    
    return new Chart(ctx, {
        type: 'line', 
        data: {
            labels: labels,
            datasets: [{
                label: `${label} (Станция ${stationID})`,
                data: values,
                borderColor: color,
                backgroundColor: color.replace('1)', '0.2)'), // легкая заливка
                fill: true,
                tension: 0.3,
                pointRadius: 2,
                borderWidth: 2
            }]
        },
        options: {
            responsive: true, 
            maintainAspectRatio: false, 
            scales: {
                y: {
                    beginAtZero: false, // Чтобы лучше видеть колебания
                    ticks: {
                        color: '#666'
                    }
                },
                x: {
                    ticks: {
                        maxRotation: 45,
                        minRotation: 45
                    }
                }
            },
            plugins: {
                legend: {
                    position: 'top',
                }
            }
        }
    });
}

document.addEventListener('DOMContentLoaded', async () => {
    const currentPath = window.location.pathname; 
    const urlParams = new URLSearchParams(window.location.search);
    const hours = urlParams.get('hours') || '1'; 

    // Задаем ID вашей станции
    const STATION_ID = 0;

    if (currentPath === '/' || currentPath === '/index.html') {
        // Рендерим 3 разных графика для станции 0 на главной странице
        await Promise.all([
            renderChart('myChart1', STATION_ID, hours, 'rgba(255, 99, 132, 1)', 'temperature', 'Температура (°C)'),
            renderChart('myChart2', STATION_ID, hours, 'rgba(54, 162, 235, 1)', 'humidity', 'Влажность (%)'),
            renderChart('myChart3', STATION_ID, hours, 'rgba(75, 192, 192, 1)', 'co2', 'CO2 (ppm)')
        ]);

    } else {
        // Если путь динамический (например, /0), извлекаем ID из пути
        const pathSegments = currentPath.split('/').filter(Boolean);
        const dynamicId = pathSegments.length > 0 ? pathSegments[0] : STATION_ID; 

        if (document.getElementById('myChart1')) {
            await renderChart('myChart1', dynamicId, hours, 'rgba(0, 119, 255, 1)', 'temperature', 'Температура');
        }
    }
});
