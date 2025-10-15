/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
'use client'

import { useEffect, useState } from 'react'
import { Box, Button, TextField, Typography, MenuItem, Collapse, Snackbar, Alert } from '@mui/material'
import { DataGrid, type GridRenderCellParams, type GridColDef } from '@mui/x-data-grid'
import { esES } from '@mui/x-data-grid/locales'

import Services from '@/modules/Services'
import { useUserToken } from '@/hooks/use-usertoken';

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
  activo?: {
    id: number
    nombre: string
  }
  tecnico?: {
    id: number
    nombre: string
  }
}

interface Activo {
  id: number
  nombre: string
}

interface Tecnico {
  id: number
  nombre: string
  apellido: string
  especialidad: string
  email: string
  telefono: string
}

const tecnicoActualId = 1
const edificioId = 1

export default function ListasAccionesView() {
  const { user, isLoading} = useUserToken();
  const [acciones, setAcciones] = useState<Accion[]>([])
  const [activos, setActivos] = useState<Activo[]>([])
  const [tecnicos, setTecnicos] = useState<Tecnico[]>([])
  const [tecnicoSeleccionado, setTecnicoSeleccionado] = useState<number>(tecnicoActualId)
  const [expandedId, setExpandedId] = useState<number | null>(null)
  const [mensaje, setMensaje] = useState<string | null>(null);

  // Campos para nueva acción
  const [activoSeleccionado, setActivoSeleccionado] = useState('')
  const [descripcion, setDescripcion] = useState('')
  const [prioridad, setPrioridad] = useState('')
  const [tipo, setTipo] = useState('')

  // --- Cargar activos desde API ---
  useEffect(() => {
    // if (isLoading || !user) {return;}
    const fetchActivos = async () => {
      try {
        if (isLoading || !user) {
          console.log('user/token no disponibles', { isLoading, user })
          return
        }
        // console.log('user.token:', user.token)
        const data = await gs.authorizedGet("/obtener-todos-activos", user.token) as { activos: Activo[] }
        console.log("Activos cargados:", data)
        const activosArray = data.activos
        setActivos(Array.isArray(activosArray) ? activosArray : [])
        //console.log("Valor de setActivos (activos):", Array.isArray(activosArray) ? activosArray : [])
      } catch (error) {
        //console.error('Error al cargar activos:', error)
      }
    }
    void fetchActivos()
  }, [isLoading, user])

  // --- Cargar técnicos desde API ---
  useEffect(() => {
    const fetchTecnicos = async () => {
      try {
        if (isLoading || !user) {
          console.log('user/token no disponibles para técnicos', { isLoading, user })
          return
        }
        const data = await gs.authorizedGet(`/obtener-contactos-id/${edificioId}`, user.token) as { contactos: Tecnico[]; total: number }
        console.log("Técnicos cargados:", data)
        const tecnicosArray = data.contactos
        setTecnicos(Array.isArray(tecnicosArray) ? tecnicosArray : [])
        
        // Seleccionar el primer técnico de la lista por defecto
        if (Array.isArray(tecnicosArray) && tecnicosArray.length > 0) {
          setTecnicoSeleccionado(tecnicosArray[0].id)
        }
      } catch (error) {
        console.error('Error al cargar técnicos:', error)
      }
    }
    void fetchTecnicos()
  }, [isLoading, user])

  // --- API Acciones ---
  const fetchAcciones = async () => {
    try {
      if (!user) {
        console.log('Usuario no disponible para cargar acciones')
        return
      }

      const res = await gs.authorizedGet(`/obtener-acciones/${tecnicoSeleccionado}`, user.token) as Accion[]
      console.log("Acciones cargadas:", res)
      const accionesProcesadas = res.map((accion) => ({
        ...accion,
        tecnico_nombre: accion.tecnico?.nombre || 'Sin técnico',
        activo_nombre: accion.activo?.nombre || 'Sin activo'
      }))
      setAcciones(accionesProcesadas)
    } catch (err) {
      console.error('Error cargando acciones:', err)
    }
  }

  const crearAccion = async () => {
    if (!activoSeleccionado || !descripcion || !tipo || !prioridad) { 
      setMensaje('Completa todos los campos.'); 
      return; 
    }

    if (!user) {
      setMensaje('Usuario no autenticado');
      return;
    }

    try {
      const data = {
        titulo: `Mantenimiento ${tipo} - ${activos.find(a => a.id === Number(activoSeleccionado))?.nombre || 'Activo'}`,
        tecnico_id: tecnicoSeleccionado,
        tipo,
        descripcion,
        prioridad,
      }

      console.log("Datos para crear acción:", data)

      const res = await gs.authorizedPost(`/accion-mantenimiento-id/${activoSeleccionado}`, data, user.token) as { error?: boolean; mensaje?: string }

      if (!res.error) {
        // Acción creada correctamente
        setMensaje('Acción creada exitosamente')
        await fetchAcciones()
        setActivoSeleccionado('')
        setDescripcion('')
        setPrioridad('')
        setTipo('')
      } else {
        // Manejo de error
        setMensaje(`Error al crear la acción: ${res.mensaje || 'Error desconocido'}`)
      }
    } catch (err) {
      console.error('Error creando acción: ', err)
      setMensaje('Error inesperado al crear la acción')
    }
  }


  const actualizarEstado = async (id: number, nuevoEstado: string) => {
    try {
      if (!user) {
        setMensaje('Usuario no autenticado');
        return;
      }

      const res = await gs.authorizedPut(`/actualizar-estado-accion/${id}`, { estado: nuevoEstado }, user.token) as { error?: boolean; mensaje?: string };

      if (!res.error) {
        // éxito
        setMensaje('Estado actualizado exitosamente');
        await fetchAcciones()
      } else {
        // console.error('Error actualizando estado:', res.error)
        setMensaje(`Error al actualizar el estado: ${res.mensaje || 'Error desconocido'}`)
      }
    } catch (err) {
      console.error('Error actualizando estado:', err)
      setMensaje('Error inesperado al actualizar el estado')
    }
  }

  useEffect(() => {
    void fetchAcciones()
  }, [tecnicoSeleccionado])

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
      renderCell: (params: GridRenderCellParams<Accion>) => (
        <Button
          size="small"
          onClick={() => { setExpandedId(expandedId === params.row.id ? null : params.row.id); }}
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
          onChange={(e) => { setActivoSeleccionado(e.target.value); }}
          fullWidth
          sx={{ minWidth: 200 }}
        >
          {activos.map((a) => (
            <MenuItem key={a.id} value={a.id}>{a.nombre}</MenuItem>
          ))}
        </TextField>
        <TextField label="Descripción" value={descripcion} onChange={(e) => { setDescripcion(e.target.value); }} fullWidth />

        {/* Tipo como dropdown */}
        <TextField
          label="Tipo"
          select
          value={tipo}
          onChange={(e) => { setTipo(e.target.value); }}
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
          onChange={(e) => { setPrioridad(e.target.value); }}
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

      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="subtitle1">Acciones Asignadas</Typography>
        
        <TextField
          label="Filtrar por Técnico"
          select
          value={tecnicoSeleccionado}
          onChange={(e) => { setTecnicoSeleccionado(Number(e.target.value)); }}
          sx={{ minWidth: 250 }}
          size="small"
        >
          {tecnicos.map((t) => (
            <MenuItem key={t.id} value={t.id}>
              {t.nombre} {t.apellido} - {t.especialidad}
            </MenuItem>
          ))}
        </TextField>
      </Box>

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
  )
}
