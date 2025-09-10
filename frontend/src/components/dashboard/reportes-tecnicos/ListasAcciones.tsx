'use client'

import { useEffect, useState } from 'react'
import { Box, Button, TextField, Typography, MenuItem, Collapse } from '@mui/material'
import { DataGrid, GridColDef } from '@mui/x-data-grid'
import { esES } from '@mui/x-data-grid/locales'

import Services from '@/modules/Services'
const gs = new Services()

interface Accion {
  id: number
  activo_id: number
  tecnico_id: number
  tipo: string
  descripcion: string
  estado: string
  prioridad: string
  fecha_inicio: string
  creado_en: string
  activo?: string
  tecnico?: {
    id: number
    nombre: string
  }
}

interface Activo {
  id: number
  nombre: string
}

const tecnicoActualId = 1

export default function ListasAccionesView() {
  const [acciones, setAcciones] = useState<Accion[]>([])
  const [activos, setActivos] = useState<Activo[]>([])
  const [expandedId, setExpandedId] = useState<number | null>(null)

  // Campos para nueva acción
  const [activoSeleccionado, setActivoSeleccionado] = useState('')
  const [descripcion, setDescripcion] = useState('')
  const [prioridad, setPrioridad] = useState('')
  const [tipo, setTipo] = useState('')

  // --- Cargar activos desde API ---
  useEffect(() => {
    const fetchActivos = async () => {
      try {
        const data = await gs.get("/obtener-activos")
        console.log("Activos cargados:", data)
        const activosArray = data.activos
        setActivos(Array.isArray(activosArray) ? activosArray : [])
        console.log("Valor de setActivos (activos):", Array.isArray(activosArray) ? activosArray : [])
      } catch (error) {
        console.error('Error al cargar activos:', error)
      }
    }
    fetchActivos()
  }, [])

  // --- API Acciones ---
  const fetchAcciones = async () => {
    try {
      const res = await await gs.get(`/gestion/acciones/tecnico/${tecnicoActualId}`)
      console.log("Acciones cargadas:", res)
      const accionesProcesadas = res.map((accion: any) => ({
        ...accion,
        tecnico_nombre: accion.tecnico?.nombre || 'Sin técnico',
        activo_nombre: accion.activo?.nombre || 'Sin activo'
      }))
      // const data = await res.json()
      setAcciones(accionesProcesadas)
    } catch (err) {
      console.error('Error cargando acciones:', err)
    }
  }

  const crearAccion = async () => {
    if (!activoSeleccionado || !descripcion) return alert('Completa todos los campos.')

    try {
      const data = {
        titulo: "prueba",
        activo_id: Number(activoSeleccionado),
        tecnico_id: tecnicoActualId,
        tipo,
        descripcion,
        prioridad,
      }

      console.log("Datos para crear acción:", data)
      
      const res = await gs.post('/gestion/acciones', data)

      if (!res.error) {
        // Acción creada correctamente
        await fetchAcciones()
        setActivoSeleccionado('')
        setDescripcion('')
        setPrioridad('')
        setTipo('')
      } else {
        // Manejo de error
        console.error('Error en creación:', res.error)
        alert('Error al crear la acción: ' + (res.mensaje || 'Error desconocido'))
      }
    } catch (err) {
      console.error('Error creando acción:', err)
      alert('Error inesperado al crear la acción')
    }
  }


  const actualizarEstado = async (id: number, nuevoEstado: string) => {
    try {
      const res = await gs.put(`/gestion/acciones/${id}/estado`, { estado: nuevoEstado })

      if (!res.error) {
        // éxito
        await fetchAcciones()
      } else {
        console.error('Error actualizando estado:', res.error)
        alert('Error al actualizar el estado: ' + (res.mensaje || 'Error desconocido'))
      }
    } catch (err) {
      console.error('Error actualizando estado:', err)
      alert('Error inesperado al actualizar el estado')
    }
  }

  useEffect(() => {
    fetchAcciones()
  }, [])

  // --- Columnas ---
  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 70, flex:0.3},
    { field: 'tecnico_nombre', headerName: 'Técnico', width: 100, flex:0.5},
    { field: 'descripcion', headerName: 'Descripción', flex: 2 },
    { field: 'tipo', headerName: 'Tipo', flex: 0.7 , minWidth: 70},
    { field: 'estado', headerName: 'Estado', flex: 0.7, minWidth: 70},
    { field: 'prioridad', headerName: 'Prioridad', minWidth: 50 , flex:0.5},
    {
      field: 'acciones',
      headerName: 'Acciones',
      flex: 1,
      renderCell: (params) => (
        <Button
          size="small"
          onClick={() => setExpandedId(expandedId === params.row.id ? null : params.row.id)}
        >
          {expandedId === params.row.id ? 'Ocultar' : 'Cambiar Estado'}
        </Button>
      ),
    },
  ]

  return (
    <Box>
      <Typography variant="h6" gutterBottom>Crear Nueva Acción</Typography>

      <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
        <TextField
          label="Activo"
          select
          value={activoSeleccionado}
          onChange={(e) => setActivoSeleccionado(e.target.value)}
          fullWidth
          sx={{ minWidth: 200 }}
        >
          {activos.map((a) => (
            <MenuItem key={a.id} value={a.id}>{a.nombre}</MenuItem>
          ))}
        </TextField>
        <TextField label="Descripción" value={descripcion} onChange={(e) => setDescripcion(e.target.value)} fullWidth />

        {/* Tipo como dropdown */}
        <TextField
          label="Tipo"
          select
          value={tipo}
          onChange={(e) => setTipo(e.target.value)}
          fullWidth
          sx={{ minWidth: 150 }}
        >
          <MenuItem value="">Seleccione tipo</MenuItem>
          <MenuItem value="preventivo">Preventivo</MenuItem>
          <MenuItem value="correctivo">Correctivo</MenuItem>
        </TextField>

        {/* Prioridad como dropdown */}
        <TextField
          label="Prioridad"
          select
          value={prioridad}
          onChange={(e) => setPrioridad(e.target.value)}
          fullWidth
          sx={{ minWidth: 100 }}
        >
          <MenuItem value="">Seleccione prioridad</MenuItem>
          <MenuItem value="alta">Alta</MenuItem>
          <MenuItem value="media">Media</MenuItem>
          <MenuItem value="baja">Baja</MenuItem>
        </TextField>

        <Button variant="contained" onClick={crearAccion}>Guardar</Button>
      </Box>

      <Typography variant="subtitle1" gutterBottom>Acciones Asignadas</Typography>

      {acciones.map((a) => (
        <Collapse key={a.id} in={expandedId === a.id}>
          <Box sx={{ p: 2, mt: 2, mb: 2, border: '1px solid #ddd', borderRadius: 2 }}>
            <Typography variant="subtitle2">Detalle de Acción</Typography>
            <Typography>Técnico: {a.tecnico?.nombre ?? 'Sin técnico'}</Typography>
            <Typography>Descripción: {a.descripcion}</Typography>
            <Typography>Tipo: {a.tipo}</Typography>
            <Typography>Estado: {a.estado}</Typography>
            <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
              <Button size="small" variant="outlined" onClick={() => actualizarEstado(a.id, 'en_progreso')}>En progreso</Button>
              <Button size="small" variant="outlined" onClick={() => actualizarEstado(a.id, 'completado')}>Completado</Button>
            </Box>
          </Box>
        </Collapse>
      ))}

      <DataGrid
        rows={acciones}
        columns={columns}
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
                  paginationRowsPerPage: 'Acciones por página',
                  noRowsLabel: 'No hay acciones asociadas',
                  footerTotalRows: 'Total de acciones:',
                  footerTotalVisibleRows: (visibleCount, totalCount) =>
                  `${visibleCount.toLocaleString()} de ${totalCount.toLocaleString()}`,
                  footerRowSelected: (count) =>                          count > 1
                    ? `${count.toLocaleString()} acciones seleccionadas`
                    : `${count.toLocaleString()} acción seleccionada`,
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
  )
}
