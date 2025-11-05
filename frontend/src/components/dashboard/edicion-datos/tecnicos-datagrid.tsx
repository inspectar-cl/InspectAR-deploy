'use client';

import * as React from 'react';
import {
  DataGrid,
  GridColDef,
  GridActionsCellItem,
  GridRowId
} from '@mui/x-data-grid';
import {
  Box,
  Alert,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
  Tooltip,
  Chip,
} from '@mui/material';
import { useRouter } from 'next/navigation';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { paths } from '@/paths';
import type { EspecialidadTecnico } from '@/types/formulario';
import { useUserToken } from '@/hooks/use-usertoken';

// Interfaz para los datos que vienen de la API
interface TecnicoAPI {
  id: number;
  nombre: string;
  correo: string;
  telefono: string;
  especialidad: EspecialidadTecnico;
  activosAsociados: number[]; // Array de IDs
}

export function TecnicosDataGrid() {
  const router = useRouter();
  const { user } = useUserToken();
  const [rows, setRows] = React.useState<TecnicoAPI[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
  const [selectedId, setSelectedId] = React.useState<GridRowId | null>(null);

  React.useEffect(() => {
    if (!user?.token) return;

    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        // !!! Reemplaza con tu endpoint de API real !!!
        const response = await fetch('/api/tecnicos/listar', {
          headers: { Authorization: `Bearer ${user.token}` },
        });
        if (!response.ok) throw new Error('Error al cargar los técnicos');
        const data = await response.json();
        setRows(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error desconocido');
      } finally {
        setIsLoading(false);
      }
    };
    fetchData();
  }, [user?.token]);

  const handleEdit = (id: GridRowId) => {
    // Redirige al formulario de "agregar-datos" en modo edición
    const url = `${paths.dashboard.agregarDatos('tecnico')}?id=${id}`;
    router.push(url);
  };

  const handleDelete = (id: GridRowId) => {
    setSelectedId(id);
    setOpenDeleteDialog(true);
  };

  const confirmDelete = async () => {
    if (!selectedId || !user?.token) return;
    try {
      // api para borrar tecnicos por id
      // await fetch(`/api/tecnicos/${selectedId}`, { method: 'DELETE', ... });
      console.log(`Simulando borrado de ID: ${selectedId}`);
      
      // Actualiza la UI
      setRows((prevRows) => prevRows.filter((row) => row.id !== selectedId));
    } catch (err) {
      setError('No se pudo eliminar el elemento.');
    } finally {
      setOpenDeleteDialog(false);
      setSelectedId(null);
    }
  };

  const columns: GridColDef<TecnicoAPI>[] = [
    { field: 'id', headerName: 'ID', width: 90 },
    { field: 'nombre', headerName: 'Nombre', flex: 1, minWidth: 200 },
    {
      field: 'especialidad',
      headerName: 'Especialidad',
      width: 150,
      renderCell: (params) => (
        <Chip label={params.value} color="info" size="small" />
      ),
    },
    { field: 'correo', headerName: 'Correo', flex: 1, minWidth: 200 },
    { field: 'telefono', headerName: 'Teléfono', flex: 1, minWidth: 150 },
    {
      field: 'activosAsociados',
      headerName: 'Activos Asignados',
      width: 150,
      type: 'number',
      align: 'left',
      headerAlign: 'left',
      
      valueGetter: (params: { row: TecnicoAPI }) =>
        params.row.activosAsociados.length,

      renderCell: (params: { row: TecnicoAPI; value?: number }) => (
        <Tooltip title={params.row.activosAsociados.join(', ')}>
          <span>{params.value ?? 0} activos</span>
        </Tooltip>
      ),
    },
    {
      field: 'actions',
      type: 'actions',
      headerName: 'Acciones',
      width: 100,
      cellClassName: 'actions',
      getActions: ({ id }) => [
        <Tooltip title="Editar" key="edit">
          <GridActionsCellItem
            icon={<EditIcon />}
            label="Editar"
            onClick={() => handleEdit(id)}
            color="primary"
          />
        </Tooltip>,
        <Tooltip title="Eliminar" key="delete">
          <GridActionsCellItem
            icon={<DeleteIcon />}
            label="Eliminar"
            onClick={() => handleDelete(id)}
            color="inherit"
          />
        </Tooltip>,
      ],
    },
  ];

  if (isLoading) return <CircularProgress />;
  if (error) return <Alert severity="error">{error}</Alert>;

  return (
    <Box sx={{ height: 600, width: '100%' }}>
      <DataGrid
        rows={rows}
        columns={columns}
        initialState={{
          pagination: { paginationModel: { pageSize: 10 } },
        }}
        pageSizeOptions={[10, 25, 50]}
        showToolbar
        slotProps={{
          toolbar: {
            showQuickFilter: true,
            quickFilterProps: { debounceMs: 500 },
          },
        }}
        
        disableRowSelectionOnClick
      />

      {/* --- Diálogo de Confirmación --- */}
      <Dialog open={openDeleteDialog} onClose={() => setOpenDeleteDialog(false)}>
        <DialogTitle>Confirmar Eliminación</DialogTitle>
        <DialogContent>
          <DialogContentText>
            ¿Estás seguro de que deseas eliminar este elemento? Esta acción no se
            puede deshacer.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDeleteDialog(false)}>Cancelar</Button>
          <Button onClick={confirmDelete} color="error" autoFocus>
            Eliminar
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}