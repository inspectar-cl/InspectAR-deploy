-- Base de datos para microservicio de documentación
-- Microservicio: Documentación Técnica
-- Puerto: 5434

-- Tabla principal de documentos técnicos
CREATE TABLE IF NOT EXISTS documentos (
    id SERIAL PRIMARY KEY,
    activo_id INTEGER NOT NULL, -- Referencia al activo (del microservicio gestión)
    tecnico_id INTEGER, -- Referencia al técnico que subió (opcional)
    nombre VARCHAR(255) NOT NULL,
    descripcion TEXT,
    categoria VARCHAR(50) NOT NULL CHECK (categoria IN ('ficha_tecnica', 'manual_fabricante', 'reporte_mantenimiento', 'diagnostico', 'certificacion')),
    tipo_archivo VARCHAR(10) NOT NULL CHECK (tipo_archivo IN ('pdf', 'docx')),
    ruta_archivo VARCHAR(500) NOT NULL, -- Ruta en MinIO o storage local
    tamano_bytes BIGINT NOT NULL,
    fecha_emision DATE,
    subido_por VARCHAR(100) NOT NULL,
    palabras_clave TEXT, -- Separadas por comas para búsqueda
    es_ficha_tecnica BOOLEAN DEFAULT FALSE, -- Marca especial para fichas técnicas
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    actualizado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla para análisis de IA (resultados de Gemini)
CREATE TABLE IF NOT EXISTS analisis_ia (
    id SERIAL PRIMARY KEY,
    documento_id INTEGER NOT NULL REFERENCES documentos(id) ON DELETE CASCADE,
    resumen TEXT NOT NULL,
    puntos_claves JSONB, -- Array JSON de puntos importantes extraídos por IA
    graficos_detectados JSONB, -- Array JSON de descripciones de gráficos
    estado VARCHAR(20) NOT NULL DEFAULT 'procesando' CHECK (estado IN ('procesando', 'completado', 'error')),
    procesado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla para consultas interactivas con IA (NUEVA FUNCIONALIDAD)
CREATE TABLE IF NOT EXISTS consultas_ia (
    id SERIAL PRIMARY KEY,
    documento_id INTEGER NOT NULL REFERENCES documentos(id) ON DELETE CASCADE,
    pregunta TEXT NOT NULL,
    respuesta TEXT NOT NULL,
    confianza DECIMAL(3,2) CHECK (confianza >= 0 AND confianza <= 1),
    fuentes JSONB, -- Array JSON con las fuentes/secciones del documento
    tiempo_respuesta_ms INTEGER, -- Tiempo de procesamiento en milisegundos
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices para optimizar consultas
CREATE INDEX IF NOT EXISTS idx_documentos_activo_id ON documentos(activo_id);
CREATE INDEX IF NOT EXISTS idx_documentos_categoria ON documentos(categoria);
CREATE INDEX IF NOT EXISTS idx_documentos_tecnico_id ON documentos(tecnico_id);
CREATE INDEX IF NOT EXISTS idx_documentos_fecha_emision ON documentos(fecha_emision);
CREATE INDEX IF NOT EXISTS idx_documentos_es_ficha_tecnica ON documentos(es_ficha_tecnica);

-- Índice compuesto para fichas técnicas por activo (consulta más frecuente)
CREATE INDEX IF NOT EXISTS idx_documentos_activo_ficha ON documentos(activo_id, es_ficha_tecnica) WHERE es_ficha_tecnica = true;

-- Agregar columna para búsqueda de texto completo
ALTER TABLE documentos ADD COLUMN IF NOT EXISTS texto_busqueda tsvector;

-- Trigger para actualizar automáticamente el vector de búsqueda
CREATE OR REPLACE FUNCTION actualizar_texto_busqueda() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.texto_busqueda := to_tsvector('spanish', 
        COALESCE(NEW.nombre, '') || ' ' || 
        COALESCE(NEW.descripcion, '') || ' ' || 
        COALESCE(NEW.palabras_clave, '')
    );
    NEW.actualizado_en := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger
DROP TRIGGER IF EXISTS trigger_texto_busqueda ON documentos;
CREATE TRIGGER trigger_texto_busqueda 
    BEFORE INSERT OR UPDATE ON documentos 
    FOR EACH ROW EXECUTE FUNCTION actualizar_texto_busqueda();

-- Índice GIN para búsqueda de texto completo súper rápida
CREATE INDEX IF NOT EXISTS idx_documentos_texto_busqueda ON documentos USING GIN(texto_busqueda);

-- Índices para tabla de análisis IA
CREATE INDEX IF NOT EXISTS idx_analisis_documento_id ON analisis_ia(documento_id);
CREATE INDEX IF NOT EXISTS idx_analisis_estado ON analisis_ia(estado, procesado_en);

-- Índices para tabla de consultas IA (NUEVA)
CREATE INDEX IF NOT EXISTS idx_consultas_documento_id ON consultas_ia(documento_id);
CREATE INDEX IF NOT EXISTS idx_consultas_fecha ON consultas_ia(creado_en);
CREATE INDEX IF NOT EXISTS idx_consultas_confianza ON consultas_ia(confianza) WHERE confianza IS NOT NULL;

-- Insertar datos de ejemplo
INSERT INTO documentos (
    activo_id, tecnico_id, nombre, descripcion, categoria, 
    tipo_archivo, ruta_archivo, tamano_bytes, fecha_emision, 
    subido_por, palabras_clave, es_ficha_tecnica
) VALUES 
-- Ficha técnica de caldera
(1, NULL, 'Ficha Técnica Caldera Principal Bosch', 'Especificaciones técnicas oficiales del fabricante para caldera industrial de 500kW', 'ficha_tecnica', 'pdf', '1_ficha_tecnica_caldera_bosch.pdf', 2048576, '2024-01-15', 'admin_sistema', 'caldera,bosch,ficha tecnica,500kw,especificaciones,fabricante', true),

-- Manual de operación
(1, 2, 'Manual de Operación y Mantenimiento', 'Manual completo de operación, mantenimiento preventivo y correctivo', 'manual_fabricante', 'pdf', '1_manual_operacion_caldera.pdf', 5242880, '2024-01-15', 'tecnico_juan', 'manual,operacion,mantenimiento,caldera,procedimientos', false),

-- Reporte de mantenimiento
(1, 1, 'Reporte Mantenimiento Preventivo Enero 2024', 'Inspección mensual realizada según programa de mantenimiento', 'reporte_mantenimiento', 'pdf', '1_reporte_enero_2024.pdf', 1048576, '2024-01-30', 'tecnico_maria', 'mantenimiento,preventivo,enero,inspeccion,caldera', false),

-- Certificación de seguridad
(1, NULL, 'Certificación de Seguridad Industrial', 'Certificado de cumplimiento de normas de seguridad industrial', 'certificacion', 'pdf', '1_certificacion_seguridad.pdf', 815432, '2024-02-01', 'inspector_seguridad', 'certificacion,seguridad,industrial,normas,cumplimiento', false),

-- Documentos para bomba hidráulica (activo 2)
(2, 3, 'Ficha Técnica Bomba Hidráulica Grundfos', 'Especificaciones técnicas bomba centrífuga 50HP', 'ficha_tecnica', 'pdf', '2_ficha_tecnica_bomba_grundfos.pdf', 1524288, '2023-12-10', 'admin_sistema', 'bomba,hidraulica,grundfos,ficha tecnica,50hp,centrifuga', true),

-- Diagnóstico de falla
(2, 2, 'Diagnóstico de Vibración Excesiva', 'Análisis de causa raíz para vibración anormal en bomba', 'diagnostico', 'pdf', '2_diagnostico_vibracion.pdf', 2097152, '2024-02-15', 'tecnico_carlos', 'diagnostico,vibracion,bomba,analisis,causa raiz,falla', false);

-- Insertar algunos análisis de IA de ejemplo
INSERT INTO analisis_ia (documento_id, resumen, puntos_claves, graficos_detectados, estado) VALUES 
(1, 
 'Ficha técnica de caldera industrial Bosch de 500kW con especificaciones completas de operación, mantenimiento y seguridad. Incluye parámetros críticos de funcionamiento y requisitos de instalación.',
 '["Capacidad térmica: 500kW", "Presión máxima operación: 150 PSI", "Temperatura máxima: 180°C", "Combustible: Gas natural/GLP", "Eficiencia energética: 92%", "Mantenimiento preventivo cada 90 días", "Certificación ISO 9001 y CE"]',
 '["Diagrama esquemático del sistema de combustión", "Tabla de especificaciones técnicas principales", "Gráfico de curva de eficiencia vs carga", "Plano de dimensiones y conexiones"]',
 'completado'),

(3,
 'Reporte de mantenimiento preventivo que documenta la inspección completa de la caldera, incluyendo verificación de sistemas de seguridad, limpieza de componentes y calibración de instrumentos.',
 '["Estado general: Satisfactorio", "Presión de prueba: 145 PSI (Normal)", "Limpieza de quemadores realizada", "Calibración de termostatos completada", "Reemplazo de filtros de aire", "Próxima inspección: Abril 2024"]',
 '["Tabla de mediciones de presión", "Checklist de inspección visual", "Fotografías de componentes inspeccionados"]',
 'completado');

-- Insertar consultas de ejemplo para demostrar funcionalidad
INSERT INTO consultas_ia (documento_id, pregunta, respuesta, confianza, fuentes, tiempo_respuesta_ms) VALUES
(1,
 '¿Cuál es la capacidad térmica de esta caldera?',
 'La caldera Bosch tiene una capacidad térmica de 500kW. Esta información se encuentra especificada en la tabla de características técnicas principales del documento.',
 0.98,
 '["Tabla de especificaciones técnicas principales", "Sección 2.1 - Características técnicas"]',
 1250),

(1,
 '¿Qué eficiencia energética tiene este equipo?',
 'La caldera tiene una eficiencia energética del 92%, lo que la clasifica como un equipo de alta eficiencia. Está certificada bajo normas ISO 9001 y CE.',
 0.96,
 '["Sección de eficiencia energética", "Certificaciones"]',
 980),

(3,
 '¿Cuándo fue la última mantención?',
 'La última mantención preventiva fue realizada el 30 de enero de 2024. Se ejecutaron tareas de limpieza de quemadores, calibración de termostatos y reemplazo de filtros de aire. El estado general del equipo fue evaluado como satisfactorio.',
 0.95,
 '["Fecha del reporte", "Checklist de actividades realizadas", "Evaluación del estado general"]',
 1450),

(5,
 '¿Cuánta presión maneja esta bomba?',
 'La bomba hidráulica Grundfos de 50HP tiene una presión máxima de operación de 150 PSI. Se recomienda operar entre 120-140 PSI para óptimo rendimiento y mayor vida útil del equipo.',
 0.97,
 '["Tabla de parámetros operacionales", "Especificaciones de presión"]',
 1100);

-- Comentarios para documentación
COMMENT ON TABLE documentos IS 'Almacena metadatos de todos los documentos técnicos asociados a activos industriales';
COMMENT ON COLUMN documentos.activo_id IS 'ID del activo al que pertenece el documento (referencia externa al microservicio gestión)';
COMMENT ON COLUMN documentos.categoria IS 'Tipo de documento: ficha_tecnica, manual_fabricante, reporte_mantenimiento, diagnostico, certificacion';
COMMENT ON COLUMN documentos.es_ficha_tecnica IS 'Marca especial para identificar rápidamente las fichas técnicas oficiales';
COMMENT ON COLUMN documentos.texto_busqueda IS 'Vector de búsqueda de texto completo generado automáticamente';

COMMENT ON TABLE consultas_ia IS 'Consultas interactivas realizadas por usuarios sobre documentos específicos usando IA';
COMMENT ON COLUMN consultas_ia.confianza IS 'Nivel de confianza de la respuesta (0.0-1.0)';
COMMENT ON COLUMN consultas_ia.fuentes IS 'Array JSON con las secciones del documento de donde se extrajo la información';
COMMENT ON COLUMN consultas_ia.tiempo_respuesta_ms IS 'Tiempo de procesamiento de la consulta en milisegundos';

-- Función para buscar consultas similares (evitar preguntas repetidas)
CREATE OR REPLACE FUNCTION buscar_consultas_similares(doc_id INTEGER, pregunta_texto TEXT)
RETURNS TABLE(
    id INTEGER,
    pregunta TEXT,
    respuesta TEXT,
    confianza DECIMAL,
    similaridad REAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        c.id,
        c.pregunta,
        c.respuesta,
        c.confianza,
        similarity(c.pregunta, pregunta_texto) as similaridad
    FROM consultas_ia c
    WHERE c.documento_id = doc_id
      AND similarity(c.pregunta, pregunta_texto) > 0.3
    ORDER BY similaridad DESC
    LIMIT 5;
END;
$$ LANGUAGE plpgsql;

-- Función para buscar documentos por texto
CREATE OR REPLACE FUNCTION buscar_documentos_texto(termino_busqueda TEXT)
RETURNS TABLE(
    id INTEGER,
    activo_id INTEGER,
    nombre VARCHAR(255),
    categoria VARCHAR(50),
    relevancia REAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        d.id,
        d.activo_id,
        d.nombre,
        d.categoria,
        ts_rank(d.texto_busqueda, plainto_tsquery('spanish', termino_busqueda)) as relevancia
    FROM documentos d
    WHERE d.texto_busqueda @@ plainto_tsquery('spanish', termino_busqueda)
    ORDER BY relevancia DESC;
END;
$$ LANGUAGE plpgsql;
