-- Insertar firmas digitales de ejemplo
-- Archivo: 003_insert_firmas.sql
-- 
-- IMPORTANTE: Este script asume que los archivos de firma ya están en el directorio
-- /app/storage/firmas/ del contenedor de gestion-service

-- Función auxiliar para insertar firmas con datos binarios
-- Las rutas de los archivos deben coincidir con la ubicación real en el contenedor

-- Usuario 1 (analista@example.com) - ID: 1
INSERT INTO firmas_digitales 
(usuario_id, nombre_archivo, ruta_archivo, tipo_mime, formato, tamano_bytes, es_predeterminada)
VALUES 
(1, 'Firma Analista', '/app/storage/firmas/firma_usuario_1.jpg', 'image/jpeg', 'jpg', 0, true)
ON CONFLICT DO NOTHING;

-- Usuario 2 (tecnico@example.com) - ID: 2
INSERT INTO firmas_digitales 
(usuario_id, nombre_archivo, ruta_archivo, tipo_mime, formato, tamano_bytes, es_predeterminada)
VALUES 
(2, 'Firma Técnico', '/app/storage/firmas/firma_usuario_2.jpg', 'image/jpeg', 'jpg', 0, true)
ON CONFLICT DO NOTHING;

-- Usuario 3 (residente@example.com) - ID: 3
INSERT INTO firmas_digitales 
(usuario_id, nombre_archivo, ruta_archivo, tipo_mime, formato, tamano_bytes, es_predeterminada)
VALUES 
(3, 'Firma Residente', '/app/storage/firmas/firma_usuario_3.jpg', 'image/jpeg', 'jpg', 0, true)
ON CONFLICT DO NOTHING;

-- Usuario 4 (admin@example.com) - ID: 4
INSERT INTO firmas_digitales 
(usuario_id, nombre_archivo, ruta_archivo, tipo_mime, formato, tamano_bytes, es_predeterminada)
VALUES 
(4, 'Firma Administrador', '/app/storage/firmas/firma_usuario_4.jpg', 'image/jpeg', 'jpg', 0, true)
ON CONFLICT DO NOTHING;

-- Actualizar tamaños de archivos después de la inserción
-- Esto es solo un placeholder, los tamaños reales se actualizarán cuando se suban archivos reales
UPDATE firmas_digitales SET tamano_bytes = 50000 WHERE tamano_bytes = 0;

-- Mensaje de confirmación
DO $$
BEGIN
    RAISE NOTICE 'Firmas digitales de ejemplo insertadas correctamente';
    RAISE NOTICE 'Total de firmas: %', (SELECT COUNT(*) FROM firmas_digitales);
END $$;
