'use client';

import * as React from 'react';
import { DataGrid, GridColDef, GridActionsCellItem, GridRowId } from '@mui/x-data-grid';
import { Box, Alert, CircularProgress, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button, Tooltip } from '@mui/material';
import { useRouter } from 'next/navigation';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { paths } from '@/paths';
import type { TipoActivo } from '@/types/formulario'; // Asumo que importas tus tipos
import { useUserToken } from '@/hooks/use-usertoken';

// Asumo un tipo de dato que viene de la API
interface ActivoAPI {
  id: number;
  tipoActivo: TipoActivo;
  edificioId: number;
  edificioNombre?: string;
  ubicacion: string;
  descripcion?: string;
}

export function ActivosDataGrid() {
  const router = useRouter();
  const { user } = useUserToken();
  const [rows, setRows] = React.useState<ActivoAPI[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
  const [selectedId, setSelectedId] = React.useState<GridRowId | null>(null);

  // Carga de Datos
  React.useEffect(() => {
    if (!user?.token) return;
    
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        // Aqui debe hacerse llamado de la api, ojala con nombre de edificio tambien
        const response = await fetch('/api/activos/listar', {
            headers: { 'Authorization': `Bearer ${user.token}` }
        });
        if (!response.ok) throw new Error('Error al cargar los activos');
        
        // Simulación de datos si la API no trae el nombre
        const data: ActivoAPI[] = (await response.json()).map((activo: ActivoAPI) => ({
            ...activo,
            edificioNombre: activo.edificioNombre || `Edificio ID ${activo.edificioId}` // Fallback
        }));
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
    const url = `${paths.dashboard.agregarDatos('activo')}?id=${id}`;
    router.push(url);
  };

  const handleDelete = (id: GridRowId) => {
    setSelectedId(id);
    setOpenDeleteDialog(true);
  };

  const confirmDelete = async () => {
    if (!selectedId || !user?.token) return;
    try {
      // Aqui api para el borrado de activos
      // await fetch(`/api/activos/${selectedId}`, { method: 'DELETE', ... });
      console.log(`Simulando borrado de ID: ${selectedId}`);
      setRows((prevRows) => prevRows.filter((row) => row.id !== selectedId));
    } catch (err) {
      setError('No se pudo eliminar el elemento.');
    } finally {
      setOpenDeleteDialog(false);
      setSelectedId(null);
    }
  };

  // Definición de Columnas
  const columns: GridColDef[] = [
    { field: 'id', headerName: 'ID', width: 90 },
    { field: 'tipoActivo', headerName: 'Tipo', flex: 1, minWidth: 150 },
    { field: 'ubicacion', headerName: 'Ubicación', flex: 2, minWidth: 200 },
    { 
      field: 'edificioNombre', 
      headerName: 'Edificio', 
      flex: 1, 
      minWidth: 150,
      // Aqui se usa el id del edificio si solo poseemos eso
      // valueGetter: (params) => `Edificio ID ${params.row.edificioId}`, 
    },
    { field: 'descripcion', headerName: 'Descripción', flex: 2, minWidth: 250 },
    {
      field: 'actions',
      type: 'actions',
      headerName: 'Acciones',
      width: 100,
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
      {/* Diálogo de Confirmación */}
      <Dialog open={openDeleteDialog} onClose={() => setOpenDeleteDialog(false)}>
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
      </Dialog>
    </Box>
  );
}