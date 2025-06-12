'use client'

import * as React from 'react'
import { DataGrid, GridColDef, GridRenderCellParams } from '@mui/x-data-grid'
import { Paper, Box, Chip, TextField, MenuItem, Stack, Checkbox, Tooltip   } from '@mui/material'
import { esES } from '@mui/x-data-grid/locales'

type Activo = {
  id: string
  tipoActivo: string
  estado: 'OK' | 'Medio' | 'Crítico'
  descripcion: string
  ubicacion: string
}

type Alerta = {
  id: string
  tipoActivo: string
  descripcion: string
  estado: 'Medio' | 'Crítico'
  ubicacion: string
  fecha: string
  atendida: boolean
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

const generarAlertas = (activos: Activo[]): Alerta[] =>
  activos
    .filter((a) => a.estado !== 'OK')
    .map((a) => ({
      id: a.id,
      tipoActivo: a.tipoActivo,
      descripcion: a.descripcion,
      estado: a.estado,
      ubicacion: a.ubicacion,
      fecha: generarFechaProxima(),
      atendida: false,
    }))

// Fecha estimada de ocurrencia: aleatoria entre 1 y 10 días desde hoy
const generarFechaProxima = (): string => {
  const hoy = new Date()
  hoy.setDate(hoy.getDate() + Math.floor(Math.random() * 10) + 1)
  return hoy.toLocaleDateString('es-CL')
}

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

export default function AlertasTable() {
  const [alertas, setAlertas] = React.useState<Alerta[]>(generarAlertas(activos))
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
            <MenuItem value="todas">Todas</MenuItem>
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
