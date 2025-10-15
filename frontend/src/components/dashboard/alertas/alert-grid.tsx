/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
'use client';

import * as React from 'react';
import { paths } from '@/paths';
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import { type SxProps, type Theme } from '@mui/material/styles';
import {
  DataGrid,
  type GridColDef,
  type GridFilterModel,
  type GridColumnVisibilityModel,
  type GridRowParams,
  type GridCellParams,
  type GridRenderCellParams,
} from '@mui/x-data-grid';
import { esES } from '@mui/x-data-grid/locales';

import { useUserToken } from '@/hooks/use-usertoken';
import { renderStatus, STATUS_OPTIONS } from './status';
import Services from '@/modules/Services';


interface Activo {
  id: number;
  edificio_id: number;
  creado_en: string;
  nombre: string;
  estado: 'OK' | 'Medio' | 'Crítico' | 'NN';
  tipo: string;
  ubicacion: string;
}

interface ApiActivosResponse {
  activos?: Activo[];
  total?: number;
  error?: { mensaje: string };
}

const gs = new Services();

const columns: GridColDef<Activo>[] = [
  { field: 'id', headerName: 'ID', width: 90 },
  {
    field: 'edificio_id',
    headerName: 'ID edificio',
    flex: 0.5,
    minWidth: 80,
    valueGetter: (params: GridCellParams<Activo>) => params.row.edificio_id,
  },
  {
    field: 'tipo',
    headerName: 'Tipo de Activo',
    flex: 1.2,
    minWidth: 120,
  },
  {
    field: 'ubicacion',
    headerName: 'Ubicación',
    flex: 1.5,
    minWidth: 160,
  },
  {
    field: 'estado',
    headerName: 'Estado',
    type: 'singleSelect',
    valueOptions: STATUS_OPTIONS,
    flex: 1,
    minWidth: 100,
    renderCell: (params: GridRenderCellParams<Activo>) => {
      return renderStatus(params as any);
    },
  },
  {
    field: 'nombre',
    headerName: 'Nombre',
    flex: 1.5,
    minWidth: 180,
  },
];

export default function DataGridDemo({
  sx,
  edificioSeleccionado,
}: {
  sx?: SxProps<Theme>;
  edificioSeleccionado?: string | null;
}) {
  const [filterModel, setFilterModel] = React.useState<GridFilterModel>({ items: [] });
  const { user } = useUserToken();
  const [activosList, setActivosList] = React.useState<Activo[]>([]);
  const [error, setError] = React.useState<string | null>(null);

  const getActivos = async (authToken: string) => {
    if (!authToken) {
      console.error("Token no disponible para getActivos.");
      return;
    }

    try {
      console.log(authToken)
      const response = await gs.authorizedGet('/obtener-todos-activos', authToken) as ApiActivosResponse;

      if (response.error) {
        setError(response.error.mensaje || 'Error al obtener activos.');
        setActivosList([]); //mocks en caso de
        return;
      }

      if (!response.activos || !Array.isArray(response.activos)) {
        setError('Respuesta inválida del servidor.');
        setActivosList([]);
        return;
      }

      const transformados: Activo[] = response.activos.map((item, idx) => ({
        id: typeof item.id === 'number' ? item.id : idx + 1,
        edificio_id: typeof item.edificio_id === 'number' ? item.edificio_id : Number(item.edificio_id) || 0,
        creado_en: item.creado_en || new Date().toISOString(),
        nombre: item.nombre || 'Activo sin nombre',
        estado: (item.estado || 'NN') as Activo['estado'],
        tipo: item.tipo || 'Desconocido',
        ubicacion: item.ubicacion || 'Ubicación no especificada',
      }));

      setActivosList(transformados);
      setError(null);
    } catch (err) {
      console.error('Error al obtener activos:', err);
      setError('Error de conexión. No se mostrarán activos.');
      setActivosList([]);
    }
  };

  React.useEffect(() => {
    const authToken = user?.token;
    if (authToken) {
      void getActivos(authToken);
    } else {
      console.log("Esperando token de usuario...");
      setActivosList([]); 
      setError("Autenticación pendiente o fallida.");
    }
  }, [user?.token]);

  React.useEffect(() => {
    if (edificioSeleccionado) {
      setFilterModel({
        items: [
          {
            field: 'edificio_id',
            operator: 'equals',
            value: edificioSeleccionado,
          },
        ],
      });
    } else {
      setFilterModel({ items: [] });
    }
  }, [edificioSeleccionado]);

  const [columnVisibilityModel, setColumnVisibilityModel] = React.useState<GridColumnVisibilityModel>({
    id: false,
    edificio_id: false,
  });

  return (
    <Card sx={sx}>
      <CardHeader title="Estado de Activos" subheader={error ? `⚠️ ${error}` : undefined} />
      <CardContent>
        <Box sx={{ height: 600, width: '100%' }}>
          <DataGrid
            columnVisibilityModel={columnVisibilityModel}
            onColumnVisibilityModelChange={(newModel) => setColumnVisibilityModel(newModel)}
            onRowClick={(params: GridRowParams<Activo>) => {
              window.location.href = paths.dashboard.activoDetail(params.row.id.toString());
            }}
            showToolbar
            rows={activosList}
            columns={columns}
            initialState={{
              pagination: { paginationModel: { pageSize: 9 } },
            }}
            pageSizeOptions={[9]}
            disableRowSelectionOnClick
            filterModel={filterModel}
            onFilterModelChange={(newModel) => setFilterModel(newModel)}
            localeText={{
              ...esES.components.MuiDataGrid.defaultProps.localeText,
              filterPanelInputLabel: 'Valor a filtrar',
              filterPanelOperator: 'Operador',
              filterPanelColumns: 'Filtrar por columna',
              toolbarColumns: 'Columnas visibles',
              toolbarFilters: 'Filtros',
              toolbarExport: 'Exportar',
              paginationRowsPerPage: 'Activos por página',
              noRowsLabel: 'No hay activos disponibles',
              footerTotalRows: 'Total de activos:',
              footerTotalVisibleRows: (visibleCount, totalCount) =>
                `${visibleCount.toLocaleString()} de ${totalCount.toLocaleString()}`,
              footerRowSelected: (count) =>
                count > 1
                  ? `${count.toLocaleString()} activos seleccionados`
                  : `${count.toLocaleString()} activo seleccionado`,
              paginationDisplayedRows: ({ from, to, count, estimated }) => {
                if (!estimated) {
                  return `${from}–${to} de ${count !== -1 ? count : `más de ${to}`}`;
                }
                const estimatedLabel = estimated && estimated > to ? `alrededor de ${estimated}` : `más de ${to}`;
                return `${from}–${to} de ${count !== -1 ? count : estimatedLabel}`;
              },
            }}
          />
        </Box>
      </CardContent>
    </Card>
  );
}