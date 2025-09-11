/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
import {useState, useEffect } from "react";
import { Box, Button, TextField, Select, MenuItem, Snackbar, Alert, FormControl, InputLabel, Checkbox, ListItemText, OutlinedInput } from "@mui/material";
import Services from '@/modules/Services';
import axios from "axios";

const gs = new Services();
const BASE_URL = process.env.NEXT_PUBLIC_API_GATEWAY_URL || '/api'

interface ActivoReporte {
  id: string;
  nombre: string;
}

interface ActivosResponse {
  activos?: ActivoReporte[];
}

function GenerarReporte() {
  const [activos, setActivos] = useState<ActivoReporte[]>([]);
  const [activo, setActivo] = useState("");
  const [observaciones, setObservaciones] = useState("");
  //  CAMBIO: ahora es un array de reportes
  //const [reportes, setReportes] = useState<any[]>([]);
  const [pdfUrl, setPdfUrl] = useState<string | null>(null); // <<< CAMBIO: estado para vista previa PDF
  const [mensaje, setMensaje] = useState<string | null>(null);
  
  // Nuevo estado para los campos seleccionables del reporte
  const [camposSeleccionados, setCamposSeleccionados] = useState<string[]>([
    "ubicacion",
    "historial_mantenimientos", 
    "ultima_acciones",
    "datos_sensores"
  ]);

  // Opciones disponibles para el dropdown
  const opcionesCampos = [
    { value: "ubicacion", label: "Ubicación" },
    { value: "historial_mantenimientos", label: "Historial de Mantenimientos" },
    { value: "ultima_acciones", label: "Últimas Acciones" },
    { value: "datos_sensores", label: "Datos de Sensores" }
  ];

  // Cargar activos al montar el componente
  useEffect(() => {
    const fetchActivos = async () => {
      try {
        const data = await gs.get("/obtener-activos") as ActivosResponse | ActivoReporte[];

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
    void fetchActivos();
  }, []);

  // Exportar PDF
  // Este método usa axios porque el gs tiene limitaciones con blobs (y otros tipos de respuesta)
  const handleExportarPDF = async () => {
  if (!activo) return;

  try {
    // Preparar el payload con los campos seleccionados
    const payload = {
      campos: camposSeleccionados
    };

    const response = await axios.post(
      `${BASE_URL}/gestion/reportes/activo/${activo}`,
      payload,
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
    setMensaje("No se pudo exportar el PDF");
  }
};

// <<< CAMBIO: función para vista previa PDF
const handleVistaPreviaPDF = async () => {
  if (!activo) return;
  try {
    // Preparar el payload con los campos seleccionados
    const payload = {
      campos: camposSeleccionados
    };

    const response = await axios.post(
      `${BASE_URL}/gestion/reportes/activo/${activo}`,
      payload,
      { responseType: "blob" }
    );
    const blob = new Blob([response.data], { type: "application/pdf" });
    const url = window.URL.createObjectURL(blob);
    setPdfUrl(url);
  } catch (error) {
    //console.error("Error al generar vista previa PDF:", error);
    setMensaje("No se pudo generar la vista previa");
  }
};



  return (
    <Box>
      {/* Filtros */}
      <Box sx={{ display: "flex", gap: 2, mb: 3, flexWrap: "wrap", alignItems: "flex-start" }}>
        <FormControl sx={{ minWidth: 200 }}>
          <InputLabel>Seleccionar Activo</InputLabel>
          <Select
            value={activo}
            onChange={(e) => { setActivo(e.target.value); }}
            label="Seleccionar Activo"
          >
            <MenuItem value="">Ninguno</MenuItem>
            {Array.isArray(activos) &&
              activos.map((a) => (
                <MenuItem key={a.id} value={a.id}>
                  {a.nombre}
                </MenuItem>
              ))}
          </Select>
        </FormControl>


        <Button variant="contained" onClick={handleVistaPreviaPDF} disabled={!activo}>
          Vista Previa PDF
        </Button>
        <Button variant="contained" onClick={handleExportarPDF} disabled={!activo}>
          Exportar PDF
        </Button>
      </Box>      {/* <<< CAMBIO: iframe para vista previa PDF */}
      {pdfUrl ? <Box sx={{ mt: 3, borderRadius: 2 }}>
          <iframe
            src={pdfUrl}
            style={{ width: "100%", height: "500px" }}
            title="Vista Previa PDF"
          />
        </Box> : null}
      <Box sx={{mt: 2}}>
        <FormControl sx={{ minWidth: 300 }}>
          <InputLabel>Campos del Reporte</InputLabel>
          <Select
            multiple
            value={camposSeleccionados}
            onChange={(e) => {
              const value = e.target.value;
              setCamposSeleccionados(typeof value === 'string' ? value.split(',') : value);
            }}
            input={<OutlinedInput label="Campos del Reporte" />}
            renderValue={(selected) => {
              return opcionesCampos
                .filter(opcion => selected.includes(opcion.value))
                .map(opcion => opcion.label)
                .join(', ');
            }}
          >
            {opcionesCampos.map((opcion) => (
              <MenuItem key={opcion.value} value={opcion.value}>
                <Checkbox checked={camposSeleccionados.includes(opcion.value)} />
                <ListItemText primary={opcion.label} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>

      </Box>
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
      <Snackbar
        open={Boolean(mensaje)}
        autoHideDuration={4000}
        onClose={() => { setMensaje(null); }}
      >
        <Alert severity="warning" onClose={() => { setMensaje(null); }}>
          {mensaje}
        </Alert>
      </Snackbar>    
    </Box>
  );
}

export default GenerarReporte;
