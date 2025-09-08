-- Script para agregar columnas de observaciones editables si no existen
-- Este script se puede ejecutar sin problemas múltiples veces

-- Agregar columna de observaciones del revisor si no existe
DO $$ 
BEGIN
    IF NOT EXISTS(SELECT * FROM information_schema.columns WHERE table_name='reportes' AND column_name='observaciones_revision') THEN
        ALTER TABLE reportes ADD COLUMN observaciones_revision TEXT DEFAULT '';
        RAISE NOTICE 'Columna observaciones_revision agregada';
    ELSE
        RAISE NOTICE 'Columna observaciones_revision ya existe';
    END IF;
END $$;

-- Crear índice adicional si no existe
CREATE INDEX IF NOT EXISTS idx_reportes_fecha_revision ON reportes(fecha_revision);

-- Verificar estructura final
SELECT 
    column_name, 
    data_type, 
    is_nullable, 
    column_default
FROM information_schema.columns 
WHERE table_name = 'reportes' 
AND column_name IN (
    'observaciones_analista', 'autor_analista', 'estructura_informe', 
    'metadata_informe', 'version_reporte', 'estado_revision', 
    'fecha_revision', 'revisor', 'observaciones_revision'
)
ORDER BY column_name;

-- Mostrar resumen
SELECT 
    'Verificación completada: Columnas de observaciones editables' as resultado,
    COUNT(*) as total_reportes,
    COUNT(CASE WHEN observaciones_analista IS NOT NULL AND observaciones_analista != '' THEN 1 END) as reportes_con_observaciones
FROM reportes;
