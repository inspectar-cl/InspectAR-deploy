-- Schema para el microservicio de documentación
-- Base de datos: documentacion_db

-- Tabla de documentos (HdU05, HdU23)
CREATE TABLE IF NOT EXISTS documentos (
    id SERIAL PRIMARY KEY,
    activo_id INTEGER NOT NULL, -- Referencia al activo (del microservicio de gestión)
    tecnico_id INTEGER, -- Referencia al técnico que subió el documento (opcional)
    nombre VARCHAR(255) NOT NULL,
    descripcion TEXT,
    categoria VARCHAR(50) NOT NULL CHECK (categoria IN ('ficha_tecnica', 'informe_mantenimiento', 'diagnostico', 'manual_fabricante', 'certificacion')),
    tipo_archivo VARCHAR(10) NOT NULL CHECK (tipo_archivo IN ('pdf', 'docx')),
    ruta_archivo VARCHAR(500) NOT NULL, -- Ruta en el storage (MinIO o local)
    tamano_bytes BIGINT NOT NULL,
    fecha_emision DATE NOT NULL,
    subido_por VARCHAR(100) NOT NULL, -- Usuario que subió el documento
    palabras_clave TEXT, -- Separadas por comas para búsqueda
    es_ficha_tecnica BOOLEAN DEFAULT FALSE, -- Para identificar fichas técnicas (HdU23)
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    actualizado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de análisis de IA (HdU19 - Gemini)
CREATE TABLE IF NOT EXISTS analisis_ia (
    id SERIAL PRIMARY KEY,
    documento_id INTEGER NOT NULL REFERENCES documentos(id) ON DELETE CASCADE,
    resumen TEXT, -- Resumen generado por IA
    puntos_claves JSONB, -- Array JSON de puntos importantes
    graficos JSONB, -- Array JSON de descripciones de gráficos
    estado VARCHAR(20) NOT NULL DEFAULT 'procesando' CHECK (estado IN ('procesando', 'completado', 'error')),
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices para mejorar rendimiento
CREATE INDEX IF NOT EXISTS idx_documentos_activo_id ON documentos(activo_id);
CREATE INDEX IF NOT EXISTS idx_documentos_categoria ON documentos(categoria);
CREATE INDEX IF NOT EXISTS idx_documentos_fecha_emision ON documentos(fecha_emision);
CREATE INDEX IF NOT EXISTS idx_documentos_tecnico_id ON documentos(tecnico_id);
CREATE INDEX IF NOT EXISTS idx_documentos_es_ficha_tecnica ON documentos(es_ficha_tecnica);
CREATE INDEX IF NOT EXISTS idx_documentos_palabras_clave ON documentos USING gin(to_tsvector('spanish', palabras_clave));

-- Índice para búsqueda de texto completo
CREATE INDEX IF NOT EXISTS idx_documentos_busqueda ON documentos USING gin(
    to_tsvector('spanish', nombre || ' ' || COALESCE(descripcion, '') || ' ' || COALESCE(palabras_clave, ''))
);

CREATE INDEX IF NOT EXISTS idx_analisis_documento_id ON analisis_ia(documento_id);
CREATE INDEX IF NOT EXISTS idx_analisis_estado ON analisis_ia(estado);

-- Trigger para actualizar timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.actualizado_en = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_documentos_updated_at BEFORE UPDATE ON documentos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insertar datos de ejemplo
INSERT INTO documentos (
    activo_id, nombre, descripcion, categoria, tipo_archivo, 
    ruta_archivo, tamano_bytes, fecha_emision, subido_por, 
    palabras_clave, es_ficha_tecnica
) VALUES 
(1, 'Ficha Técnica Caldera Principal', 'Especificaciones técnicas oficiales del fabricante', 'ficha_tecnica', 'pdf', 'ficha_caldera_principal.pdf', 2048576, '2024-01-15', 'admin', 'caldera,especificaciones,fabricante', true),
(1, 'Informe Mantenimiento Enero 2024', 'Reporte de mantenimiento preventivo realizado', 'informe_mantenimiento', 'pdf', 'informe_enero_2024.pdf', 1024576, '2024-01-30', 'tecnico_juan', 'mantenimiento,preventivo,enero', false),
(2, 'Manual del Fabricante - Bomba Hidráulica', 'Manual completo de operación y mantenimiento', 'manual_fabricante', 'pdf', 'manual_bomba_hidraulica.pdf', 5242880, '2023-12-10', 'admin', 'bomba,hidraulica,manual,operacion', false),
(3, 'Diagnóstico Sistema Eléctrico', 'Análisis completo del sistema eléctrico del activo', 'diagnostico', 'pdf', 'diagnostico_electrico.pdf', 3145728, '2024-02-05', 'tecnico_maria', 'diagnostico,electrico,sistema,analisis', false);

-- Comentarios para documentación
COMMENT ON TABLE documentos IS 'Almacena todos los documentos técnicos asociados a activos industriales';
COMMENT ON COLUMN documentos.categoria IS 'Categoría del documento: ficha_tecnica, informe_mantenimiento, diagnostico, manual_fabricante, certificacion';
COMMENT ON COLUMN documentos.es_ficha_tecnica IS 'Indica si el documento es una ficha técnica oficial del fabricante (HdU23)';
COMMENT ON COLUMN documentos.palabras_clave IS 'Palabras clave separadas por comas para facilitar búsquedas';

COMMENT ON TABLE analisis_ia IS 'Almacena los análisis realizados por IA (Gemini) sobre los documentos';
COMMENT ON COLUMN analisis_ia.puntos_claves IS 'Array JSON con puntos importantes extraídos por IA';
COMMENT ON COLUMN analisis_ia.graficos IS 'Array JSON con descripciones de gráficos y diagramas detectados por IA';
