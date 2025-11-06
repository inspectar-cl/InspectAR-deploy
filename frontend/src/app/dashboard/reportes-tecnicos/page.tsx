/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/

// app/dashboard/reportes-tecnicos/page.tsx
'use client'
import * as React from 'react';
import { useState, useMemo, useEffect } from 'react';
import { 
  Container, 
  Typography, 
  Box, 
  Tabs, 
  Tab, 
  Stack, 
  Card, 
  CardContent,
  CircularProgress,
  Alert
} from '@mui/material';
import GenerarReporte from '@/components/dashboard/reportes-tecnicos/GenerarReporte';
import ListasAcciones from '@/components/dashboard/reportes-tecnicos/ListasAcciones';
import DocumentosAsociados from '@/components/dashboard/reportes-tecnicos/DocumentosAsociados';
import { useUserToken } from '@/hooks/use-usertoken';
import Services from '@/modules/Services';

const gs = new Services();

interface TabConfig {
  label: string;
  component: React.ReactNode;
  roles: string[];
}

interface UserRoleResponse {
  scope?: string;
  error?: string;
}

export default function ReportesPage() {
  const [tab, setTab] = useState(0);
  const { user, isLoading } = useUserToken();
  
  // Estados para rol y datos
  const [userRole, setUserRole] = useState<string | null>(null);
  const [loadingRole, setLoadingRole] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Obtener el rol del usuario desde la API
  useEffect(() => {
    if (isLoading || !user?.token) {
      return;
    }

    const fetchUserRole = async () => {
      setLoadingRole(true);
      setError(null);

      try {
        const response = await gs.authorizedGet('/rol', user.token) as UserRoleResponse;
        
        if (response.error) {
          throw new Error(response.error);
        }

        // El endpoint retorna { scope: "user-type:Root" } o similar
        if (response.scope) {
          // Extraer solo el tipo de usuario después de "user-type:"
          const scopeParts = response.scope.split(':');
          const role = scopeParts.length > 1 ? scopeParts[1] : scopeParts[0];
          setUserRole(role.toLowerCase());
        } else {
          // Fallback si no viene scope
          setUserRole('root');
        }

      } catch (err) {
        console.error('Error al obtener el rol del usuario:', err);
        setError('No se pudo cargar la información del usuario');
        // Fallback: usar un rol por defecto o extraer del token
        setUserRole('root');
      } finally {
        setLoadingRole(false);
      }
    };

    void fetchUserRole();
  }, [user, isLoading]);

  // Definir tabs según el rol
  const allTabs: TabConfig[] = useMemo(() => [
    {
      label: 'Generar Reporte',
      component: <GenerarReporte />,
      roles: ['root', 'tecnico', 'analista']
    },
    {
      label: 'Listas de Acciones',
      component: <ListasAcciones />,
      roles: ['root', 'tecnico']
    },
    {
      label: 'Documentos Asociados',
      component: <DocumentosAsociados />,
      roles: ['root', 'tecnico', 'analista']
    }
  ], []);

  // Filtrar tabs según el rol del usuario
  const availableTabs = useMemo(() => {
    if (!userRole) return [];
    return allTabs.filter(t => t.roles.includes(userRole));
  }, [allTabs, userRole]);

  // Ajustar el tab activo si está fuera de rango
  useEffect(() => {
    if (tab >= availableTabs.length && availableTabs.length > 0) {
      setTab(0);
    }
  }, [availableTabs, tab]);

  const handleChange = (_: React.SyntheticEvent, newValue: number) => { 
    setTab(newValue); 
  };

  // Estado de carga
  if (isLoading || loadingRole) {
    return (
      <Stack spacing={2} alignItems="center" justifyContent="center" sx={{ minHeight: '50vh' }}>
        <CircularProgress />
        <Typography>Cargando información del usuario...</Typography>
      </Stack>
    );
  }

  // Estado de error
  if (error) {
    return (
      <Stack spacing={2}>
        <Typography variant="h4">Reportes Técnicos</Typography>
        <Alert severity="error">
          {error}
        </Alert>
      </Stack>
    );
  }

  // Sin permisos
  if (!userRole || availableTabs.length === 0) {
    return (
      <Stack spacing={2}>
        <Typography variant="h4">Reportes Técnicos</Typography>
        <Alert severity="warning">
          No tienes permisos para acceder a esta sección.
        </Alert>
      </Stack>
    );
  }

  return (
    <Stack spacing={2}>
      <div>
        <Typography variant="h4">Reportes Técnicos</Typography>
        {userRole && (
          <Typography variant="caption" color="text.secondary">
            Rol: {userRole.charAt(0).toUpperCase() + userRole.slice(1)}
          </Typography>
        )}
      </div>
      
      <Container maxWidth="lg">
        <Box sx={{ mt: 3, width: '100%' }}>
          <Tabs value={tab} onChange={handleChange} centered>
            {availableTabs.map((t, index) => (
              <Tab key={index} label={t.label} />
            ))}
          </Tabs>

          <Box sx={{ mt: 3 }}>
            <Card sx={{ width: '100%', height: '100%' }}>
              <CardContent>
                {availableTabs[tab]?.component}
              </CardContent>
            </Card>
          </Box>
        </Box>
      </Container>
    </Stack>
  );
}
