import { useState, useEffect } from "react";
import { Box, Button, TextField, Select, MenuItem } from "@mui/material";
import Services from '@/modules/Services';
import axios from "axios";

const gs = new Services();

function GenerarReporte() {
  const [activos, setActivos] = useState<any[]>([]);
  const [activo, setActivo] = useState("");
  const [observaciones, setObservaciones] = useState("");
  //  CAMBIO: ahora es un array de reportes
  //const [reportes, setReportes] = useState<any[]>([]);
  const [pdfUrl, setPdfUrl] = useState<string | null>(null); // <<< CAMBIO: estado para vista previa PDF


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
          //console.warn("La respuesta no tiene activos válidos");
          setActivos([]);
        }
      } catch (error) {
        //console.error("Error al conectar con API Gateway:", error);
        setActivos([]);
      }
    };
    fetchActivos();
  }, []);

  // Exportar PDF
  // Este método usa axios porque el gs tiene limitaciones con blobs (y otros tipos de respuesta)
  const handleExportarPDF = async () => {
  if (!activo) return;

  try {
    const response = await axios.post(
      `/api/gestion/reportes/activo/${activo}`,
      {},
      { responseType: "blob" } // importante
    );

    const blob = new Blob([response.data], { type: "application/pdf" });
    const url = window.URL.createObjectURL(blob);

    // Generar nombre dinámico: "Reporte_AscensorNorte_20250908.pdf"
    const activoSeleccionado = activos.find(a => a.id === activo);
    const nombreActivo = activoSeleccionado ? activoSeleccionado.nombre.replace(/\s+/g, "_") : "Reporte";
    const fecha = new Date().toISOString().split("T")[0].replace(/-/g, "");
    const nombreArchivo = `Reporte_${nombreActivo}_${fecha}.pdf`;

    // Crear link temporal y hacer click
    const link = document.createElement("a");
    link.href = url;
    link.download = nombreArchivo;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    // Liberar objeto URL
    setTimeout(() => { window.URL.revokeObjectURL(url); }, 10000);

  } catch (error) {
    //console.error("Error al exportar PDF:", error);
    alert("No se pudo generar el PDF");
  }
};

  // <<< CAMBIO: función para vista previa PDF
  const handleVistaPreviaPDF = async () => {
    if (!activo) return;
    try {
      const response = await axios.post(
        `/api/gestion/reportes/activo/${activo}`,
        {},
        { responseType: "blob" }
      );
      const blob = new Blob([response.data], { type: "application/pdf" });
      const url = window.URL.createObjectURL(blob);
      setPdfUrl(url);
    } catch (error) {
      //console.error("Error al generar vista previa PDF:", error);
      alert("No se pudo generar la vista previa");
    }
  };

  // Función para actualizar observaciones
// Función para actualizar observaciones del analista
{/*const handleActualizarObservaciones = async () => {
  if (!reportes || reportes.length === 0) {
    console.log("No hay reportes para actualizar"); // <<< log agregado
    return;
  }

  const reporteId = reportes[0].id;
  console.log("Intentando actualizar observaciones para reporte ID:", reporteId);
  console.log("Observaciones actuales:", observaciones);

  try {
    const data = await gs.put(`/reportes/${reporteId}/observaciones`, {
      observaciones_analista: observaciones, // ✅ campo correcto
      autor_analista: "Tu Nombre"            // ✅ remplazar por quien actualiza
    });
    console.log("Respuesta API:", data);
    alert("Observaciones actualizadas correctamente");
  } catch (error) {
    console.error("Error al actualizar observaciones:", error);
    alert("No se pudo actualizar las observaciones");
  }
};*/}



  return (
    <Box>
      {/* Filtros */}
      <Box sx={{ display: "flex", gap: 2, mb: 3 }}>
        <Select
          value={activo}
          onChange={(e) => { setActivo(e.target.value); }}
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
        <Button variant="contained" onClick={handleVistaPreviaPDF}>Vista Previa PDF</Button>
        {/* <Button variant="contained" onClick={handleGenerar}>Generar</Button> */}
        <Button variant="contained" onClick={handleExportarPDF}>Exportar PDF</Button>
      </Box>

      {/* <<< CAMBIO: iframe para vista previa PDF */}
      {pdfUrl ? <Box sx={{ mt: 3, borderRadius: 2 }}>
          <iframe
            src={pdfUrl}
            style={{ width: "100%", height: "500px" }}
            title="Vista Previa PDF"
          />
        </Box> : null}

      {/* Observaciones */}
      <Box sx={{ mt: 2 }}>
        <TextField
          label="Observaciones"
          multiline
          rows={3}
          fullWidth
          value={observaciones}
          onChange={(e) => { setObservaciones(e.target.value); }}
        />
        <Button 
    variant="contained" 
    sx={{ mt: 1 }}
    //onClick={handleActualizarObservaciones}
  >
    Guardar Observaciones
  </Button>
      </Box>
    </Box>
  );
}

export default GenerarReporte;
