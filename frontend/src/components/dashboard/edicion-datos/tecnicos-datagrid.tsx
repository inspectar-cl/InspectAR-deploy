'use client';

import * as React from 'react';
import {
  DataGrid,
  GridColDef,
  GridActionsCellItem,
  GridRowId,
  GridRenderCellParams,
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
import type { SolicitudFormData, TecnicoData, EspecialidadTecnico } from '@/types/edicion-data';
import { useForm, FormProvider } from 'react-hook-form';
import { useUserToken } from '@/hooks/use-usertoken';
import { TecnicoEditForm } from './tecnico-edit-form';

// Interfaz para los datos que vienen de la API
interface TecnicoAPI {
  id: number;
  nombre: string;
  apellido: string;
  email: string;
  telefono: string;
  especialidad: string;
  autorizado: boolean;
  empresa_id: number;
  creado_en: string;
  // 'activosAsociados' no viene en la respuesta de la API, se omite
}

function transformarApiAForm(tecnico: TecnicoAPI): SolicitudFormData {
  const nombreCompleto = `${tecnico.nombre} ${tecnico.apellido}`;
  
  let especialidadForm: EspecialidadTecnico;
  if (tecnico.especialidad === "Electricidad Industrial") {
    especialidadForm = "Eléctrico";
  } else if (tecnico.especialidad === "Sistemas HVAC") {
    especialidadForm = "Climatización";
  } else {
    especialidadForm = "Mecánico"; // Ejemplo de fallback
  }

  return {
    tipoSolicitud: 'Técnico',
    asunto: '',
    detalles: '',
    datosEspecificos: {
      nombre: nombreCompleto,
      correo: tecnico.email,
      telefono: tecnico.telefono,
      especialidad: especialidadForm, 
    },
  } as unknown as SolicitudFormData;
}

export function TecnicosDataGrid() {
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
        const response = await fetch('/api/gestion/tecnicos', { // Ruta de tu API
          headers: { Authorization: `Bearer ${user.token}` },
        });
        if (!response.ok) throw new Error('Error al cargar los técnicos');
        
        const data = await response.json();
        
        if (!Array.isArray(data)) {
           throw new Error("El formato de respuesta de la API es incorrecto, se esperaba un array.");
        }
        
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

    const datosDelFormulario = data.datosEspecificos as TecnicoData;

    // Asumo que tu API de edición espera este payload
    const payload = {
      nombre: datosDelFormulario.nombre,
      email: datosDelFormulario.correo,
      telefono: datosDelFormulario.telefono,
      especialidad: datosDelFormulario.especialidad,
    };
    
    console.log('Enviando actualización para ID:', editingId, payload);
    modalFormMethods.clearErrors();

    try {
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
            ? { ...row, ...payload, email: payload.email } // Asegúrate de actualizar los campos correctos
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
      const response = await fetch(`/api/eliminar-tecnico/${selectedId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${user.token}` },
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
    {
      field: 'especialidad',
      headerName: 'Especialidad',
      width: 180,
      renderCell: (params: GridRenderCellParams<TecnicoAPI>) => (
        <Chip label={params.value} color="info" size="small" />
      ),
    },
    { field: 'email', headerName: 'Email', flex: 1, minWidth: 200 }, // Usa 'email'
    { field: 'telefono', headerName: 'Teléfono', flex: 1, minWidth: 150 },
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