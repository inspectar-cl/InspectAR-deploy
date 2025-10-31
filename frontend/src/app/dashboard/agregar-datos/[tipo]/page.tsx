// app/gestion/crear/[tipo]/page.tsx
'use client';

import * as React from 'react';
import { useSearchParams, useRouter, useParams } from 'next/navigation';
import { FormProvider, useForm } from 'react-hook-form';
import type { SolicitudFormData, TipoSolicitud } from '@/types/formulario';
import { 
    Container, Paper, Typography, Box, 
    CircularProgress, Alert, Tabs, Tab, Button
} from '@mui/material';

import { SLUG_TO_TIPO, TABS } from '@/utils/solicitud-utils';
import { EdificioDataForm } from '@/components/dashboard/agregar-datos/edificio-datos-form';
import { ActivoForm } from '@/components/dashboard/agregar-datos/activo-datos-form';
import { SensorForm } from '@/components/dashboard/agregar-datos/sensor-datos-form';
import { TecnicoForm } from '@/components/dashboard/agregar-datos/tecnico-datos-form';

function GestionFormulario() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const params = useParams();
  const tipoSlug = params.tipo as string; // 'params.tipo' viene de [tipo]

  const activeTabTipo = SLUG_TO_TIPO[tipoSlug];
  const [parseError, setParseError] = React.useState<string | null>(null);

  const defaultValues = React.useMemo(() => {
    const dataParam = searchParams.get('data');
    if (dataParam && activeTabTipo) {
      try {
        const decodedString = decodeURIComponent(dataParam);
        const solicitud = JSON.parse(decodedString);
        
        if (solicitud.tipoSolicitud !== activeTabTipo) {
          throw new Error('Conflicto de datos y URL.');
        }

        return {
          tipoSolicitud: solicitud.tipoSolicitud,
          asunto: solicitud.asunto,
          detalles: solicitud.detalles,
          datosEspecificos: solicitud.datosEspecificos,
        } as SolicitudFormData;
      } catch (e) {
        console.error("Error al parsear datos:", e);
        setParseError("No se pudieron cargar los datos de la solicitud.");
        return { tipoSolicitud: activeTabTipo } as SolicitudFormData;
      }
    }
    // Formulario nuevo (sin datos)
    return { tipoSolicitud: activeTabTipo, datosEspecificos: {} } as SolicitudFormData;
  }, [searchParams, activeTabTipo]);

  const methods = useForm<SolicitudFormData>({
    defaultValues: defaultValues,
  });

  //Resetear el formulario si el defaultValues (derivado de la URL) cambia
  React.useEffect(() => {
    methods.reset(defaultValues);
  }, [defaultValues, methods]);


  const handleTabChange = (event: React.SyntheticEvent, newSlug: string) => {
    router.push(`/dashboard/agregar-datos/${newSlug}`);
  };

  const onSubmit = async (data: SolicitudFormData) => {
    console.log(`Enviando datos para: ${data.tipoSolicitud}`, data);
    // Aqui debe estar la logica de la API para guardar la data
    alert(`Simulación: ${data.tipoSolicitud} creado.`);
    // router.push('/dashboard/solicitudes'); //Agregar?
  };

  if (!activeTabTipo) {
    return <Alert severity="error">Tipo de formulario no válido.</Alert>;
  }
  if (parseError) {
    return <Alert severity="error">{parseError}</Alert>;
  }

  // 7. Renderizado
  return (
    <FormProvider {...methods}>
      <form onSubmit={methods.handleSubmit(onSubmit)}>
        
        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Tabs value={tipoSlug} onChange={handleTabChange}>
            {TABS.map((tab) => (
              <Tab key={tab.slug} label={tab.label} value={tab.slug} />
            ))}
          </Tabs>
        </Box>

        <Box>
          {activeTabTipo === 'Edificio' && <EdificioDataForm />}
          {activeTabTipo === 'Activo' && <ActivoForm />}
          {activeTabTipo === 'Sensor' && <SensorForm />}
          {activeTabTipo === 'Técnico' && <TecnicoForm />}
        </Box>
        
        <Button 
          type="submit" 
          variant="contained" 
          fullWidth 
          sx={{ mt: 4 }}
          disabled={methods.formState.isSubmitting}
        >
          Guardar {activeTabTipo}
        </Button>
      </form>
    </FormProvider>
  );
}

export default function PaginaGestionCrear({ params }: { params: { tipo: string } }) {
  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Paper sx={{ p: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Centro de Gestión de Datos
        </Typography>

        <React.Suspense fallback={
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
            <CircularProgress />
          </Box>
        }>
          <GestionFormulario />
        </React.Suspense>
        
      </Paper>
    </Container>
  );
}