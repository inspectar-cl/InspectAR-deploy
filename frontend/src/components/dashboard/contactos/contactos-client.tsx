'use client'

import * as React from 'react';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { DownloadIcon } from '@phosphor-icons/react/dist/ssr/Download';
import { PlusIcon } from '@phosphor-icons/react/dist/ssr/Plus';
import { UploadIcon } from '@phosphor-icons/react/dist/ssr/Upload';

import { ContactosFilters } from './contactos-filters';
import { ContactosTable, type Contacto } from './contactos-table';
import { CustomAlert } from '@/components/dashboard/alert-popups/customAlertPopup';

import { useEffect, useState } from 'react';
import Services from '@/modules/Services'

interface Props {
  contactos: Contacto[];
  especialidades: string[];
}
const gs = new Services()

export function ContactosClient({ contactos, especialidades }: Props) {
  const [filter, setFilter] = React.useState({ search: '', especialidad: null as string | null });
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(5);

  const [alertState, setAlertState] = React.useState({
    open: false,
    title: '',
    message: '',
    severity: 'success' as 'success' | 'info' | 'warning' | 'error',
  });

  const handlePageChange = (_: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleRowsPerPageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleFilterChange = (newFilters: {search: string, especialidad: null | string}, resetPage = false) => {
    setFilter(newFilters);
    if (resetPage) setPage(0);
  };

  const handleSendRequest = async (selectedRows: Contacto[]) => {
    console.log('Contactos seleccionados para enviar solicitud:', selectedRows);
    
    try {
      // Enviar notificación a cada contacto seleccionado
      const promises = selectedRows.map(async (contacto) => {
        const payload = {
          activo_id: 1, // Valor temporal como int
          message: "Solicitud de contacto técnico enviada exitosamente",
          priority: "high", // Valor temporal
          technician_email: contacto.email,
          user_email: "admin@inspectar.cl", // Valor temporal
          user_name: "Admin InspectAR" // Valor temporal
        };
        
        console.log(`Enviando notificación a ${contacto.email}:`, payload);
        return await gs.post("/notificacion/technician/contact", payload);
      });

      // Esperar a que todas las llamadas se completen
      const results = await Promise.all(promises);
      console.log('Resultados de envío:', results);

      // Realizar esta alerta si sale bien el envio por la API
      setAlertState({
        open: true,
        title: '¡Éxito!',
        message: `Se ha enviado la solicitud a ${selectedRows.length} contacto(s).`,
        severity: 'success',
      });

    } catch (error) {
      console.error('Error al enviar notificaciones:', error);
      
      // Realizar esta alerta si sale mal el envio por API
      setAlertState({
        open: true,
        title: '¡Ups!',
        message: `Hubo un problema enviando las solicitud(es).`,
        severity: 'error',
      });
    }
    
    setTimeout(() => {
      setAlertState(prev => ({ ...prev, open: false }));
    }, 4000);
  };

  const handleAlertClose = () => {
    setAlertState(prev => ({ ...prev, open: false }));
  };

  const filtered = contactos.filter((c) => {
    const matchesSearch =
        c.name.toLowerCase().includes(filter.search.toLowerCase()) ||
        c.email.toLowerCase().includes(filter.search.toLowerCase());

    const matchesEspecialidad =
        !filter.especialidad || c.especialidad === filter.especialidad;

    return matchesSearch && matchesEspecialidad;
  });

  const paginated = applyPagination(filtered, page, rowsPerPage);

  return (
    <Stack spacing={3}>
      <CustomAlert
        title={alertState.title}
        message={alertState.message}
        severity={alertState.severity}
        open={alertState.open}
        onClose={handleAlertClose}
      />
      <Stack direction="row" spacing={3}>
        <Stack spacing={1} sx={{ flex: '1 1 auto' }}>
          <Typography variant="h4">Contactos</Typography>
          <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
            { /* Botones de importar y exportar, descomentar si es necesario
            <Button color="inherit" startIcon={<UploadIcon fontSize="var(--icon-fontSize-md)" />}>
              Import
            </Button>
            <Button color="inherit" startIcon={<DownloadIcon fontSize="var(--icon-fontSize-md)" />}>
              Export
            </Button>
            */ }
          </Stack>
        </Stack>
        {/* <div>
          <Button startIcon={<PlusIcon fontSize="var(--icon-fontSize-md)" />} variant="contained">
            Añadir
          </Button>
        </div> */}
      </Stack>
      <ContactosFilters onFilterChange={handleFilterChange} especialidades={especialidades} />
      <ContactosTable
        count={filtered.length}
        page={page}
        rows={paginated}
        rowsPerPage={rowsPerPage}
        onPageChange={handlePageChange}
        onRowsPerPageChange={handleRowsPerPageChange}
        onSendRequest={handleSendRequest}
      />
    </Stack>
  );
}

function applyPagination(rows: Contacto[], page: number, rowsPerPage: number): Contacto[] {
  return rows.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);
}