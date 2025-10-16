-- Configuración inicial de la base de datos
-- Este script siempre se ejecuta al inicializar el contenedor

-- Crear usuario y base de datos si no existen
DO $$
BEGIN
    -- Verificar si la base de datos existe, si no, crearla
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gestion_db') THEN
        PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE gestion_db');
    END IF;
EXCEPTION WHEN OTHERS THEN
    -- Si hay error (probablemente porque la base ya existe), continuar
    NULL;
END $$;

-- Asegurar que el usuario tenga todos los permisos
GRANT ALL PRIVILEGES ON DATABASE gestion_db TO gestion_user;
