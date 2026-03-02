
CREATE TABLE IF NOT EXISTS meteostation (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_out TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_out TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_out TIMESTAMP WITH TIME ZONE -- Это поле может быть NULL, если запись не удалена
);

CREATE TABLE IF NOT EXISTS meteodata (
    id SERIAL PRIMARY KEY,
    meteostation_id INTEGER NOT NULL,

    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    humidity    REAL    NOT NULL, 
    temperature REAL    NOT NULL, 
    pressure    REAL    NOT NULL, 
    CO2         SMALLINT NOT NULL, 
    TVOC        SMALLINT NOT NULL, 
    

    FOREIGN KEY (meteostation_id) REFERENCES meteostation (id) ON DELETE CASCADE
);