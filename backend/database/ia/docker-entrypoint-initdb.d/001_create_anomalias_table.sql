-- Crear tabla de anomalías detectadas por IA
CREATE TABLE IF NOT EXISTS anomalias (
    id SERIAL PRIMARY KEY,
    activo_id INTEGER NOT NULL,
    sensor_id VARCHAR(100),
    timestamp TIMESTAMP NOT NULL,
    anomaly_score FLOAT NOT NULL,
    anomaly_likelihood FLOAT NOT NULL,
    severidad VARCHAR(20) NOT NULL CHECK (severidad IN ('baja', 'media', 'alta', 'critica')),
    descripcion TEXT,
    threshold FLOAT NOT NULL,
    is_anomaly BOOLEAN NOT NULL DEFAULT false,
    most_influential_variable FLOAT,
    contribution_magnitude FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para mejorar el rendimiento de las consultas
CREATE INDEX idx_anomalias_activo_id ON anomalias(activo_id);
CREATE INDEX idx_anomalias_sensor_id ON anomalias(sensor_id);
CREATE INDEX idx_anomalias_timestamp ON anomalias(timestamp DESC);
CREATE INDEX idx_anomalias_is_anomaly ON anomalias(is_anomaly) WHERE is_anomaly = true;
CREATE INDEX idx_anomalias_severidad ON anomalias(severidad);
CREATE INDEX idx_anomalias_activo_timestamp ON anomalias(activo_id, timestamp DESC);

-- Comentarios para documentación
COMMENT ON TABLE anomalias IS 'Almacena las anomalías detectadas por el sistema de IA en los datos de sensores';
COMMENT ON COLUMN anomalias.activo_id IS 'ID del activo relacionado';
COMMENT ON COLUMN anomalias.sensor_id IS 'Identificador del sensor que generó la anomalía';
COMMENT ON COLUMN anomalias.timestamp IS 'Momento en que se detectó la anomalía';
COMMENT ON COLUMN anomalias.anomaly_score IS 'Puntuación de anomalía calculada por el modelo';
COMMENT ON COLUMN anomalias.anomaly_likelihood IS 'Probabilidad de que sea una anomalía real';
COMMENT ON COLUMN anomalias.severidad IS 'Nivel de gravedad: baja, media, alta, critica';
COMMENT ON COLUMN anomalias.descripcion IS 'Descripción detallada de la anomalía detectada';
COMMENT ON COLUMN anomalias.threshold IS 'Umbral utilizado para detectar la anomalía';
COMMENT ON COLUMN anomalias.is_anomaly IS 'Indica si se confirmó como anomalía verdadera';
COMMENT ON COLUMN anomalias.most_influential_variable IS 'Variable que más contribuyó a la detección de la anomalía';
COMMENT ON COLUMN anomalias.contribution_magnitude IS 'Magnitud de la contribución de la variable más influyente';

-- Trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_anomalias_updated_at BEFORE UPDATE ON anomalias
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();