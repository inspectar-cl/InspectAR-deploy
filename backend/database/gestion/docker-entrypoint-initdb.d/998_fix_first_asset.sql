-- Script para asegurar que el primer activo tenga la información correcta
-- Archivo: 998_fix_first_asset.sql

-- Actualizar el primer activo para que coincida con la especificación
UPDATE activos SET 
    nombre = 'BombaDeAgua #1 ML',
    tipo = 'bomba de agua',
    ubicacion = 'Planta B',
    estado = 'operativo'
WHERE id = 1;

-- Verificar que el cambio se aplicó
DO $$
DECLARE
    first_asset_name VARCHAR(255);
    first_asset_type VARCHAR(100);
BEGIN
    SELECT nombre, tipo INTO first_asset_name, first_asset_type FROM activos WHERE id = 1;
    
    IF first_asset_name = 'BombaDeAgua #1 ML' AND first_asset_type = 'bomba de agua' THEN
        RAISE NOTICE '✅ Primer activo actualizado correctamente: % (tipo: %)', first_asset_name, first_asset_type;
    ELSE
        RAISE WARNING '❌ Error al actualizar primer activo. Nombre: %, Tipo: %', first_asset_name, first_asset_type;
    END IF;
END $$;
