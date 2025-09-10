'use client';

import * as React from 'react'
import { DataGrid, GridColDef, GridRenderCellParams } from '@mui/x-data-grid'
import { Paper, Box, Chip, TextField, MenuItem, Stack, Checkbox, Tooltip   } from '@mui/material'
import { esES } from '@mui/x-data-grid/locales'
import { Alerta } from '@/types/'
import { activosMock } from '@/mocks/'


const activos = activosMock

export function AlertasTable(): React.JSX.Element {
  const [alertas, setAlertas] = React.useState<Alerta[]>([])
  const [filtroTipo, setFiltroTipo] = React.useState('')
  const [filtroUbicacion, setFiltroUbicacion] = React.useState('')
  const [filtroEstado, setFiltroEstado] = React.useState('')
  const [filtroAtendidas, setFiltroAtendidas] = React.useState('')

  const handleAtender = (id: string) => {
    setAlertas((prev) =>
      prev.map((a) =>
        a.id === id ? { ...a, atendida: true } : a
      )
    )
  }

  const alertasFiltradas = alertas.filter((a) => {
    const tipoOK = !filtroTipo || a.tipoActivo.toLowerCase().includes(filtroTipo.toLowerCase())
    const ubicacionOK = !filtroUbicacion || a.ubicacion.toLowerCase().includes(filtroUbicacion.toLowerCase())
    const estadoOK = !filtroEstado || a.estado === filtroEstado
    const atencionOK =
      filtroAtendidas === '' ||
      (filtroAtendidas === 'pendientes' && !a.atendida) ||
      (filtroAtendidas === 'atendidas' && a.atendida)
    return tipoOK && ubicacionOK && estadoOK && atencionOK
  })

  const columns: GridColDef[] = [
    {
      field: 'atendida',
      headerName: 'Atendida',
      sortable: false,
      flex: 0.5,
      minWidth: 80,
      renderCell: (params: GridRenderCellParams) => (
        <Checkbox
          checked={params.value}
          disabled={params.value}
          onChange={() => handleAtender(params.row.id)}
          color="primary"
        />
      ),
    },
    { field: 'fecha', headerName: 'Fecha de alerta', flex: 1, minWidth: 140 },
    { field: 'id', headerName: 'ID', flex: 0.5, minWidth: 50 },
    { field: 'tipoActivo', headerName: 'Tipo de activo', flex: 1.5, minWidth: 120 },
    {
      field: 'estado',
      headerName: 'Nivel de riesgo',
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
    {
      field: 'descripcion',
      headerName: 'Descripción',
      flex: 2,
      minWidth: 220,
      renderCell: (params: GridRenderCellParams) => (
        <Tooltip title={params.value}>
          <span style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
            {params.value}
          </span>
        </Tooltip>
      ),
    },
    { field: 'ubicacion', headerName: 'Ubicación', flex: 1.5, minWidth: 160 },
  ]

  const estados = ['Medio', 'Crítico']
  const paginas = 8

  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <Paper elevation={3} sx={{ p: 2, mb: 2 }}>
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
          <TextField
            label="Estado de alerta"
            select
            value={filtroAtendidas}
            onChange={(e) => setFiltroAtendidas(e.target.value)}
            size="small"
            fullWidth
          >
            <MenuItem value="">Todas</MenuItem>
            <MenuItem value="pendientes">Pendientes</MenuItem>
            <MenuItem value="atendidas">Atendidas</MenuItem>
          </TextField>
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

      <Paper elevation={3} sx={{ p: 2 , height: '100%' }}>
        <DataGrid
          rows={alertasFiltradas}
          columns={columns}
          getRowId={(row) => row.id}
          initialState={{
            pagination: {
              paginationModel: { pageSize: paginas, page: 0 },
            },
          }}
          pageSizeOptions={Array.from({ length: Math.ceil(alertasFiltradas.length / paginas) }, (_, i) => (i + 1) * paginas)}
          disableRowSelectionOnClick
          disableColumnFilter
          localeText={{
            ...esES.components.MuiDataGrid.defaultProps.localeText,
            paginationRowsPerPage: 'Alertas por página',
            noRowsLabel: 'No hay alertas activas',
            footerTotalRows: 'Total de alertas:',
            footerTotalVisibleRows: (visibleCount, totalCount) =>
              `${visibleCount.toLocaleString()} de ${totalCount.toLocaleString()}`,
            footerRowSelected: (count) =>
              count > 1
                ? `${count.toLocaleString()} alertas seleccionadas`
                : `${count.toLocaleString()} alerta seleccionada`,
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
