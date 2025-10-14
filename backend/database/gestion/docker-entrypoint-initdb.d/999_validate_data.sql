-- Script de validación de datos de gestión
-- Archivo: 999_validate_data.sql

-- Mostrar estadísticas de los datos insertados
DO $$
DECLARE
    empresas_count INTEGER;
    edificios_count INTEGER;
    tecnicos_count INTEGER;
    activos_count INTEGER;
    solicitudes_count INTEGER;
    reportes_count INTEGER;
    acciones_count INTEGER;
    tipos_falla_count INTEGER;
    comentarios_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO empresas_count FROM empresas;
    SELECT COUNT(*) INTO edificios_count FROM edificios;
    SELECT COUNT(*) INTO tecnicos_count FROM tecnicos;
    SELECT COUNT(*) INTO activos_count FROM activos;
    SELECT COUNT(*) INTO solicitudes_count FROM solicitudes_tecnico;
    SELECT COUNT(*) INTO reportes_count FROM reportes;
    SELECT COUNT(*) INTO acciones_count FROM acciones_mantenimiento;
    SELECT COUNT(*) INTO tipos_falla_count FROM tipos_falla;
    SELECT COUNT(*) INTO comentarios_count FROM comentarios;
    
    RAISE NOTICE '📊 === ESTADÍSTICAS DE DATOS DE GESTIÓN ===';
    RAISE NOTICE '🏢 Empresas: %', empresas_count;
    RAISE NOTICE '🏗️ Edificios: %', edificios_count;
    RAISE NOTICE '👥 Técnicos: %', tecnicos_count;
    RAISE NOTICE '🔧 Activos: %', activos_count;
    RAISE NOTICE '📝 Solicitudes: %', solicitudes_count;
    RAISE NOTICE '📊 Reportes: %', reportes_count;
    RAISE NOTICE '⚙️ Acciones de Mantenimiento: %', acciones_count;
    RAISE NOTICE '🚨 Tipos de Falla (HU22): %', tipos_falla_count;
    RAISE NOTICE '💬 Comentarios (HU22): %', comentarios_count;
    
    -- Verificar que tenemos datos mínimos
    IF tecnicos_count >= 8 AND activos_count >= 8 AND empresas_count >= 3 AND edificios_count >= 4 
       AND tipos_falla_count >= 15 AND comentarios_count >= 38 THEN
        RAISE NOTICE '✅ Base de datos inicializada correctamente con datos de ejemplo';
    ELSE
        RAISE WARNING '❌ Faltan datos de ejemplo. Verificar scripts de inicialización';
    END IF;
END $$;

-- Mostrar algunos datos de ejemplo
SELECT 'TÉCNICOS' as tabla;
SELECT id, nombre, apellido, especialidad FROM tecnicos LIMIT 5;

SELECT 'ACTIVOS' as tabla;
SELECT id, nombre, tipo FROM activos LIMIT 5;

SELECT 'EMPRESAS' as tabla;
SELECT nombre, rut FROM empresas;
