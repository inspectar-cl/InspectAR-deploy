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
      // const res = await gs.post('/gestion/acciones', {
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify({
      //     activo_id: Number(activoSeleccionado), // usamos el id del activo
      //     tecnico_id: tecnicoActualId,
      //     tipo,
      //     descripcion,
      //     prioridad,
      //   }),
      // })
      if (res.ok) {
        await fetchAcciones()
        setActivoSeleccionado('')
        setDescripcion('')
        setPrioridad('')
        setTipo('')
      }
    } catch (err) {
      console.error('Error creando acción:', err)
    }
  }

  const actualizarEstado = async (id: number, nuevoEstado: string) => {
    try {
      await gs.put(`/gestion/acciones/${id}/estado`, {
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ estado: nuevoEstado }),
      })
      await fetchAcciones()
    } catch (err) {
      console.error('Error actualizando estado:', err)
    }
  }

  useEffect(() => {
    fetchAcciones()
  }, [])

  // --- Columnas ---
  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 70 },
    { field: 'tecnico_nombre', headerName: 'Técnico', width: 100},
    { field: 'descripcion', headerName: 'Descripción', flex: 2 },
    { field: 'estado', headerName: 'Estado', flex: 1 },
    { field: 'prioridad', headerName: 'Prioridad', width: 100 },
    {
      field: 'acciones',
      headerName: 'Acciones',
      flex: 1,
      renderCell: (params) => (
        <Button
          size="small"
          onClick={() => setExpandedId(expandedId === params.row.id ? null : params.row.id)}
        >
          {expandedId === params.row.id ? 'Ocultar' : 'Ver / Cambiar Estado'}
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
          size="small"
          sx={{ minWidth: 200 }}
        >
          {activos.map((a) => (
            <MenuItem key={a.id} value={a.id}>{a.nombre}</MenuItem>
          ))}
        </TextField>
        <TextField label="Descripción" value={descripcion} onChange={(e) => setDescripcion(e.target.value)} fullWidth />
        <TextField label="Tipo" value={tipo} onChange={(e) => setTipo(e.target.value)} size="small" />
        <TextField label="Prioridad" value={prioridad} onChange={(e) => setPrioridad(e.target.value)} size="small" />
        <Button variant="contained" onClick={crearAccion}>Guardar</Button>
      </Box>

      <Typography variant="subtitle1" gutterBottom>Acciones Asignadas</Typography>
      <DataGrid
        autoHeight
        rows={acciones}
        columns={columns}
        pageSizeOptions={[5]}
        disableRowSelectionOnClick
        localeText={esES.components.MuiDataGrid.defaultProps.localeText}
      />

      {acciones.map((a) => (
        <Collapse key={a.id} in={expandedId === a.id}>
          <Box sx={{ p: 2, mt: 1, border: '1px solid #ddd', borderRadius: 2 }}>
            <Typography variant="subtitle2">Detalle de Acción</Typography>
            <Typography>Técnico: {a.tecnico.nombre}</Typography>
            <Typography>Descripción: {a.descripcion}</Typography>
            <Typography>Estado: {a.estado}</Typography>
            <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
              <Button size="small" variant="outlined" onClick={() => actualizarEstado(a.id, 'en_progreso')}>En progreso</Button>
              <Button size="small" variant="outlined" onClick={() => actualizarEstado(a.id, 'completado')}>Completado</Button>
            </Box>
          </Box>
        </Collapse>
      ))}
    </Box>
  )
}
