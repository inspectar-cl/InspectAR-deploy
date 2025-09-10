/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
import { useState, useEffect } from "react";
import { Box, Button, TextField, Typography, Select, MenuItem, Snackbar, Alert } from "@mui/material";
import { DataGrid, type GridColDef } from '@mui/x-data-grid'
import { esES } from '@mui/x-data-grid/locales';
import Services from '@/modules/Services';

const gs = new Services();

function DocumentosAsociados() {
  const [documentos, setDocumentos] = useState<any[]>([]);
  const [activos, setActivos] = useState<any[]>([]);
  const [activo, setActivo] = useState("");
  const [categorias, setCategorias] = useState<string[]>([]);
  const [categoriaSeleccionada, setCategoriaSeleccionada] = useState("");
  const [archivo, setArchivo] = useState<File | null>(null);
  const [nombreDocumento, setNombreDocumento] = useState("");
  const [descripcionDocumento, setDescripcionDocumento] = useState("");
  const [palabrasClave, setPalabrasClave] = useState("");
  const [mensaje, setMensaje] = useState<string | null>(null);

  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', flex:0.5, filterable: false},
    { field: 'nombre', headerName: 'Nombre', flex: 1.5, filterable: false},
    { field: 'categoria', headerName: 'Categoría', flex: 1, filterable: true },
    { field: 'activo_nombre', headerName: 'Activo', flex: 1, filterable: false },
    { field: 'palabras_clave', headerName: 'Palabras clave', flex: 1, filterable: true },
    { field: 'creado_en_formateado', headerName: 'Fecha', flex: 1, filterable: true },
  ];

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

  // Cargar categorías al montar el componente
  useEffect(() => {
    const fetchCategorias = async () => {
      try {
        const data = await gs.get("/documentacion/documentos");
        //console.log("Respuesta documentos:", data);

        const documentos = data.documentos || [];

        // Extraer categorías únicas
        const categoriasUnicas = Array.from(
          new Set(documentos.map((doc: any) => doc.categoria))
        );

        setCategorias(categoriasUnicas.filter(c => typeof c === "string"));
      } catch (error) {
        //console.error("Error al conectar con API Gateway:", error);
        setCategorias([]);
      }
    };

    fetchCategorias();
  }, []);

  // Cargar documentos
  useEffect(() => {
    const fetchDocumentos = async () => {
      try {
        const data = await gs.get("/documentacion/documentos");

        const docsArray = data?.documentos || [];
        const documentosConFecha = docsArray.map((doc: any) => {
          const fecha = new Date(doc.fecha_emision); // 🔹 usa tu campo real
          const fechaFormateada = new Intl.DateTimeFormat('es-CL').format(fecha); // dd/mm/aaaa
          return {
            ...doc,
            creado_en_formateado: fechaFormateada,
            palabras_clave: doc.palabras_clave || "-",
            activo_nombre: activos.find((a: any) => a.id === doc.activo_id)?.nombre || `ID ${doc.activo_id}`,
          };
        });

        setDocumentos(documentosConFecha);

        // Extraer categorías únicas
        const cats = Array.from(new Set(docsArray.map((d: any) => d.categoria)));
        setCategorias(cats.filter(c => typeof c === "string"));
        } catch (error) {
          //console.error("Error al cargar documentos:", error);
        }
    };
    if (activos.length > 0) fetchDocumentos(); // 🔹 Solo corre cuando ya hay activos
}, [activos]); // 🔹 Dependencia en activos

  const handleGuardarDocumento = async () => {
    if (!archivo) { setMensaje("Debe seleccionar un archivo"); return; }
    if (!categoriaSeleccionada) { setMensaje("Debe seleccionar una categoría"); return; }
    if (!activo) { setMensaje("Debe seleccionar un activo"); return; }
    if (!nombreDocumento || !descripcionDocumento || !palabrasClave) {
      setMensaje("Debe completar todos los campos"); return;
    }
    // const formData = {
    //   archivo: archivo,
    //   activo_id: activo,
    //   categoria: categoriaSeleccionada,
    //   nombre: nombreDocumento,
    //   descripcion: descripcionDocumento,
    //   palabras_clave: palabrasClave,
    // }
    const formData = new FormData();
    formData.append("archivo", archivo);
    formData.append("activo_id", activo);
    formData.append("categoria", categoriaSeleccionada);
    formData.append("nombre", nombreDocumento);
    formData.append("descripcion", descripcionDocumento);
    formData.append("palabras_clave", palabrasClave);
    try {
      await gs.post("/documentacion/documentos", formData);
      // Subido correctamente, puedes recargar documentos
      const data = await gs.get("/documentacion/documentos");
      const docsArray = data?.documentos || [];
      const documentosConFecha = docsArray.map((doc: any) => {
        const fecha = new Date(doc.fecha_emision); 
        const fechaFormateada = new Intl.DateTimeFormat('es-CL').format(fecha); // dd/mm/aaaa
        return {
          ...doc,
          creado_en_formateado: fechaFormateada,
          palabras_clave: doc.palabras_clave || "-",
          activo_nombre: activos.find((a: any) => a.id === doc.activo_id)?.nombre || `ID ${doc.activo_id}`,
        };
      });
      setDocumentos(documentosConFecha);

      // Limpiar formulario
      setArchivo(null);
      setNombreDocumento("");
      setDescripcionDocumento("");
      setPalabrasClave("");
      setActivo("");
      setCategoriaSeleccionada("");
    } catch (error) {
      setMensaje("Error al subir documento");
    }
  };


  return (
    <Box>
      <Typography variant="h6" gutterBottom>Documentos Asociados</Typography>

      <Box sx={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 2, mb: 2 }}>
        {/* Selector de archivo */}
        <Button variant="outlined" component="label">
          Subir Archivo
          <input
            type="file"
            hidden
            onChange={(e) => { setArchivo(e.target.files ? e.target.files[0] : null); }}
          />
        </Button>

        {/* Selector de categoría */}
        <Select
          value={categoriaSeleccionada}
          onChange={(e) => { setCategoriaSeleccionada(e.target.value); }}
          displayEmpty
        >
          <MenuItem value="">Categoría</MenuItem>
          {categorias.map((cat) => (
            <MenuItem key={cat} value={cat}>
              {cat.replace(/_/g, " ")} {/* para mostrarlo bonito */}
            </MenuItem>
          ))}
        </Select>

        {/* Selector de activo */}
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
      </Box>

      {/* Inputs adicionales para subir documento */}
      <Box sx={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 2, mb: 2 }}>
        <TextField
          label="Nombre del Documento"
          value={nombreDocumento}
          onChange={(e) => { setNombreDocumento(e.target.value); }}
          fullWidth
        />
        <TextField
          label="Palabras Clave (separadas por coma)"
          value={palabrasClave}
          onChange={(e) => { setPalabrasClave(e.target.value); }}
          fullWidth
        />
        <TextField
          label="Descripción"
          value={descripcionDocumento}
          onChange={(e) => { setDescripcionDocumento(e.target.value); }}
          fullWidth
        />
      </Box>
      {/* Botón para guardar documento */}
      <Button variant="contained" onClick={handleGuardarDocumento}>
        Guardar Documento
        <input
          type="file"
          hidden
          onChange={(e) => { setArchivo(e.target.files ? e.target.files[0] : null); }}
        />
      </Button>

      <Box sx={{ mt:3 }}>
        <DataGrid
            rows={documentos}
            columns={columns}
            pageSizeOptions={[5]}
            disableRowSelectionOnClick
            showToolbar
            slotProps={{
              toolbar: {
                showQuickFilter: true,      // 🔹 CAMBIO: activa barra de búsqueda
                quickFilterProps: { debounceMs: 300 },
              },
            }}
            localeText={{
                ...esES.components.MuiDataGrid.defaultProps.localeText,
                  filterPanelInputLabel: 'Valor a filtrar',
                  filterPanelOperator: 'Operador',
                  filterPanelColumns: 'Filtrar por columna',
                  toolbarColumns: 'Columnas visibles',
                  toolbarFilters: 'Filtros',
                  toolbarExport: 'Exportar',
                  // Traducción de paginación
                  paginationRowsPerPage: 'Documentos por página',
                  noRowsLabel: 'No hay documentos asociados',
                  footerTotalRows: 'Total de documentos:',
                  footerTotalVisibleRows: (visibleCount, totalCount) =>
                  `${visibleCount.toLocaleString()} de ${totalCount.toLocaleString()}`,
                  footerRowSelected: (count) =>                          count > 1
                    ? `${count.toLocaleString()} documentos seleccionados`
                    : `${count.toLocaleString()} documento seleccionado`,
                  paginationDisplayedRows: ({ from, to, count, estimated }) => {
                    if (!estimated) {                            
                      return `${from}–${to} de ${count !== -1 ? count : `más de ${to}`}`;
                    }
                    const estimatedLabel = estimated && estimated > to ? `alrededor de ${estimated}` : `más de ${to}`;
                    return `${from}–${to} de ${count !== -1 ? count : estimatedLabel}`;
                    },
            }}
        />
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

export default DocumentosAsociados;
