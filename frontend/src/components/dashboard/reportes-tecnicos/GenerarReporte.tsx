/* eslint-disable @typescript-eslint/explicit-function-return-type -- Componente React, tipos inferidos automáticamente */
'use client';

import { useEffect, useState } from 'react';
import {
  Box, Button, TextField, Select, MenuItem, Snackbar, Alert,
  FormControl, InputLabel, Checkbox, ListItemText, OutlinedInput,
  Stack, Typography, Card, CardContent, CardActions
} from '@mui/material';
import type {SelectChangeEvent} from '@mui/material'
import Services from '@/modules/Services';
import type { SavedSignature } from '@/components/dashboard/reportes-tecnicos/SignatureDialog'; 
import SignatureDialog from '@/components/dashboard/reportes-tecnicos/SignatureDialog';
import { useSignatures } from '@/hooks/use-signatures';
import { useUserToken } from '@/hooks/use-usertoken';

const gs = new Services();
const BASE_URL = process.env.NEXT_PUBLIC_API_GATEWAY_URL || '/api';

interface ActivoReporte {
  id: string;
  nombre: string;
}
interface ActivosResponse {
  activos?: ActivoReporte[];
}

function GenerarReporte() {
  const {user, isLoading} = useUserToken();
  const [activos, setActivos] = useState<ActivoReporte[]>([]);
  const [activo, setActivo] = useState<string>('');
  const [observaciones, setObservaciones] = useState<string>('');
  const [pdfUrl, setPdfUrl] = useState<string | null>(null);
  const [mensaje, setMensaje] = useState<string | null>(null);
  const [camposSeleccionados, setCamposSeleccionados] = useState<string[]>([
    'ubicacion',
    'historial_mantenimientos',
    'ultima_acciones',
    'datos_sensores',
  ]);
  const opcionesCampos = [
    { value: 'ubicacion', label: 'Ubicación' },
    { value: 'historial_mantenimientos', label: 'Historial de Mantenimientos' },
    { value: 'ultima_acciones', label: 'Últimas Acciones' },
    { value: 'datos_sensores', label: 'Datos de Sensores' },
  ];

  const [signatureDialogOpen, setSignatureDialogOpen] = useState(false);
  const { signatures, addSignature, removeSignature } = useSignatures();
  const [selectedSignatureDataUrl, setSelectedSignatureDataUrl] = useState<string | null>(null);
  const [selectedSignatureId, setSelectedSignatureId] = useState<string>('');

  // cargar activos
  useEffect(() => {
    const fetchActivos = async () => {
      try {
        if (isLoading || !user) {
          return
        }

        const data = await gs.authorizedGet('/obtener-todos-activos', user.token) as ActivosResponse;
        

        if (data && Array.isArray(data.activos)) {
          setActivos(data.activos);
        } else {
          setActivos([]);
        }
      } catch (error) {
        setActivos([]);
      }
    };
    void fetchActivos();
  }, [isLoading, user]);

  // actualizar dataUrl al elegir una firma guardada
  useEffect(() => {
    if (!selectedSignatureId) {
      setSelectedSignatureDataUrl(null);
      return;
    }
    const found = signatures.find((s) => s.id === selectedSignatureId);
    setSelectedSignatureDataUrl(found?.dataUrl ?? null);
  }, [selectedSignatureId, signatures]);

  // limpiar firma al cambiar de activo 
  useEffect(() => {
    setSelectedSignatureId('');
    setSelectedSignatureDataUrl(null);
  }, [activo]);

  // callback del diálogo
  const handleUseSignature = (dataUrl: string, persist?: boolean) => {
    const id = crypto.randomUUID();
    setSelectedSignatureId(id);
    setSelectedSignatureDataUrl(dataUrl);

    if (persist) {
      const item: SavedSignature = {
        id,
        name: `Firma ${new Date().toLocaleString()}`,
        dataUrl,
        source: 'drawn',
        createdAt: new Date().toISOString(),
      };
      addSignature(item);
    }
    setSignatureDialogOpen(false);
  };

  const handleRemoveSelectedSignature = () => {
    if (selectedSignatureId) {
      removeSignature(selectedSignatureId);
      setSelectedSignatureId('');
    }
    setSelectedSignatureDataUrl(null);
  };

  // exportar PDF
  const handleExportarPDF = async () => {
    if (!activo) return;
    if (isLoading || !user) {
      setMensaje('Usuario no autenticado');
      return;
    }
    
    try {
      const payload = {
        campos: camposSeleccionados,
      };


      // Hacer petición con fetch para obtener el blob
      const response = await fetch(`${BASE_URL}/pdf-reporte/${activo}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${user.token}`,
        },
        body: JSON.stringify(payload),
      });


      if (!response.ok) {
        throw new Error(`Error ${response.status}: ${response.statusText}`);
      }

      // Obtener el PDF como blob
      const blob = await response.blob();

      // Crear URL local para el blob
      const url = window.URL.createObjectURL(blob);

      const activoSeleccionado = activos.find((a) => a.id === activo);
      const nombreActivo = activoSeleccionado ? activoSeleccionado.nombre.replace(/\s+/g, '_') : 'Reporte';
      const fecha = new Date().toISOString().split('T')[0].replace(/-/g, '');
      const nombreArchivo = `Reporte_${nombreActivo}_${fecha}.pdf`;

      const link = document.createElement('a');
      link.href = url;
      link.download = nombreArchivo;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      // Limpiar la URL después de un tiempo
      setTimeout(() => {
        window.URL.revokeObjectURL(url);
      }, 10000);

    } catch (error) {
      setMensaje('No se pudo exportar el PDF');
    }
  };

  // vista previa PDF
  const handleVistaPreviaPDF = async () => {
    if (!activo) return;
    if (isLoading || !user) {
      setMensaje('Usuario no autenticado');
      return;
    }
    
    try {
      const payload = {
        campos: camposSeleccionados,
      };

      // Hacer petición con fetch para obtener el blob
      const response = await fetch(`${BASE_URL}/pdf-reporte/${activo}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${user.token}`,
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`Error ${response.status}: ${response.statusText}`);
      }

      // Obtener el PDF como blob
      const blob = await response.blob();

      // Crear URL local para el blob
      const url = window.URL.createObjectURL(blob);

      // Establecer la URL para la vista previa
      setPdfUrl(url);

    } catch (error) {
      setMensaje('No se pudo generar la vista previa');
    }
  };

  return (
    <Box>
      {/* aqui van los activos de la BD */}
      <Box sx={{ display: 'flex', gap: 2, mb: 3, flexWrap: 'wrap', alignItems: 'flex-start' }}>
        <FormControl sx={{ minWidth: 240 }}>
          <InputLabel id="activo-label" shrink>Seleccionar Activo</InputLabel>
          <Select
            labelId="activo-label"
            value={activo}
            label="Seleccionar Activo"
            onChange={(e) => {
              setActivo(e.target.value);
            }}
            input={<OutlinedInput label="Seleccionar Activo" />}
            displayEmpty
            renderValue={(selected) => {
              if (!selected) return <em>Ninguno</em>;
              const a = activos.find((x) => x.id === selected);
              return a?.nombre ?? selected;
            }}
          >
            <MenuItem value="">
              <em>Ninguno</em>
            </MenuItem>
            {Array.isArray(activos) &&
              activos.map((a) => (
                <MenuItem key={a.id} value={a.id}>
                  {a.nombre}
                </MenuItem>
              ))}
          </Select>
        </FormControl>

        <Button variant="contained" onClick={handleVistaPreviaPDF} disabled={!activo}>
          Vista Previa PDF
        </Button>
        <Button variant="contained" onClick={handleExportarPDF} disabled={!activo}>
          Exportar PDF
        </Button>
      </Box>

      {/* Firma del técnico */}
      <Card variant="outlined" sx={{ mb: 2 }}>
        <CardContent>
          <Typography variant="h6" sx={{ mb: 1 }}>
            Firma del técnico
          </Typography>

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems="center">
            <Button variant="outlined" onClick={() => { setSignatureDialogOpen(true); }}>
              Agregar firma
            </Button>

            <FormControl sx={{ minWidth: 240 }}>
              <InputLabel id="firmas-guardadas-label" shrink>
                Firmas guardadas
              </InputLabel>
              <Select
                labelId="firmas-guardadas-label"
                value={selectedSignatureId}
                onChange={(e) => { setSelectedSignatureId(e.target.value); }}
                input={<OutlinedInput label="Firmas guardadas" />}
                displayEmpty
                renderValue={(selected) => {
                  if (!selected) return <em>Ninguna</em>;
                  const sig = signatures.find((s) => s.id === selected);
                  return sig?.name ?? selected;
                }}
              >
                <MenuItem value="">
                  <em>Ninguna</em>
                </MenuItem>
                {signatures.map((sig) => (
                  <MenuItem key={sig.id} value={sig.id}>
                    {sig.name || sig.id}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            {selectedSignatureDataUrl ? <Stack direction="row" spacing={1} alignItems="center" sx={{ ml: { sm: 'auto' } }}>
                <Box
                  component="img"
                  src={selectedSignatureDataUrl}
                  alt="Firma seleccionada"
                  sx={{
                    height: 60,
                    maxWidth: 240,
                    border: '1px solid',
                    borderColor: 'divider',
                    borderRadius: 1,
                    p: 0.5,
                    background: '#fff',
                  }}
                />
                <Button color="error" variant="outlined" onClick={handleRemoveSelectedSignature}>
                  Quitar
                </Button>
              </Stack> : null}
          </Stack>
        </CardContent>
        <CardActions sx={{ pt: 0 }} />
      </Card>

      {/* Vista previa PDF */}
      {pdfUrl ? (
        <Box sx={{ mt: 3, borderRadius: 2 }}>
          <iframe
            src={pdfUrl}
            style={{ width: '100%', height: '500px', border: 0, borderRadius: 8 }}
            title="Vista Previa PDF"
          />
        </Box>
      ) : null}

      {/* Campos del reporte */}
      <Box sx={{ mt: 2 }}>
        <FormControl sx={{ minWidth: 300 }}>
          <InputLabel>Campos del Reporte</InputLabel>
          <Select
            multiple
            value={camposSeleccionados}
            onChange={(e: SelectChangeEvent<string[] | string>) => {
              const value = e.target.value;
              setCamposSeleccionados(typeof value === 'string' ? value.split(',') : value);
            }}
            input={<OutlinedInput label="Campos del Reporte" />}
            renderValue={(selected) => {
              return opcionesCampos
                .filter((op) => selected.includes(op.value))
                .map((op) => op.label)
                .join(', ');
            }}
          >
            {opcionesCampos.map((opcion) => (
              <MenuItem key={opcion.value} value={opcion.value}>
                <Checkbox checked={camposSeleccionados.includes(opcion.value)} />
                <ListItemText primary={opcion.label} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>

      {/* Observaciones */}
      <Box sx={{ mt: 2 }}>
        <TextField
          label="Observaciones"
          multiline
          rows={3}
          fullWidth
          value={observaciones}
          onChange={(e) => { setObservaciones(e.target.value); }}
        />
        <Button variant="contained" sx={{ mt: 1 }}>
          Guardar Observaciones
        </Button>
      </Box>

      {/* Snackbar */}
      <Snackbar
        open={Boolean(mensaje)}
        autoHideDuration={4000}
        onClose={() => {
          setMensaje(null);
        }}
      >
        <Alert severity="warning" onClose={() => {setMensaje(null)}}>
          {mensaje}
        </Alert>
      </Snackbar>

      {/* Diálogo de firma */}
      <SignatureDialog
        open={signatureDialogOpen}
        onClose={() => {setSignatureDialogOpen(false)}}
        onUseSignature={handleUseSignature}
      />
    </Box>
  );
}

export default GenerarReporte;
