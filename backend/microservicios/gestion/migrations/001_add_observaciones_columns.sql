-- Migración para agregar columnas de observaciones editables
-- Sistema de Reportes con Observaciones - Microservicio de Gestión
-- Fecha: 2025-09-08

-- Agregar columnas para observaciones editables
ALTER TABLE reportes 
ADD COLUMN IF NOT EXISTS observaciones_analista TEXT,
ADD COLUMN IF NOT EXISTS autor_analista VARCHAR(255),
ADD COLUMN IF NOT EXISTS estructura_informe JSONB,
ADD COLUMN IF NOT EXISTS metadata_informe JSONB,
ADD COLUMN IF NOT EXISTS version_reporte INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS estado_revision VARCHAR(50) DEFAULT 'pendiente',
ADD COLUMN IF NOT EXISTS fecha_revision TIMESTAMP,
ADD COLUMN IF NOT EXISTS revisor VARCHAR(255),
ADD COLUMN IF NOT EXISTS observaciones_revision TEXT;

-- Crear índices para mejor performance
CREATE INDEX IF NOT EXISTS idx_reportes_activo_observaciones 
ON reportes(activo_id) WHERE observaciones_analista IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_reportes_estado_revision 
ON reportes(estado_revision);

CREATE INDEX IF NOT EXISTS idx_reportes_autor_analista 
ON reportes(autor_analista);

CREATE INDEX IF NOT EXISTS idx_reportes_fecha_revision 
ON reportes(fecha_revision);

-- Agregar constraint para estados de revisión válidos
ALTER TABLE reportes 
ADD CONSTRAINT IF NOT EXISTS check_estado_revision 
CHECK (estado_revision IN ('pendiente', 'en_revision', 'aprobado', 'rechazado'));

-- Comentarios en las columnas para documentación
COMMENT ON COLUMN reportes.observaciones_analista IS 'Observaciones editables del analista sobre el reporte';
COMMENT ON COLUMN reportes.autor_analista IS 'Nombre del analista que realizó las observaciones';
COMMENT ON COLUMN reportes.estructura_informe IS 'Estructura del informe en formato JSON';
COMMENT ON COLUMN reportes.metadata_informe IS 'Metadata adicional del informe en formato JSON';
COMMENT ON COLUMN reportes.version_reporte IS 'Versión del reporte para control de cambios';
COMMENT ON COLUMN reportes.estado_revision IS 'Estado de revisión del reporte';
COMMENT ON COLUMN reportes.fecha_revision IS 'Fecha de la última revisión del reporte';
COMMENT ON COLUMN reportes.revisor IS 'Nombre de la persona que realizó la revisión';
COMMENT ON COLUMN reportes.observaciones_revision IS 'Observaciones del revisor sobre el reporte';

-- Actualizar reportes existentes con valores por defecto
UPDATE reportes 
SET 
    version_reporte = 1,
    estado_revision = 'pendiente'
WHERE version_reporte IS NULL OR estado_revision IS NULL;

-- Mostrar información de la migración
SELECT 
    'Migración completada: Columnas de observaciones agregadas a la tabla reportes' as resultado,
    count(*) as reportes_actualizados
FROM reportes;
