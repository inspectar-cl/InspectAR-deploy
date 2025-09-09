-- Script para asegurar que el primer activo tenga la información correcta
-- Archivo: 998_fix_first_asset.sql

-- Script para asegurar que los documentos del primer activo tengan la información correcta
-- Archivo: 998_fix_first_asset.sql
-- Nota: En la DB de documentación no existe tabla 'activos', solo documentos asociados a activo_id

-- Verificar que existen documentos para el activo_id = 1 (BombaDeAgua #1 ML)
DO $$
DECLARE
    docs_count INTEGER;
    ficha_tecnica_name VARCHAR(255);
BEGIN
    -- Contar documentos para activo 1
    SELECT COUNT(*) INTO docs_count FROM documentos WHERE activo_id = 1;
    
    -- Obtener nombre de la ficha técnica
    SELECT nombre INTO ficha_tecnica_name 
    FROM documentos 
    WHERE activo_id = 1 AND es_ficha_tecnica = true 
    LIMIT 1;
    
    IF docs_count > 0 THEN
        RAISE NOTICE '✅ Documentación para activo 1 (BombaDeAgua #1 ML): % documentos encontrados', docs_count;
        IF ficha_tecnica_name IS NOT NULL THEN
            RAISE NOTICE '✅ Ficha técnica encontrada: %', ficha_tecnica_name;
        END IF;
    ELSE
        RAISE WARNING '❌ No se encontraron documentos para el activo 1';
    END IF;
END $$;
