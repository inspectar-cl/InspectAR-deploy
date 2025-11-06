-- Eliminar trigger si existe antes de crearlo
DROP TRIGGER IF EXISTS trigger_update_numero_activos ON activos;

CREATE TRIGGER trigger_update_numero_activos
    AFTER INSERT OR DELETE ON activos
    FOR EACH ROW EXECUTE FUNCTION update_numero_activos();
