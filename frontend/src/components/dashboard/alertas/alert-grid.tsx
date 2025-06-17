'use client';

import * as React from 'react';
import { paths } from '@/paths';
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import { DataGrid, GridColDef, GridFilterModel, GridColumnVisibilityModel} from '@mui/x-data-grid';
import { esES } from '@mui/x-data-grid/locales';

import { activosMock, alertasMock} from '@/mocks/'

//Import de estatus personalizado
import {
  renderStatus,
  STATUS_OPTIONS,
} from './status';

const activos = activosMock;

const columns: GridColDef<(typeof activos)[number]>[] = [
  { field: 'id', headerName: 'ID', width: 90 },
  {
    field: 'id_edificio',
    headerName: 'ID edificio',
    width: 90,
  },
  {
    field: 'tipoActivo',
    headerName: 'Tipo de Activo',
    width: 150,
  },
  {
    field: 'ubicacion',
    headerName: 'Ubicación',
    width: 160,
  },
  {
    field: 'estado',
    renderCell: renderStatus,
    headerName: 'Estado',
    type: 'singleSelect',
    valueOptions: STATUS_OPTIONS,
    width: 110,
  },
  {
    field: 'descripcion',
    headerName: 'Descripción',
    description: 'Descripcion de estado del activo',
    width: 320,
  },
];

export default function DataGridDemo({sx, edificioSeleccionado,}: {sx?: any; edificioSeleccionado?: string | null;}) {
  const [filterModel, setFilterModel] = React.useState<GridFilterModel>({
    items: [],
  });

  React.useEffect(() => {
    if (edificioSeleccionado) {
      setFilterModel({
        items: [
          {
            field: 'id_edificio',
            operator: 'equals',
            value: edificioSeleccionado,
          },
        ],
      });
    } else {
      setFilterModel({ items: [] }); // Quitar filtro
    }
  }, [edificioSeleccionado]);

  //Modelo de columnas invisibles al inicio
  const [columnVisibilityModel, setColumnVisibilityModel] =
    React.useState<GridColumnVisibilityModel>({
      id: false,
      id_edificio: false,
    });

  return (
    <Card sx={sx}>
      <CardHeader title="Estado de Activos" />
      <CardContent>
        <Box sx={{ height: 600, width: '100%' }}>
          <DataGrid
            columnVisibilityModel={columnVisibilityModel}
            onColumnVisibilityModelChange={(newModel) =>
              setColumnVisibilityModel(newModel)
            }
            onRowClick={(params) => {
              window.location.href = paths.dashboard.activoDetail(params.row.id);
            }}
            showToolbar
            rows={activos}
            columns={columns}
            initialState={{
              pagination: {
                paginationModel: {
                  pageSize: 9,
                },
              },
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
            }}
          />
        </Box>
      </CardContent>
    </Card>
  );
}
