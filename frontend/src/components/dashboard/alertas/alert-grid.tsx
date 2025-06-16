'use client';

import * as React from 'react';
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import { DataGrid, GridColDef, QuickFilter} from '@mui/x-data-grid';
import { esES } from '@mui/x-data-grid/locales';

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

const activos: Activo[] = [
  { id: 'A1', tipoActivo: 'Ascensor #1', estado: 'OK', descripcion: 'Funciona correctamente', ubicacion: 'Edificio A, Santiago'},
  { id: 'A2', tipoActivo: 'Ascensor #2', estado: 'OK', descripcion: 'Funciona correctamente', ubicacion: 'Edificio A, Santiago' },
  { id: 'A3', tipoActivo: 'Ascensor #3', estado: 'Medio', descripcion: 'Mantenimiento programado no realizado', ubicacion: 'Edificio A, Santiago' },
  { id: 'A4', tipoActivo: 'Ascensor #4', estado: 'Medio', descripcion: 'Puerta atascada, requiere revisión', ubicacion: 'Edificio B, Santiago' },
  { id: 'C1', tipoActivo: 'Caldera #1', estado: 'Crítico', descripcion: 'Temperatura fuera de rango, riesgo de daño', ubicacion: 'Edificio A, Santiago' },
  { id: 'C2', tipoActivo: 'Caldera #2', estado: 'Crítico', descripcion: 'Anomalía detectada – revisar urgentemente', ubicacion: 'Edificio B, Santiago' },
  { id: 'C3', tipoActivo: 'Caldera #3', estado: 'Crítico', descripcion: 'Fuga de gas detectada', ubicacion: 'Edificio C, Santiago' },
  { id: 'C4', tipoActivo: 'Caldera #4', estado: 'Crítico', descripcion: 'Presión excesiva, riesgo de explosión', ubicacion: 'Edificio C, Santiago' },
  { id: 'B1', tipoActivo: 'Bomba de Agua #1', estado: 'Crítico', descripcion: 'Motor sobrecalentado, riesgo de falla', ubicacion: 'Edificio A, Santiago' },
  { id: 'B2', tipoActivo: 'Bomba de Agua #2', estado: 'Crítico', descripcion: 'Fuga de presión, riesgo de parada', ubicacion: 'Edificio A, Santiago' },
  { id: 'B3', tipoActivo: 'Bomba de Agua #3', estado: 'Crítico', descripcion: 'Riesgo crítico de falla en 7 días', ubicacion: 'Edificio B, Santiago' },
  { id: 'B4', tipoActivo: 'Bomba de Agua #4', estado: 'Crítico', descripcion: 'Nivel de agua crítico', ubicacion: 'Edificio B, Santiago' },
  { id: 'E1', tipoActivo: 'Sistema Eléctrico #1', estado: 'Medio', descripcion: 'Sobrecarga detectada, monitorear consumo', ubicacion: 'Edificio A, Santiago' },
  { id: 'E2', tipoActivo: 'Sistema Eléctrico #2', estado: 'Medio', descripcion: 'Pico de voltaje registrado', ubicacion: 'Edificio B, Santiago' },
  { id: 'E3', tipoActivo: 'Sistema Eléctrico #3', estado: 'Medio', descripcion: 'Variación de frecuencia, revisar panel', ubicacion: 'Edificio C, Santiago' },
]

const columns: GridColDef<(typeof activos)[number]>[] = [
  { field: 'id', headerName: 'ID', width: 90 },
  {
    field: 'tipoActivo',
    headerName: 'Tipo de Activo',
    width: 150,
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
    width: 160,
  },
  {
    field: 'ubicacion',
    headerName: 'Ubicación',
    width: 160,
  },
];

export default function DataGridDemo({ sx }: { sx?: any }) {
  return (
    <Card sx={sx}>
      <CardHeader title="Estado de Activos" />
        <CardContent>
            <Box sx={{ height: 400, width: '100%' }}>
                <DataGrid
                    showToolbar
                    rows={activos}
                    columns={columns}
                    initialState={{
                      pagination: {
                          paginationModel: {
                          pageSize: 5,
                          },
                      },
                    }}
                    pageSizeOptions={[5]}
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
