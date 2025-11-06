// app/dashboard/gestion/page.tsx
'use client';

import * as React from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { Container, Paper, Typography, Box, Tabs, Tab, Alert } from '@mui/material';

// --- Importarás un componente DataGrid para cada Tab ---
// (Los crearemos en el paso 3)
import { EdificiosDataGrid } from '@/components/dashboard/edicion-datos/edificios-datagrid';
import { ActivosDataGrid } from '@/components/dashboard/edicion-datos/activos-datagrid';
//import { SensoresDataGrid } from '@/components/dashboard/edicion-datos/sensores-datagrid';
import { TecnicosDataGrid } from '@/components/dashboard/edicion-datos/tecnicos-datagrid';

type TabValue = 'edificios' | 'activos' | 'sensores' | 'tecnicos';
const TABS: { value: TabValue; label: string }[] = [
  { value: 'edificios', label: 'Edificios' },
  { value: 'activos', label: 'Activos' },
  //{ value: 'sensores', label: 'Sensores' },
  { value: 'tecnicos', label: 'Técnicos' },
];

export default function PaginaGestionDatos() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [activeTab, setActiveTab] = React.useState<TabValue>('edificios');

  // Sincronizar el estado del tab con un query param ?tab=...
  React.useEffect(() => {
    const tab = searchParams.get('tab') as TabValue;
    if (tab && TABS.find(t => t.value === tab)) {
      setActiveTab(tab);
    }
  }, [searchParams]);

  const handleTabChange = (event: React.SyntheticEvent, newTab: TabValue) => {
    setActiveTab(newTab);
    // Actualiza la URL
    router.push(`/dashboard/gestion?tab=${newTab}`);
  };

  return (
    <Container maxWidth="xl" sx={{ py: 4 }}>
      <Paper sx={{ p: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Gestión de Datos
        </Typography>
        <Alert severity="info" sx={{ mb: 3 }}>
          En esta sección puedes visualizar, editar y eliminar los datos maestros del sistema.
        </Alert>

        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Tabs value={activeTab} onChange={handleTabChange} aria-label="Tabs de gestión de datos">
            {TABS.map((tab) => (
              <Tab key={tab.value} label={tab.label} value={tab.value} />
            ))}
          </Tabs>
        </Box>

        {/* Renderizado condicional de la DataGrid activa */}
        <Box>
          {activeTab === 'edificios' && <EdificiosDataGrid />}
          {activeTab === 'activos' && <ActivosDataGrid />}
          {/*activeTab === 'sensores' && <SensoresDataGrid />*/}
          {activeTab === 'tecnicos' && <TecnicosDataGrid />}
        </Box>
      </Paper>
    </Container>
  );
}