import { Box, Button, TextField, Typography, Select, MenuItem } from "@mui/material";
import { DataGrid } from '@mui/x-data-grid'
import { esES } from '@mui/x-data-grid/locales';

const DocumentosAsociados = () => {
  return (
    <Box>
      <Typography variant="h6" gutterBottom>Documentos Asociados</Typography>

      <Box sx={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 2, mb: 2 }}>
        <Button variant="outlined" component="label">
          Subir Archivo
          <input type="file" hidden />
        </Button>

        <Select defaultValue="" displayEmpty>
          <MenuItem value="">Categoría</MenuItem>
          <MenuItem value="ficha">Ficha técnica</MenuItem>
          <MenuItem value="mantencion">Informe de mantención</MenuItem>
        </Select>

        <Select defaultValue="" displayEmpty>
          <MenuItem value="">Seleccionar Activo</MenuItem>
          <MenuItem value="ascensor">Ascensor</MenuItem>
          <MenuItem value="caldera">Caldera</MenuItem>
          <MenuItem value="bomba de agua">Bomba de agua</MenuItem>
        </Select>
      </Box>

      <Button variant="contained">Guardar Documento</Button>

      <Box sx={{ mt: 3, border: "1px solid #ccc", borderRadius: 2, p: 2 }}>
        <Typography variant="subtitle1">Documentos Cargados</Typography>
        <ul>
          <li>Informe de Mantención - Ascensor - 01/07/2025</li>
          <li>Ficha Técnica - Caldera - 15/06/2025</li>
        </ul>
      </Box>
      <Box sx={{ mt:3 }}>
        <DataGrid
            autoHeight
            rows={[
                { id: 1, doc: 'Informe Mantención.pdf', categoria: 'Mantención', activo: 'Bomba', fecha: '2025-07-01' },
                { id: 2, doc: 'Ficha Técnica.docx', categoria: 'Ficha Técnica', activo: 'Caldera', fecha: '2025-06-10' },
                ]}
            columns={[
                { field: 'id', headerName: 'ID', width: 70 },
                { field: 'doc', headerName: 'Documento', flex: 1.5 },
                { field: 'categoria', headerName: 'Categoría', flex: 1 },
                { field: 'activo', headerName: 'Activo', flex: 1 },
                { field: 'fecha', headerName: 'Fecha', flex: 1 },
                ]}
            pageSizeOptions={[5]}
            disableRowSelectionOnClick
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
    </Box>
  );
};

export default DocumentosAsociados;
