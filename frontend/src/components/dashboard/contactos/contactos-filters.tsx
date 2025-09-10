/* eslint-disable @typescript-eslint/explicit-function-return-type -- Función generadora de UI, tipo inferido*/
'use client';

import * as React from 'react';
import { TextField, Stack, Card } from '@mui/material'
import Autocomplete from '@mui/material/Autocomplete';
import InputAdornment from '@mui/material/InputAdornment';
import OutlinedInput from '@mui/material/OutlinedInput';
import { MagnifyingGlassIcon } from '@phosphor-icons/react/dist/ssr/MagnifyingGlass';

interface ContactosFiltersProps  {
  onFilterChange: (filters: { search: string; especialidad: string | null }, resetPage?: boolean) => void;
  especialidades: string[];
}

export function ContactosFilters({ onFilterChange, especialidades }: ContactosFiltersProps ): React.JSX.Element {
  const [search, setSearch] = React.useState('');
  const [tipo, setTipo] = React.useState<string | null>(null);

  const updateFilters = (newSearch: string, newTipo: string | null) => {
    onFilterChange({ search: newSearch, especialidad: newTipo }, true);
  };

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    setSearch(value);
    updateFilters(value, tipo);
  };

  const handleTipoChange = (_: React.SyntheticEvent, newValue: string | null) => {
    setTipo(newValue);
    updateFilters(search, newValue);
  };

  return (
    <Card sx={{ p: 2 }}>
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={4}>
        <OutlinedInput
          value={search}
          onChange={handleSearchChange}
          fullWidth
          placeholder="Buscar contacto"
          startAdornment={
            <InputAdornment position="start">
              <MagnifyingGlassIcon fontSize="var(--icon-fontSize-md)" />
            </InputAdornment>
          }
          sx={{ maxWidth: '500px' }}
        />

        <Autocomplete
          value={tipo}
          onChange={handleTipoChange}
          options={especialidades}
          sx={{ width: 300 }}
          renderInput={(params) => <TextField {...params} label="Especialidad" />}
        />
      </Stack>
    </Card>
  );
}
