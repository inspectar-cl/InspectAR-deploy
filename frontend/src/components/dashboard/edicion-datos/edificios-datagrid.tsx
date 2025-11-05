'use client';

import * as React from 'react';
import { DataGrid, GridColDef, GridActionsCellItem, GridRowId } from '@mui/x-data-grid';
import { Box, Alert, CircularProgress, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button, Tooltip } from '@mui/material';
import { useRouter } from 'next/navigation';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { paths } from '@/paths';
import { useUserToken } from '@/hooks/use-usertoken';

// Asumo un tipo de dato simple
interface EdificioAPI {
  id: number;
  nombre: string;
  direccion: string;
  latitud: number;
  longitud: number;
}

export function EdificiosDataGrid() {
  const router = useRouter();
  const { user } = useUserToken();
  const [rows, setRows] = React.useState<EdificioAPI[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  // Estados para el diálogo de confirmación de borrado
  const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
  const [selectedId, setSelectedId] = React.useState<GridRowId | null>(null);

  // --- 1. Carga de Datos ---
  React.useEffect(() => {
    if (!user?.token) return; // Esperar token
    
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        // !!! Reemplaza con tu endpoint de API real !!!
        const response = await fetch('/api/edificios/listar', {
            headers: { 'Authorization': `Bearer ${user.token}` }
        });
        if (!response.ok) throw new Error('Error al cargar los edificios');
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
    // Aqui se redirige a pag de formulario existente (cambiar si no es bueno)
    // pasándole un 'id' para que sepa que está en modo "Editar".
    const url = `${paths.dashboard.agregarDatos('edificio')}?id=${id}`;
    router.push(url);
  };

  const handleDelete = (id: GridRowId) => {
    setSelectedId(id);
    setOpenDeleteDialog(true);
  };

  const confirmDelete = async () => {
    if (!selectedId || !user?.token) return;
    try {
      // Api para borrar edificio por ID
      // await fetch(`/api/edificios/${selectedId}`, { 
      //    method: 'DELETE', 
      //    headers: { 'Authorization': `Bearer ${user.token}` }
      // });
      console.log(`Simulando borrado de ID: ${selectedId}`);
      setRows((prevRows) => prevRows.filter((row) => row.id !== selectedId));
    } catch (err) {
      setError('No se pudo eliminar el elemento.');
    } finally {
      setOpenDeleteDialog(false);
      setSelectedId(null);
    }
  };

  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 90 },
    { field: 'nombre', headerName: 'Nombre', flex: 1, minWidth: 200 },
    { field: 'direccion', headerName: 'Dirección', flex: 2, minWidth: 300 },
    { field: 'latitud', headerName: 'Latitud', width: 150 },
    { field: 'longitud', headerName: 'Longitud', width: 150 },
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

      {/* --- Diálogo de Confirmación de Borrado --- */}
      <Dialog open={openDeleteDialog} onClose={() => setOpenDeleteDialog(false)}>
        <DialogTitle>Confirmar Eliminación</DialogTitle>
        <DialogContent>
          <DialogContentText>
            ¿Estás seguro de que deseas eliminar este elemento? Esta acción no se puede deshacer.
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