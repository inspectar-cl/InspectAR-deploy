/* eslint-disable @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-return, @typescript-eslint/no-unsafe-call -- API integration requires type flexibility */
// app/dashboard/agregar-datos/page.tsx
'use client';

import * as React from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { FormProvider, useForm } from 'react-hook-form';
import type { 
  SolicitudFormData, 
  ActivoData, 
  EdificioData, 
  TecnicoData 
} from '@/types/formulario';
import {
  Container,
  Paper,
  Typography,
  Box,
  CircularProgress,
  Alert,
  Tabs,
  Tab,
  Button,
} from '@mui/material';
import { paths } from '@/paths';
import { useUserToken } from '@/hooks/use-usertoken';

// Importa tus formularios
import { EdificioDataForm } from '@/components/dashboard/agregar-datos/edificio-datos-form';
import { ActivoForm } from '@/components/dashboard/agregar-datos/activo-datos-form';
import { TecnicoForm } from '@/components/dashboard/agregar-datos/tecnico-datos-form';

type TabValue = 'Edificio' | 'Activo' | 'Técnico';
const TABS: { value: TabValue; label: string }[] = [
  { value: 'Edificio', label: 'Edificio' },
  { value: 'Activo', label: 'Activo' },
  { value: 'Técnico', label: 'Técnico' },
];

// --- Helper para transformar datos de Activo ---
function transformarDatosActivo(datosFormulario: ActivoData) {
  let tipoApi: string;
  switch (datosFormulario.tipoActivo) {
    case 'BombaDeAgua': tipoApi = 'bomba de agua'; break;
    case 'Ascensor': tipoApi = 'ascensor'; break;
    case 'PanelElectrico': tipoApi = 'panel electrico'; break;
    default: tipoApi = 'NN';
  }
  return {
    nombre: datosFormulario.nombre,
    tipo: tipoApi,
    descripcion: datosFormulario.descripcion || '',
    ubicacion: datosFormulario.ubicacion,
    edificio_id: datosFormulario.edificioId,
  };
}
// --- Helper para transformar datos de Edificio ---
function transformarDatosEdificio(datosFormulario: EdificioData) {
  return {
    nombre: datosFormulario.nombre,
    direccion: datosFormulario.direccion,
    latitud: datosFormulario.latitud,
    longitud: datosFormulario.longitud,
  };
}
// --- (Añade un 'transformarDatosTecnico' si es necesario) ---


export default function PaginaAgregarDatos() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user } = useUserToken();
  const [apiError, setApiError] = React.useState<string | null>(null);
  
  // --- 1. Estado del Tab (controlado por query param) ---
  const tabParam = searchParams.get('tab') as TabValue;
  const dataParam = searchParams.get('data'); // <-- ¡Lee el 'data' param!
  
  const [activeTab, setActiveTab] = React.useState<TabValue>(tabParam || 'Edificio');

  const defaultValues = React.useMemo(() => {
    // Caso 1: Vienes de una solicitud (pre-llenar)
    if (dataParam) {
      try {
        const decodedString = decodeURIComponent(dataParam);
        const solicitud = JSON.parse(decodedString) as SolicitudFormData;
        
        // Sincroniza el tab con los datos recibidos
        setActiveTab(solicitud.tipoSolicitud); 
        
        return {
          tipoSolicitud: solicitud.tipoSolicitud,
          asunto: solicitud.asunto,
          detalles: solicitud.detalles,
          datosEspecificos: solicitud.datosEspecificos,
        } as SolicitudFormData;
      } catch (e) {
        // Fallback a formulario vacío si los datos son corruptos
        return { tipoSolicitud: activeTab, datosEspecificos: {} } as SolicitudFormData;
      }
    }
    
    // Caso 2: Vienes del Sidebar (formulario vacío)
    return { tipoSolicitud: activeTab, datosEspecificos: {} } as SolicitudFormData;
  }, [dataParam, activeTab]);

  const methods = useForm<SolicitudFormData>({
    defaultValues,
  });

  // Sincronizar el estado del tab con un query param ?tab=...
  React.useEffect(() => {
    const tab = searchParams.get('tab') as TabValue;
    if (tab && TABS.find((t) => t.value === tab)) {
      setActiveTab(tab);
    }
  }, [searchParams]);

  React.useEffect(() => {
    // Resetea el formulario solo si los defaultValues cambian
    // (Ej: al cargar la página con ?data=... o al cambiar de tab)
    methods.reset(defaultValues);
  }, [defaultValues, methods]);

  const handleTabChange = (event: React.SyntheticEvent, newTab: TabValue) => {
    setActiveTab(newTab);
    // Limpia la URL (quita el ?data=...) al cambiar de tab manualmente
    router.push(`${paths.dashboard.agregardatos}?tab=${newTab}`);
  };

  // --- 4. Lógica de onSubmit ---
  const onSubmit = async (data: SolicitudFormData): Promise<void> => {
    setApiError(null);
    if (!user?.token) {
      setApiError('No estás autenticado.');
      return;
    }
    
    // Verificamos que el tipo de dato coincida con el tab activo
    if (data.tipoSolicitud !== activeTab) {
      setApiError(`Error: El formulario no coincide con el tab (${data.tipoSolicitud} vs ${activeTab})`);
      return;
    }

    try {
      let endpoint = '';
      let payload: any;

      switch (data.tipoSolicitud) {
        case 'Activo':
          endpoint = '/api/crear-activo';
          payload = transformarDatosActivo(data.datosEspecificos as ActivoData);
          break;
        case 'Edificio':
          endpoint = '/api/crear-edificio'; 
          payload = transformarDatosEdificio(data.datosEspecificos as EdificioData);
          break;
        case 'Técnico':
          endpoint = '/api/crear-tecnico';
          payload = data.datosEspecificos as TecnicoData; // Asume transformación si es necesaria
          break;
        default:
          throw new Error('Tipo de formulario no reconocido');
      }


      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${user.token}`,
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => null);
        throw new Error(errorData?.error || `Error del servidor: ${response.status}`);
      }

      // alert(`${data.tipoSolicitud} creado con éxito.`);
      methods.reset();

    } catch (err) {
      setApiError(err instanceof Error ? err.message : 'Ocurrió un error desconocido');
    }
  };

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Paper sx={{ p: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Centro de Gestión de Datos
        </Typography>

        {/* --- 5. FormProvider envuelve TODO --- */}
        <FormProvider {...methods}>
          <form onSubmit={methods.handleSubmit(onSubmit)}>
            
            <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
              <Tabs value={activeTab} onChange={handleTabChange}>
                {TABS.map((tab) => (
                  <Tab key={tab.value} label={tab.label} value={tab.value} />
                ))}
              </Tabs>
            </Box>

            {apiError ? <Alert severity="error" sx={{ mb: 3 }}>
                {apiError}
              </Alert> : null}

            {/* Renderizado condicional del formulario */}
            <Box>
              {activeTab === 'Edificio' && <EdificioDataForm />}
              {activeTab === 'Activo' && <ActivoForm />}
              {activeTab === 'Técnico' && <TecnicoForm />}
            </Box>

            <Button
              type="submit"
              variant="contained"
              fullWidth
              sx={{ mt: 4 }}
              disabled={methods.formState.isSubmitting}
            >
              {methods.formState.isSubmitting ? (
                <CircularProgress size={24} color="inherit" />
              ) : (
                `Guardar ${activeTab}`
              )}
            </Button>
          </form>
        </FormProvider>
      </Paper>
    </Container>
  );
}