-- Script para asegurar que el primer activo tenga la información correcta
-- Archivo: 998_fix_first_asset.sql

-- Actualizar el primer activo para que coincida con la especificación
-- Nota: En la DB de notificaciones solo actualizamos el nombre ya que solo tiene: id, nombre, edificio_id
UPDATE activos SET 
    nombre = 'BombaDeAgua #1 ML'
WHERE id = 1;

-- Verificar que el cambio se aplicó
DO $$
DECLARE
    first_asset_name VARCHAR(255);
    first_asset_edificio_id INTEGER;
BEGIN
    SELECT nombre, edificio_id INTO first_asset_name, first_asset_edificio_id FROM activos WHERE id = 1;
    
    IF first_asset_name = 'BombaDeAgua #1 ML' THEN
        RAISE NOTICE '✅ Primer activo actualizado correctamente: % (edificio_id: %)', first_asset_name, first_asset_edificio_id;
    ELSE
        RAISE WARNING '❌ Error al actualizar primer activo. Nombre actual: %', first_asset_name;
    END IF;
END $$;
