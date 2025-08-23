// app/dashboard/reportes-tecnicos/page.tsx
'use client'

import { useState } from 'react';
import { Container, Typography, Box, Tabs, Tab, Button, TextField, MenuItem, Stack, Card, CardContent } from '@mui/material';
import GenerarReporte from '@/components/dashboard/reportes-tecnicos/GenerarReporte';
import ListasAcciones from '@/components/dashboard/reportes-tecnicos/ListasAcciones';
import DocumentosAsociados from '@/components/dashboard/reportes-tecnicos/DocumentosAsociados';


export default function ReportesPage() {
  const [tab, setTab] = useState(0)
  const handleChange = (_: React.SyntheticEvent, newValue: number) => setTab(newValue)

  return (
    <Stack spacing={2}>
      <div>
        <Typography variant="h4">Reportes Técnicos</Typography>
      </div>
    <Container maxWidth="lg" >
      <Box sx={{ mt: 3, width: '100%' }}>
        <Tabs value={tab} onChange={handleChange} centered>
          <Tab label="Generar Reporte" />
          <Tab label="Listas de Acciones" />
          <Tab label="Documentos Asociados" />
        </Tabs>

        <Box sx={{ mt: 3 }}>
        <Card sx={{ width: '100%', height: '100%' }}>
            <CardContent>
          {tab === 0 && <GenerarReporte />}

          {tab === 1 && <ListasAcciones />}

          {tab === 2 && <DocumentosAsociados />}
                </CardContent>
        </Card>
        </Box>
      </Box>
    </Container>
    </Stack>
  )
}
