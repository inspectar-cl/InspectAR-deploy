'use client';

import * as React from 'react';
import { DataGrid, GridColDef, GridActionsCellItem, GridRowId, GridRenderCellParams } from '@mui/x-data-grid';
import { Box, Alert, CircularProgress, Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button, Tooltip, Chip } from '@mui/material';
import { useRouter } from 'next/navigation';
import { Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { paths } from '@/paths';
//import type { TipoActivo } from '@/types/formulario'; // Asumo que importas tus tipos
import { useUserToken } from '@/hooks/use-usertoken';
import type { SolicitudFormData } from '@/types/formulario';
import { useForm, FormProvider } from 'react-hook-form';
import { ActivoEditForm } from './activo-edit-form';

// Asumo un tipo de dato que viene de la API
interface SensorAnidado {
  sensor_id: string;
  tipo: string;
  estado: string;
}
interface ActivoAPI {
  id: number;
  nombre: string;
  tipo: string;
  estado: string;
  ubicacion: string;
  edificio_id: number;
  creado_en: string;
  sensores: SensorAnidado[];
  descripcion?: string; // Asumo que puede venir
}

const getEstadoChipColor = (
  estado: 'Medio' | 'OK' | 'Crítico' | string
) => {
  if (estado === 'Crítico') return 'error';
  if (estado === 'Medio') return 'warning';
  if (estado === 'OK') return 'success';
  return 'default';
};

function transformarApiAForm(activo: ActivoAPI): SolicitudFormData {
  return {
    tipoSolicitud: 'Activo',
    asunto: '',
    detalles: '',
    datosEspecificos: {
      nombre: activo.nombre,
      ubicacion: activo.ubicacion,
      descripcion: activo.descripcion || '',
    },
  } as unknown as SolicitudFormData;
}

export function ActivosDataGrid() {
  const router = useRouter();
  const { user } = useUserToken();
  const [rows, setRows] = React.useState<ActivoAPI[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
  const [selectedId, setSelectedId] = React.useState<GridRowId | null>(null);

  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [editingId, setEditingId] = React.useState<number | null>(null);
  const modalFormMethods = useForm<SolicitudFormData>();

  // Carga de Datos
  React.useEffect(() => {
    if (!user?.token) return;

    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await fetch(
          '/api/obtener-todos-activos?sensores=true', // URL correcta
          {
            headers: { Authorization: `Bearer ${user.token}` },
          }
        );
        if (!response.ok) throw new Error('Error al cargar los activos');

        // La API devuelve un objeto { activos: [...] }, extraemos el array
        const data: { activos: ActivoAPI[] } = await response.json();

        if (!data.activos) {
          throw new Error("La respuesta de la API no tiene el formato esperado.");
        }
        
        setRows(data.activos); // <-- Asignamos el array 'activos'
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error desconocido');
      } finally {
        setIsLoading(false);
      }
    };
    fetchData();
  }, [user?.token]);

  const handleEdit = (id: GridRowId) => {
    const idNum = Number(id);
    const itemToEdit = rows.find((row) => row.id === idNum);
    if (!itemToEdit) return;

    // Transformamos los datos de la API al formato del formulario
    const datosFormulario = transformarApiAForm(itemToEdit);

    // Guardamos el ID que estamos editando
    setEditingId(idNum);
    // Pre-llenamos el formulario del modal con los datos
    modalFormMethods.reset(datosFormulario);
    // Abrimos el modal
    setIsModalOpen(true);
  };

  const handleDelete = (id: GridRowId) => {
    setSelectedId(id);
    setOpenDeleteDialog(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingId(null);
    modalFormMethods.reset(); // Limpia el formulario al cerrar
  };

  const onModalSubmit = async (data: SolicitudFormData) => {
    if (!editingId || !user?.token) return;

    console.log('Enviando actualización para ID:', editingId, data.datosEspecificos);
    try {
      const response = await fetch(`/api/editar-activo/${editingId}`, {
        method: 'PUT',
        headers: { 
          'Authorization': `Bearer ${user.token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data.datosEspecificos),
      });
      if (!response.ok) throw new Error('Falló la actualización');
      
      // Simulación de éxito
      alert('Activo actualizado');

      // Actualizar la fila en la tabla localmente (para no recargar todo)
      setRows((prevRows) =>
        prevRows.map((row) =>
          row.id === editingId
            ? { ...row, ...(data.datosEspecificos as any) } // Actualiza la fila
            : row
        )
      );
      
      handleCloseModal(); // Cierra el modal
    } catch (err) {
      console.error('Error al actualizar:', err);
      alert('Error al actualizar el activo.');
    }
  };

  const confirmDelete = async () => {
    if (!selectedId || !user?.token) return;
    try {
      // Aqui api para el borrado de activos
      await fetch(`/api/eliminar-activo/${selectedId}`, 
      { method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${user.token}`,
        },
      });
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
  const columns: GridColDef<ActivoAPI>[] = [
    { field: 'id', headerName: 'ID', width: 90 },
    { field: 'nombre', headerName: 'Nombre', flex: 1, minWidth: 200 },
    { field: 'tipo', headerName: 'Tipo', flex: 1, minWidth: 150 },
    {
      field: 'estado',
      headerName: 'Estado',
      width: 130,
      renderCell: (params: GridRenderCellParams<ActivoAPI>) => (
        <Chip
          label={params.value}
          color={getEstadoChipColor(params.value)}
          size="small"
          variant="outlined"
        />
      ),
    },
    { field: 'ubicacion', headerName: 'Ubicación', flex: 2, minWidth: 200 },
    { field: 'edificio_id', headerName: 'ID Edificio', width: 120 },
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
        
        // --- 6. Sintaxis moderna para la Toolbar (Buscador/Filtros) ---
        
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
        <DialogTitle>Editar Activo (ID: {editingId})</DialogTitle>
        <FormProvider {...modalFormMethods}>
          <form onSubmit={modalFormMethods.handleSubmit(onModalSubmit)}>
            <DialogContent>
              {/* ¡Aquí re-usamos tu formulario! */}
              <ActivoEditForm />
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
      
      {/* --- 7. Diálogo de Confirmación (Corregido, sin anidar) --- */}
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