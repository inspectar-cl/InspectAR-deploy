'use client'

import * as React from 'react'
import { DataGrid, GridColDef, GridRenderCellParams } from '@mui/x-data-grid'
import { Paper, Box, Chip, TextField, MenuItem, Stack} from '@mui/material'
import { esES } from '@mui/x-data-grid/locales'

type Activo = {
  id: string
  tipoActivo: string
  estado: 'OK' | 'Medio' | 'Crítico'
  descripcion: string
  ubicacion: string
}

const activos: Activo[] = [
  { id: 'A1', tipoActivo: 'Ascensor #1', estado: 'OK', descripcion: 'Funciona correctamente', ubicacion: 'Edificio A, Santiago' },
  { id: 'A2', tipoActivo: 'Ascensor #2', estado: 'OK', descripcion: 'Funciona correctamente', ubicacion: 'Edificio A, Santiago' },
  { id: 'A3', tipoActivo: 'Ascensor #3', estado: 'Medio', descripcion: 'Mantenimiento programado no realizado', ubicacion: 'Edificio A, Santiago' },
  { id: 'A4', tipoActivo: 'Ascensor #4', estado: 'Medio', descripcion: 'Puerta atascada, requiere revisión', ubicacion: 'Edificio B, Santiago' },
  { id: 'C1', tipoActivo: 'Caldera #1', estado: 'Crítico', descripcion: 'Temperatura fuera de rango, riesgo de daño', ubicacion: 'Edificio A, Santiago' },
  { id: 'C2', tipoActivo: 'Caldera #2', estado: 'Crítico', descripcion: 'Anomalía detectada – revisar urgentemente', ubicacion: 'Edificio B, Santiago' },
  { id: 'C3', tipoActivo: 'Caldera #3', estado: 'Crítico', descripcion: 'Fuga de gas detectada', ubicacion: 'Edificio C, Santiago' },
  { id: 'C4', tipoActivo: 'Caldera #4', estado: 'Crítico', descripcion: 'Presión excesiva, riesgo de explosión', ubicacion: 'Edificio C, Santiago' },
  { id: 'B1', tipoActivo: 'Bomba de Agua #1', estado: 'Crítico', descripcion: 'Motor sobrecalentado, riesgo de falla', ubicacion: 'Edificio A, Santiago' },
  { id: 'B2', tipoActivo: 'Bomba de Agua #2', estado: 'Crítico', descripcion: 'Fuga de presión, riesgo de parada', ubicacion: 'Edificio A, Santiago' },
  { id: 'B3', tipoActivo: 'Bomba de Agua #3', estado: 'Crítico', descripcion: 'Riesgo crítico de falla en 7 días', ubicacion: 'Edificio B, Santiago' },
  { id: 'B4', tipoActivo: 'Bomba de Agua #4', estado: 'Crítico', descripcion: 'Nivel de agua crítico', ubicacion: 'Edificio B, Santiago' },
  { id: 'E1', tipoActivo: 'Sistema Eléctrico #1', estado: 'Medio', descripcion: 'Sobrecarga detectada, monitorear consumo', ubicacion: 'Edificio A, Santiago' },
  { id: 'E2', tipoActivo: 'Sistema Eléctrico #2', estado: 'Medio', descripcion: 'Pico de voltaje registrado', ubicacion: 'Edificio B, Santiago' },
  { id: 'E3', tipoActivo: 'Sistema Eléctrico #3', estado: 'Medio', descripcion: 'Variación de frecuencia, revisar panel', ubicacion: 'Edificio C, Santiago' },
]

const columns: GridColDef[] = [
  { field: 'id', headerName: 'ID', flex: 1, minWidth: 90 },
  { field: 'tipoActivo', headerName: 'Tipo de activo', flex: 1, minWidth: 160 },
  {
    field: 'estado',
    headerName: 'Estado',
    flex: 1,
    minWidth: 100,
    renderCell: (params) => {
      const color =
        params.value === 'Crítico'
          ? 'error'
          : params.value === 'Medio'
            ? 'warning'
            : 'success'
      return <Chip label={params.value} color={color} size="small" />
    },
  },
  { field: 'descripcion', headerName: 'Descripción', flex: 2, minWidth: 200 },
  { field: 'ubicacion', headerName: 'Ubicación', flex: 1.5, minWidth: 160 },
]

export default function ActivosTable() {
  const [filtroEstado, setFiltroEstado] = React.useState('')
  const [filtroTipo, setFiltroTipo] = React.useState('')
  const [filtroUbicacion, setFiltroUbicacion] = React.useState('')

  const activosFiltrados = activos.filter((a) => {
    const estadoOK = !filtroEstado || a.estado === filtroEstado
    const tipoOK = !filtroTipo || a.tipoActivo.toLowerCase().includes(filtroTipo.toLowerCase())
    const ubicacionOK = !filtroUbicacion || a.ubicacion.toLowerCase().includes(filtroUbicacion.toLowerCase())
    return estadoOK && tipoOK && ubicacionOK
  })

  const tiposActivos = Array.from(new Set(activos.map((a) => a.tipoActivo)))
  const estados = ['OK', 'Medio', 'Crítico']
  const paginas = 10

  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <Paper elevation={3} sx={{ p: 2, mb: 2 }}>
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
          <TextField
            label="Tipo de activo"
            value={filtroTipo}
            onChange={(e) => setFiltroTipo(e.target.value)}
            size="small"
            fullWidth
          />
          <TextField
            label="Estado"
            select
            value={filtroEstado}
            onChange={(e) => setFiltroEstado(e.target.value)}
            size="small"
            fullWidth
          >
            <MenuItem value="">Todos</MenuItem>
            {estados.map((estado) => (
              <MenuItem key={estado} value={estado}>
                {estado}
              </MenuItem>
            ))}
          </TextField>
          <TextField
            label="Buscar ubicación"
            value={filtroUbicacion}
            onChange={(e) => setFiltroUbicacion(e.target.value)}
            size="small"
            fullWidth
          />
        </Stack>
      </Paper>

      <Paper elevation={3} sx={{ p: 2, height: '100%' }}>
        <DataGrid
          rows={activosFiltrados}
          columns={columns}
          getRowId={(row) => row.id}
          initialState={{
            pagination: {
              paginationModel: { pageSize: paginas, page: 0 },
            },
          }}
          pageSizeOptions={Array.from({ length: Math.ceil(activosFiltrados.length / paginas) }, (_, i) => (i + 1) * paginas)}
          disableRowSelectionOnClick
          localeText={{
            ...esES.components.MuiDataGrid.defaultProps.localeText,
            paginationRowsPerPage: 'Activos por página',
            noRowsLabel: 'No hay activos disponibles',
            footerTotalRows: 'Total de activos:',
            footerTotalVisibleRows: (visibleCount, totalCount) =>
              `${visibleCount.toLocaleString()} de ${totalCount.toLocaleString()}`,
            footerRowSelected: (count) =>
              count > 1
                ? `${count.toLocaleString()} activos seleccionados`
                : `${count.toLocaleString()} activo seleccionado`,
            paginationDisplayedRows: ({ from, to, count, estimated }) => {
              if (!estimated) {
                return `${from}–${to} de ${count !== -1 ? count : `más de ${to}`}`;
              }
              const estimatedLabel = estimated && estimated > to ? `alrededor de ${estimated}` : `más de ${to}`;
              return `${from}–${to} de ${count !== -1 ? count : estimatedLabel}`;
            },
          }}
        />
      </Paper>
    </Box>
  )
}
