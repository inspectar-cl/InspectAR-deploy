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
import type { SolicitudFormData, TecnicoData, EspecialidadTecnico } from '@/types/edicion-data';
import { useForm, FormProvider } from 'react-hook-form';
import { useUserToken } from '@/hooks/use-usertoken';
import { TecnicoEditForm } from './tecnico-edit-form';

// Interfaz para los datos que vienen de la API
interface TecnicoAPI {
  id: number;
  nombre: string;
  correo: string;
  telefono: string;
  especialidad: EspecialidadTecnico;
  activosAsociados: number[]; // Array de IDs
}

const getEstadoChipColor = (
  estado: 'Medio' | 'OK' | 'Crítico' | string
) => {
  if (estado === 'Crítico') return 'error';
  if (estado === 'Medio') return 'warning';
  if (estado === 'OK') return 'success';
  return 'default';
};

function transformarApiAForm(tecnico: TecnicoAPI): SolicitudFormData {
  return {
    tipoSolicitud: 'Técnico', // Relleno
    asunto: '',
    detalles: '',
    // Solo mapea los campos editables
    datosEspecificos: {
      nombre: tecnico.nombre,
      correo: tecnico.correo,
      telefono: tecnico.telefono,
      especialidad: tecnico.especialidad,
    },
  } as unknown as SolicitudFormData;
}

export function TecnicosDataGrid() {
  const router = useRouter();
  const { user } = useUserToken();
  const [rows, setRows] = React.useState<TecnicoAPI[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
  const [selectedId, setSelectedId] = React.useState<GridRowId | null>(null);

  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [editingId, setEditingId] = React.useState<number | null>(null);
  const modalFormMethods = useForm<SolicitudFormData>();

  React.useEffect(() => {
    if (!user?.token) return;
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await fetch('/api/gestion/tecnicos', {
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
    void fetchData();
  }, [user?.token]);

  const handleEdit = (id: GridRowId) => {
    const idNum = Number(id);
    const itemToEdit = rows.find((row) => row.id === idNum);
    if (!itemToEdit) return;
    const datosFormulario = transformarApiAForm(itemToEdit);
    setEditingId(idNum);
    modalFormMethods.reset(datosFormulario);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingId(null);
    modalFormMethods.reset();
  };

  const onModalSubmit = async (data: SolicitudFormData) => {
    if (!editingId || !user?.token) return;

    // Aseguramos a TS que 'datosEspecificos' tiene la forma de TecnicoData
    const datosDelFormulario = data.datosEspecificos as TecnicoData;

    // Creamos el payload para la API (solo 4 campos)
    const payload = {
      nombre: datosDelFormulario.nombre,
      correo: datosDelFormulario.correo,
      telefono: datosDelFormulario.telefono,
      // --- ¡AQUÍ ESTÁ LA CORRECCIÓN! ---
      // Hacemos el cast del 'string' genérico al tipo 'EspecialidadTecnico'
      especialidad: datosDelFormulario.especialidad as EspecialidadTecnico,
    };
    
    console.log('Enviando actualización para ID:', editingId, payload);
    modalFormMethods.clearErrors();

    try {
      // Asumo que la ruta es /api/editar-tecnico/{id}
      const response = await fetch(`/api/editar-tecnico/${editingId}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${user.token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => null);
        throw new Error(errorData?.error || `Error del servidor: ${response.status}`);
      }

      alert('Técnico actualizado con éxito.');

      // Actualiza la fila en la tabla localmente
      setRows((prevRows) =>
        prevRows.map((row) =>
          row.id === editingId
            ? { ...row, ...payload } // Ahora 'payload' tiene el tipo correcto
            : row
        )
      );

      handleCloseModal();
    } catch (err) {
      console.error('Error al actualizar:', err);
      modalFormMethods.setError('root', {
        type: 'manual',
        message: err instanceof Error ? err.message : 'Error al guardar',
      });
    }
  };

  const handleDelete = (id: GridRowId) => {
    setSelectedId(id);
    setOpenDeleteDialog(true);
  };

  const confirmDelete = async () => {
    if (!selectedId || !user?.token) return;
    setError(null);

    try {
      // !!! Asumo esta ruta de API !!!
      const response = await fetch(`/api/eliminar-tecnico/${selectedId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${user.token}`,
        },
      });

       if (!response.ok) {
        const errorData = await response.json().catch(() => null);
        throw new Error(errorData?.error || `Error del servidor: ${response.status}`);
      }
      
      alert('Técnico eliminado con éxito.');
      setRows((prevRows) => prevRows.filter((row) => row.id !== selectedId));
    } catch (err) {
      console.error('Error al eliminar:', err);
      setError(err instanceof Error ? err.message : 'No se pudo eliminar el elemento.');
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
            onClick={() => {handleEdit(id)}}
            color="primary"
          />
        </Tooltip>,
        <Tooltip title="Eliminar" key="delete">
          <GridActionsCellItem
            icon={<DeleteIcon />}
            label="Eliminar"
            onClick={() => {handleDelete(id)}}
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

      <Dialog open={isModalOpen} onClose={handleCloseModal} maxWidth="sm" fullWidth>
        <DialogTitle>Editar Técnico (ID: {editingId})</DialogTitle>
        <FormProvider {...modalFormMethods}>
          <form onSubmit={modalFormMethods.handleSubmit(onModalSubmit)}>
            <DialogContent>
              {modalFormMethods.formState.errors.root && (
                <Alert severity="error" sx={{ mb: 2 }}>
                  {modalFormMethods.formState.errors.root.message}
                </Alert>
              )}
              {/* ¡Aquí usamos tu NUEVO formulario de edición simple! */}
              <TecnicoEditForm />
            </DialogContent>
            <DialogActions>
              <Button onClick={handleCloseModal}>Cancelar</Button>
              <Button
                type="submit"
                variant="contained"
                disabled={modalFormMethods.formState.isSubmitting}
              >
                {modalFormMethods.formState.isSubmitting ? (
                  <CircularProgress size={24} color="inherit" />
                ) : (
                  'Guardar Cambios'
                )}
              </Button>
            </DialogActions>
          </form>
        </FormProvider>
      </Dialog>

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