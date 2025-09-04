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
BEGIN
    SELECT COUNT(*) INTO empresas_count FROM empresas;
    SELECT COUNT(*) INTO edificios_count FROM edificios;
    SELECT COUNT(*) INTO tecnicos_count FROM tecnicos;
    SELECT COUNT(*) INTO activos_count FROM activos;
    SELECT COUNT(*) INTO solicitudes_count FROM solicitudes_tecnico;
    SELECT COUNT(*) INTO reportes_count FROM reportes;
    SELECT COUNT(*) INTO acciones_count FROM acciones_mantenimiento;
    
    RAISE NOTICE '📊 === ESTADÍSTICAS DE DATOS DE GESTIÓN ===';
    RAISE NOTICE '🏢 Empresas: %', empresas_count;
    RAISE NOTICE '🏗️ Edificios: %', edificios_count;
    RAISE NOTICE '👥 Técnicos: %', tecnicos_count;
    RAISE NOTICE '🔧 Activos: %', activos_count;
    RAISE NOTICE '📝 Solicitudes: %', solicitudes_count;
    RAISE NOTICE '📊 Reportes: %', reportes_count;
    RAISE NOTICE '⚙️ Acciones de Mantenimiento: %', acciones_count;
    
    -- Verificar que tenemos datos mínimos
    IF tecnicos_count >= 8 AND activos_count >= 8 AND empresas_count >= 3 AND edificios_count >= 4 THEN
        RAISE NOTICE '✅ Base de datos inicializada correctamente con datos de ejemplo';
    ELSE
        RAISE WARNING '❌ Faltan datos de ejemplo. Verificar scripts de inicialización';
    END IF;
END $$;

-- Mostrar algunos datos de ejemplo
SELECT 'TÉCNICOS' as tabla;
SELECT id, nombre, apellido, especialidad FROM tecnicos LIMIT 5;

SELECT 'ACTIVOS' as tabla;
SELECT activo_id, nombre, tipo, estado FROM activos LIMIT 5;

SELECT 'EMPRESAS' as tabla;
SELECT nombre, rut FROM empresas;
