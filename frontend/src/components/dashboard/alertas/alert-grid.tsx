/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
'use client';

import * as React from 'react';
import { paths } from '@/paths';
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import { DataGrid, type GridColDef, type GridFilterModel, type GridColumnVisibilityModel} from '@mui/x-data-grid';
import { type Activo } from '@/types/'
import { esES } from '@mui/x-data-grid/locales';

import { activosMock} from '@/mocks/'

//Import de estatus personalizado
import {
  renderStatus,
  STATUS_OPTIONS,
} from './status';

// Configuración rutas de obtención de datos desde db.
import Services from '@/modules/Services'

const activos = activosMock;

const gs = new Services()

const uris = {
  GET: `/parser/activo`
}

const columns: GridColDef<(typeof activos)[number]>[] = [
  { field: 'id', headerName: 'ID', width: 90 },
  {
    field: 'id_edificio',
    headerName: 'ID edificio',
    flex: 0.5, minWidth: 50, // diseño responsivo
  },
  {
    field: 'tipoActivo',
    headerName: 'Tipo de Activo',
    flex: 1.5, minWidth: 120, // diseño responsivo
  },
  {
    field: 'ubicacion',
    headerName: 'Ubicación',
    flex: 1.5, minWidth: 160, // diseño responsivo
  },
  {
    field: 'estado',
    renderCell: renderStatus,
    headerName: 'Estado',
    type: 'singleSelect',
    valueOptions: STATUS_OPTIONS,
    flex: 1,
    minWidth: 100, // diseño responsivo
  },
  {
    field: 'descripcion',
    headerName: 'Descripción',
    description: 'Descripcion de estado del activo',
    flex: 2,
    minWidth: 220, // diseño responsivo
  },
];

export default function DataGridDemo({sx, edificioSeleccionado,}: {sx?: any; edificioSeleccionado?: string | null;}) {
  const [filterModel, setFilterModel] = React.useState<GridFilterModel>({
    items: [],
  });

  const [activos, setActivos] = React.useState<Activo[]>([])
  
  const hasFetchedRef = React.useRef(false)

  // Función para obtener los activos
  const getActivos = async () => {
    try {
      const response = await gs.get(uris.GET)
      console.log("Response:", response)

      // Si no hay datos desde backend, usamos mocks
      if (!response || response.length === 0) {
        console.warn("No se encontraron activos en backend, usando mocks")
        setActivos(activosMock)
        return
      }

      // Transformacion de los datos de la db
      const transformados = response.map((item: any, index: number) => ({
        id: item.activo_id || `B${index + 1}`,
        tipoActivo: item.nombre || 'Activo sin nombre',
        estado: item.estado || 'NN',
        descripcion: item.descripcion || 'NN',
        ubicacion: item.ubicacion || 'Ubicación desconocida',
        id_edificio: item.id_edificio || 'ID no obtenida'
      }))

      console.log("Activos transformados:", transformados)
      setActivos(transformados)
    } catch (error) {
      console.error("Error al obtener los items desde backend, usando mocks", error)
      setActivos(activosMock) // usar mocks si falla la llamada
    }
  }

  React.useEffect(() => {
    if (!hasFetchedRef.current) {
      hasFetchedRef.current = true
      getActivos()
    }
  }, [])

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
                      { setColumnVisibilityModel(newModel); }
                    }
                    onRowClick={(params) => {window.location.href = paths.dashboard.activoDetail(params.row.id);}}
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
                    onFilterModelChange={(newModel) => { setFilterModel(newModel); }}
                    localeText={{
                      ...esES.components.MuiDataGrid.defaultProps.localeText,
                      filterPanelInputLabel: 'Valor a filtrar',
                      filterPanelOperator: 'Operador',
                      filterPanelColumns: 'Filtrar por columna',
                      toolbarColumns: 'Columnas visibles',
                      toolbarFilters: 'Filtros',
                      toolbarExport: 'Exportar',
                      // Traducción de paginación
                      paginationRowsPerPage: 'Activos por página',
                      noRowsLabel: 'No hay activos disponibles',
                      footerTotalRows: 'Total de activos:',
                      footerTotalVisibleRows: (visibleCount, totalCount) =>
                      `${visibleCount.toLocaleString()} de ${totalCount.toLocaleString()}`,
                      footerRowSelected: (count) =>                          count > 1
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
