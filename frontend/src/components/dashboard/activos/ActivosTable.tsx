'use client'

import * as React from 'react'
import { DataGrid, GridColDef, GridRenderCellParams } from '@mui/x-data-grid'
import { Paper, Box, Chip, TextField, MenuItem, Stack} from '@mui/material'
import { esES } from '@mui/x-data-grid/locales'
import { Activo } from '@/types/'
import { activosMock } from '@/mocks/'

const activos = activosMock

const columns: GridColDef[] = [
  { field: 'id', headerName: 'ID', flex: 1, minWidth: 70 },
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
  const paginas = 8

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
            label="Nivel de riesgo"
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
          disableColumnFilter 
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
