import { useState, useEffect } from "react";
import { Box, Button, TextField, Typography, Select, MenuItem } from "@mui/material";
import Services from '@/modules/Services';

const gs = new Services();

const GenerarReporte = () => {
  const [activos, setActivos] = useState([]);
  const [activo, setActivo] = useState("");
  const [observaciones, setObservaciones] = useState("");
  const [reporte, setReporte] = useState(null);

  // Cargar activos al montar el componente
  useEffect(() => {
    const fetchActivos = async () => {
      try {
        const data = await gs.get("/obtener-activos");

        // Si la API devuelve directamente un array de activos
        if (Array.isArray(data)) {
          setActivos(data);
        } 
        // Si la API devuelve un objeto con la propiedad 'activos'
        else if (data && Array.isArray(data.activos)) {
          setActivos(data.activos);
        } 
        // Si no devuelve nada útil
        else {
          console.warn("La respuesta no tiene activos válidos");
          setActivos([]);
        }
      } catch (error) {
        console.error("Error al conectar con API Gateway:", error);
        setActivos([]);
      }
    };
    fetchActivos();
  }, []);


  // Generar reporte
  const handleGenerar = async () => {
    if (!activo) return;
    // ======= LLAMADA API PARA GENERAR REPORTE =======
    try {
      console.log("Generando reporte para activo ID:", activo);
      const data = await gs.get(`/gestion/reportes/activo/${activo}`, { // await gs.get(`/reportes/activo/${activo}`
      });
      if (Array.isArray(data) && data.length > 0) {
        setReporte(data[0]); // Tomar el primer reporte
      } else {
        setReporte(null);
      }
    } catch (error) {
      console.error("Error al generar reporte:", error);
      setReporte(null);
    }
  };


  // Exportar PDF
  const handleExportarPDF = async () => {
    if (!activo) return;

    // ======= LLAMADA API PARA GENERAR PDF =======
    try {
      const response = await gs.post(`/gestion/reportes/activo/1`, {
      });

      // Lógica para descargar PDF
    } catch (error) {
      console.error("Error al exportar PDF:", error);
    }
  };

  return (
    <Box>
      {/* Filtros */}
      <Box sx={{ display: "flex", gap: 2, mb: 3 }}>
        <Select
          value={activo}
          onChange={(e) => setActivo(e.target.value)}
          displayEmpty
        >
          <MenuItem value="">Seleccionar Activo</MenuItem>
          {Array.isArray(activos) &&
            activos.map((a) => (
              <MenuItem key={a.id} value={a.id}>
                {a.nombre}
              </MenuItem>
            ))}
        </Select>
        <Button variant="contained" onClick={handleGenerar}>Generar</Button>
        <Button variant="contained" onClick={handleExportarPDF}>Exportar PDF</Button>
      </Box>

      <Box sx={{ border: "1px solid #ddd", borderRadius: 2, p: 2, minHeight: 200 }}>
        {reporte ? (
          <>
            {/* Nombre del activo y tipo de reporte */}
            <Typography variant="h6">
              {reporte.activo.nombre} - {reporte.tipo_reporte}
            </Typography>

            {/* Información adicional */}
            <Typography variant="body2" sx={{ mt: 1 }}>
              Generado en: {new Date(reporte.generado_en).toLocaleString()}
            </Typography>
            <Typography variant="body2" sx={{ mt: 1 }}>
              Estado: {reporte.estado}
            </Typography>
            <Typography variant="body2" sx={{ mt: 1 }}>
              Ubicación: {reporte.activo.ubicacion} | Tipo: {reporte.activo.tipo}
            </Typography>
            <Typography variant="body2" sx={{ mt: 1 }}>
              Edificio: {reporte.activo.edificio.nombre} - {reporte.activo.edificio.direccion}
            </Typography>

            {/* Contenido del reporte */}
            <Typography variant="body2" sx={{ mt: 1 }}>
              Contenido: {reporte.contenido}
            </Typography>
          </>
        ) : (
          <Typography variant="body2" color="text.secondary">
            [Aquí se mostrará la vista previa del reporte]
          </Typography>
        )}
      </Box>


      {/* Observaciones */}
      <Box sx={{ mt: 2 }}>
        <TextField
          label="Observaciones"
          multiline
          rows={3}
          fullWidth
          value={observaciones}
          onChange={(e) => setObservaciones(e.target.value)}
        />
      </Box>
    </Box>
  );
};

export default GenerarReporte;
