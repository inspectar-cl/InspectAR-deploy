'use client';

import * as React from 'react';
import { DataGrid, type GridColDef, GridActionsCellItem, type GridRowId } from '@mui/x-data-grid';
import { Box, Alert, CircularProgress, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button, Tooltip } from '@mui/material';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { useUserToken } from '@/hooks/use-usertoken';
import { useForm, FormProvider } from 'react-hook-form';
import type { SolicitudFormData, EdificioData } from '@/types/formulario';
import { EdificioDataForm } from '@/components/dashboard/agregar-datos/edificio-datos-form';
import Services from '@/modules/Services';

const services = new Services();

// Asumo un tipo de dato simple
interface EdificioAPI {
  id: number;
  nombre: string;
  direccion: string;
  creado_en: string;
  latitud?: number;
  longitud?: number;
}

function transformarApiAForm(edificio: EdificioAPI): SolicitudFormData {
  return {
    tipoSolicitud: 'Edificio',
    asunto: '',
    detalles: '',
    datosEspecificos: {
      nombre: edificio.nombre,
      direccion: edificio.direccion,
      // Si la API no envía latitud, el formulario recibirá 0.
      // Si la API SÍ la envía, el formulario recibirá el valor real.
      latitud: edificio.latitud ?? 0,
      longitud: edificio.longitud ?? 0,
    },
  } as unknown as SolicitudFormData;
}

export function EdificiosDataGrid() {
  // // const router = useRouter();
  const { user } = useUserToken();
  const [rows, setRows] = React.useState<EdificioAPI[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  // Estados para el diálogo de confirmación de borrado
  const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
  const [selectedId, setSelectedId] = React.useState<GridRowId | null>(null);

  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [editingId, setEditingId] = React.useState<number | null>(null);
  const modalFormMethods = useForm<SolicitudFormData>();

  
  // --- 1. Carga de Datos ---
  React.useEffect(() => {
    if (!user?.token) return;

    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const data = await services.authorizedGet('/gestion/edificios', user.token) as { edificios: EdificioAPI[] };
        
        if (!data.edificios) {
           throw new Error("El formato de respuesta de la API es incorrecto.");
        }

        setRows(data.edificios);
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

    // Transforma los datos de la API al formato del formulario
    const datosFormulario = transformarApiAForm(itemToEdit);

    // Guarda el ID y pre-llena el formulario del modal
    setEditingId(idNum);
    modalFormMethods.reset(datosFormulario);
    // Abre el modal
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingId(null);
    modalFormMethods.reset(); // Limpia el formulario
  };

  // Función de submit para el formulario del modal
  const onModalSubmit = async (data: SolicitudFormData) => {
    if (!editingId || !user?.token) return;

    const datosParaApi = data.datosEspecificos as EdificioData;

    modalFormMethods.clearErrors();

    try {
      await services.authorizedPut(`/actualizar-edificio/${editingId}`, datosParaApi, user.token);

      // Éxito: actualiza la fila en la tabla localmente
      setRows((prevRows) =>
        prevRows.map((row) =>
          row.id === editingId
            ? { ...row, ...datosParaApi } // Combina datos antiguos y nuevos
            : row
        )
      );
      handleCloseModal(); // Cierra el modal
      // alert('Edificio actualizado con éxito.');

    } catch (err) {
      // Muestra el error dentro del modal
      modalFormMethods.setError('root', { 
        type: 'manual', 
        message: err instanceof Error ? err.message : 'Error al guardar' 
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
      await services.authorizedDelete(`/eliminar-edificio/${selectedId}`, user.token);

      setRows((prevRows) => prevRows.filter((row) => row.id !== selectedId));
      // alert('Edificio eliminado con éxito.');

    } catch (err) {
      setError(err instanceof Error ? err.message : 'No se pudo eliminar el elemento.');
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
        <DialogTitle>Editar Edificio (ID: {editingId})</DialogTitle>
        <FormProvider {...modalFormMethods}>
          <form onSubmit={modalFormMethods.handleSubmit(onModalSubmit)}>
            <DialogContent>
              {/* Aquí se renderiza tu formulario de edificio existente */}
              <EdificioDataForm />
            </DialogContent>
            <DialogActions>
              <Button onClick={handleCloseModal}>Cancelar</Button>
              <Button type="submit" variant="contained">
                Guardar Cambios
              </Button>
            </DialogActions>
          </form>
        </FormProvider>
      </Dialog>

      {/* --- Diálogo de Confirmación de Borrado --- */}
      <Dialog open={openDeleteDialog} onClose={() => { setOpenDeleteDialog(false); }}>
        <DialogTitle>Confirmar Eliminación</DialogTitle>
        <DialogContent>
          <DialogContentText>
            ¿Estás seguro de que deseas eliminar este elemento? Esta acción no se puede deshacer.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => { setOpenDeleteDialog(false); }}>Cancelar</Button>
          {/* eslint-disable-next-line jsx-a11y/no-autofocus -- Delete confirmation requires focus for UX */}
          <Button onClick={confirmDelete} color="error" autoFocus>
            Eliminar
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}