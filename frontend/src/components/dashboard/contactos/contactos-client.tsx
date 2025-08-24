'use client'

import * as React from 'react';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { DownloadIcon } from '@phosphor-icons/react/dist/ssr/Download';
import { PlusIcon } from '@phosphor-icons/react/dist/ssr/Plus';
import { UploadIcon } from '@phosphor-icons/react/dist/ssr/Upload';

import { ContactosFilters } from './contactos-filters';
import { ContactosTable, Contacto } from './contactos-table';
import { CustomAlert } from '@/components/dashboard/alert-popups/customAlertPopup';

interface Props {
  contactos: Contacto[];
}

export function ContactosClient({ contactos }: Props) {
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

  const handleSendRequest = (selectedRows: Contacto[]) => {
    console.log('Contactos seleccionados para enviar solicitud:', selectedRows);
    // Aqui llamar a la API para enviar notificacion a los contactos seleccionados

    //Realizar esta alerta si sale bien el envio por la API
    setAlertState({
      open: true,
      title: '¡Éxito!',
      message: `Se ha enviado la solicitud a ${selectedRows.length} contacto(s).`,
      severity: 'success',
    });

    // Realizar esta alerta si sale mal el envio por API (descomentar)
    {/*
    setAlertState({
      open: true,
      title: '¡Ups!',
      message: `Hubo un problema enviando de solicitud(es).`,
      severity: 'error',
    });
    */}
    
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
        <div>
          <Button startIcon={<PlusIcon fontSize="var(--icon-fontSize-md)" />} variant="contained">
            Añadir
          </Button>
        </div>
      </Stack>
      <ContactosFilters onFilterChange={handleFilterChange}/>
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