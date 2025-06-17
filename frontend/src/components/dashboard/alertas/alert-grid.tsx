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

type Activo = {
  id: string
  tipoActivo: string
  estado: 'OK' | 'Medio' | 'Crítico'
  descripcion: string
  ubicacion: string
}

type Alerta = {
  id: string
  tipoActivo: string
  descripcion: string
  estado: 'Medio' | 'Crítico'
  ubicacion: string
  fecha: string
  atendida: boolean
}

// Configuración rutas de obtención de datos desde db.
import Services from '@/modules/Services'
const activos = activosMock;

const gs = new Services()

const uris = {
  GET: `/activo`
}

const columns: GridColDef<Activo>[] = [
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

export default function DataGridDemo({ sx }: { sx?: any }) {
  const [activos, setActivos] = React.useState<Activo[]>([])
  const hasFetchedRef = React.useRef(false)

  const getActivos = async () => {
    try {
      const response = await gs.get(uris.GET)
      console.log("Response:", response)

      // Transformación de los datos de la db
      const transformados = response.map((item: any, index: number) => ({
        id: item.activo_id || `B${index + 1}`, // o usa item.id si lo prefieres
        tipoActivo: item.nombre || 'Activo sin nombre',
        estado: item.estado || 'NN', // Esto puedes modificarlo con lógica si tienes
        descripcion: item.descripcion || 'NN', // Aquí también puedes usar lógica si quieres
        ubicacion: item.ubicacion || 'Ubicación desconocida'
      }));

      console.log("Activos transformados:", transformados)
      setActivos(transformados)
      return activos
    } catch (error) {
      console.error("Error al obtener los items.", error)
    }
  }

  React.useEffect(() => {
    if (!hasFetchedRef.current) {
      hasFetchedRef.current = true
      getActivos()
    }
  }, [])
  return (
    <Card sx={sx}>
      <CardHeader title="Estado de Activos" />
        <CardContent>
            <Box>
              activos: {JSON.stringify(activos, null, 2)}
            </Box>
            <Box sx={{ height: 600, width: '100%' }}>
                <DataGrid 
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
