'use client'

import { useState } from 'react'
import { Box, Button, TextField, Typography, MenuItem, Collapse, IconButton } from '@mui/material'
import { DataGrid, GridColDef } from '@mui/x-data-grid'
import DeleteIcon from '@mui/icons-material/Delete'
import { esES } from '@mui/x-data-grid/locales';

interface Accion {
  accion: string
  fecha: string
  responsable: string
  observaciones: string
}

interface Lista {
  id: number
  activo: string
  tecnico: string
  fecha: string
  acciones: Accion[]
}

const activosMock = ['Ascensor', 'Caldera', 'Bomba de Agua']
const tecnicoActual = 'Técnico Actual'

const listasMock: Lista[] = [
  {
    id: 1,
    activo: 'Ascensor',
    tecnico: 'Técnico 1',
    fecha: '2025-08-01',
    acciones: [
      { accion: 'Revisar frenos', fecha: '2025-08-01', responsable: 'Técnico 1', observaciones: '' },
      { accion: 'Lubricar poleas', fecha: '2025-08-01', responsable: 'Técnico 1', observaciones: '' },
    ],
  },
  {
    id: 2,
    activo: 'Caldera',
    tecnico: 'Técnico 2',
    fecha: '2025-07-15',
    acciones: [
      { accion: 'Cambiar válvula', fecha: '2025-07-15', responsable: 'Técnico 2', observaciones: '' },
    ],
  },
]

export default function ListasAccionesView() {
  const [listas, setListas] = useState<Lista[]>(listasMock)
  const [activoSeleccionado, setActivoSeleccionado] = useState('')
  const [accionesNuevas, setAccionesNuevas] = useState<Accion[]>([{ accion: '', fecha: '', responsable: tecnicoActual, observaciones: '' }])
  const [expandedListaId, setExpandedListaId] = useState<number | null>(null)

  // Acciones temporales por lista existente
  const [accionesTemp, setAccionesTemp] = useState<{ [key: number]: Accion[] }>({})

  // --- Nueva lista ---
  const addAccionNueva = () => {
    setAccionesNuevas([...accionesNuevas, { accion: '', fecha: '', responsable: tecnicoActual, observaciones: '' }])
  }
  const handleAccionChange = (i: number, field: keyof Accion, value: string) => {
    const nuevas = [...accionesNuevas]; nuevas[i][field] = value; setAccionesNuevas(nuevas)
  }
  const handleEliminarAccionNueva = (i: number) => {
    const nuevas = accionesNuevas.filter((_, idx) => idx !== i)
    setAccionesNuevas(nuevas.length ? nuevas : [{ accion: '', fecha: '', responsable: tecnicoActual, observaciones: '' }])
  }
  const handleGuardarNuevaLista = () => {
    if (!activoSeleccionado) return alert('Selecciona un activo.')
    if (accionesNuevas.some(a => !a.accion || !a.fecha)) return alert('Completa todos los campos obligatorios')
    const nuevaLista: Lista = { id: listas.length + 1, activo: activoSeleccionado, tecnico: tecnicoActual, fecha: new Date().toISOString().split('T')[0], acciones: accionesNuevas }
    setListas([nuevaLista, ...listas])
    setAccionesNuevas([{ accion: '', fecha: '', responsable: tecnicoActual, observaciones: '' }])
    setActivoSeleccionado('')
  }

  // --- Acciones existentes ---
  const handleAccionExistenteChange = (listaId: number, i: number, field: keyof Accion, value: string) => {
    setListas(prev => prev.map(l => l.id === listaId ? { 
      ...l, acciones: l.acciones.map((a, idx) => idx === i ? { ...a, [field]: value } : a)
    } : l))
  }
  const handleEliminarAccionExistente = (listaId: number, i: number) => {
    setListas(prev => prev.map(l => l.id === listaId ? { 
      ...l, acciones: l.acciones.filter((_, idx) => idx !== i) 
    } : l))
  }

  // --- Acciones temporales ---
  const handleAddAccionTemp = (listaId: number) => {
    const prev = accionesTemp[listaId] || []
    setAccionesTemp({ ...accionesTemp, [listaId]: [...prev, { accion: '', fecha: new Date().toISOString().split('T')[0], responsable: tecnicoActual, observaciones: '' }] })
  }
  const handleAccionTempChange = (listaId: number, i: number, field: keyof Accion, value: string) => {
    const prev = accionesTemp[listaId] || []
    const nuevas = [...prev]; nuevas[i][field] = value
    setAccionesTemp({ ...accionesTemp, [listaId]: nuevas })
  }
  const handleEliminarAccionTemp = (listaId: number, i: number) => {
    const prev = accionesTemp[listaId] || []
    const nuevas = prev.filter((_, idx) => idx !== i)
    setAccionesTemp({ ...accionesTemp, [listaId]: nuevas })
  }
  const handleGuardarAccionesTemp = (listaId: number) => {
    const nuevas = accionesTemp[listaId] || []
    if (nuevas.some(a => !a.accion || !a.fecha)) return alert('Completa todas las acciones antes de guardar')
    setListas(prev => prev.map(l => l.id === listaId ? { ...l, acciones: [...l.acciones, ...nuevas] } : l))
    setAccionesTemp({ ...accionesTemp, [listaId]: [] })
  }

  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 70 },
    { field: 'activo', headerName: 'Activo', flex: 1 },
    { field: 'tecnico', headerName: 'Técnico', flex: 1 },
    { field: 'fecha', headerName: 'Fecha', flex: 1 },
    { field: 'acciones', headerName: 'N° Acciones', flex: 1},
    { field: 'ver', headerName: 'Acciones', flex: 1,
      renderCell: (params) => (
        <Button size="small" onClick={() => setExpandedListaId(expandedListaId === params.row.id ? null : params.row.id)}>
          {expandedListaId === params.row.id ? 'Ocultar' : 'Ver / Editar'}
        </Button>
      ),
    },
  ]

  return (
    <Box>
      <Typography variant="h6" gutterBottom>Crear Nueva Lista de Acciones</Typography>
      <Box sx={{ display:'flex', gap:2, mb:2 }}>
        <TextField label="Activo" select value={activoSeleccionado} onChange={e=>setActivoSeleccionado(e.target.value)} size="small" sx={{ minWidth:200 }}>
          {activosMock.map(a=><MenuItem key={a} value={a}>{a}</MenuItem>)}
        </TextField>
      </Box>
      <Box sx={{ mb:2}}>
        {accionesNuevas.map((a,i)=>(
          <Box key={i} sx={{ display:'grid', gridTemplateColumns:'1fr 1fr 1fr 1fr auto', gap:2, mb:2 }}>
            <TextField label="Acción" value={a.accion} onChange={e=>handleAccionChange(i,'accion',e.target.value)} fullWidth/>
            <TextField type="date" label="Fecha" InputLabelProps={{shrink:true}} value={a.fecha} onChange={e=>handleAccionChange(i,'fecha',e.target.value)} fullWidth/>
            <TextField label="Responsable" value={a.responsable} onChange={e=>handleAccionChange(i,'responsable',e.target.value)} fullWidth/>
            <TextField label="Observaciones" value={a.observaciones} onChange={e=>handleAccionChange(i,'observaciones',e.target.value)} fullWidth/>
            <IconButton color="error" onClick={()=>handleEliminarAccionNueva(i)}><DeleteIcon/></IconButton>
          </Box>
        ))}
        <Box sx={{display:'flex', gap:2}}>
          <Button onClick={addAccionNueva}>+ Acción</Button>
          <Button variant="contained" onClick={handleGuardarNuevaLista}>Guardar Lista</Button>
        </Box>
      </Box>

      <Box sx={{ mt:3 }}>
        <Typography variant="subtitle1" gutterBottom>Listas Anteriores</Typography>
        <DataGrid 
          autoHeight 
          rows={listas.map(l=>({...l, acciones:l.acciones.length}))} 
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
            paginationRowsPerPage: 'Listas por página',
            noRowsLabel: 'No hay listas asociadas',
            footerTotalRows: 'Total de listas:',
            footerTotalVisibleRows: (visibleCount, totalCount) =>
            `${visibleCount.toLocaleString()} de ${totalCount.toLocaleString()}`,
            footerRowSelected: (count) =>                          count > 1
              ? `${count.toLocaleString()} listas seleccionadas`
              : `${count.toLocaleString()} lista seleccionada`,
            paginationDisplayedRows: ({ from, to, count, estimated }) => {
            if (!estimated) {                            
              return `${from}–${to} de ${count !== -1 ? count : `más de ${to}`}`;
            }
            const estimatedLabel = estimated && estimated > to ? `alrededor de ${estimated}` : `más de ${to}`;
              return `${from}–${to} de ${count !== -1 ? count : estimatedLabel}`;
            },
          }}
        />
        {listas.map(l=>(
          <Collapse key={l.id} in={expandedListaId===l.id}>
            <Box sx={{p:2, mt:1}}>
              <Typography variant="subtitle2" gutterBottom>Acciones de {l.activo}</Typography>
              {l.acciones.map((a,i)=>(
                <Box key={i} sx={{ display:'grid', gridTemplateColumns:'1fr 1fr 1fr 1fr auto', gap:2, mb:2 }}>
                  <TextField label="Acción" value={a.accion} onChange={e=>handleAccionExistenteChange(l.id,i,'accion',e.target.value)} fullWidth/>
                  <TextField type="date" label="Fecha" InputLabelProps={{shrink:true}} value={a.fecha} onChange={e=>handleAccionExistenteChange(l.id,i,'fecha',e.target.value)} fullWidth/>
                  <TextField label="Responsable" value={a.responsable} onChange={e=>handleAccionExistenteChange(l.id,i,'responsable',e.target.value)} fullWidth/>
                  <TextField label="Observaciones" value={a.observaciones} onChange={e=>handleAccionExistenteChange(l.id,i,'observaciones',e.target.value)} fullWidth/>
                  <IconButton color="error" onClick={()=>handleEliminarAccionExistente(l.id,i)}><DeleteIcon/></IconButton>
                </Box>
              ))}
              {/* Acciones temporales */}
              {(accionesTemp[l.id] || []).map((a,i)=>(
                <Box key={i} sx={{ display:'grid', gridTemplateColumns:'1fr 1fr 1fr 1fr auto', gap:2, mb:2 }}>
                  <TextField label="Acción" value={a.accion} onChange={e=>handleAccionTempChange(l.id,i,'accion',e.target.value)} fullWidth/>
                  <TextField type="date" label="Fecha" InputLabelProps={{shrink:true}} value={a.fecha} onChange={e=>handleAccionTempChange(l.id,i,'fecha',e.target.value)} fullWidth/>
                  <TextField label="Responsable" value={a.responsable} onChange={e=>handleAccionTempChange(l.id,i,'responsable',e.target.value)} fullWidth/>
                  <TextField label="Observaciones" value={a.observaciones} onChange={e=>handleAccionTempChange(l.id,i,'observaciones',e.target.value)} fullWidth/>
                  <IconButton color="error" onClick={()=>handleEliminarAccionTemp(l.id,i)}><DeleteIcon/></IconButton>
                </Box>
              ))}
              <Box sx={{ display:'flex', gap:2 }}>
                <Button onClick={()=>handleAddAccionTemp(l.id)}>+ Acción</Button>
                <Button variant="contained" onClick={()=>handleGuardarAccionesTemp(l.id)}>Guardar Acciones Nuevas</Button>
              </Box>
            </Box>
          </Collapse>
        ))}
      </Box>
    </Box>
  )
}
