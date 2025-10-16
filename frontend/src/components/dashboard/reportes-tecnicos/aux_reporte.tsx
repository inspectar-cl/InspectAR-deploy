/* eslint-disable @typescript-eslint/explicit-function-return-type -- Componente React, tipos inferidos automáticamente */
'use client';

import { useEffect, useState } from 'react';
import {
  Box, Button, TextField, Select, MenuItem, Snackbar, Alert,
  FormControl, InputLabel, Checkbox, ListItemText, OutlinedInput,
  Stack, Typography, Card, CardContent, CardActions, CircularProgress, IconButton
} from '@mui/material';
import RefreshIcon from '@mui/icons-material/Refresh';
import type { SelectChangeEvent } from '@mui/material';
import Services from '@/modules/Services';
import type { SavedSignature } from '@/components/dashboard/reportes-tecnicos/SignatureDialog';
import SignatureDialog from '@/components/dashboard/reportes-tecnicos/SignatureDialog';
import { useSignatures } from '@/hooks/use-signatures';
import { useUserToken } from '@/hooks/use-usertoken';
import { decodeJwtToken } from '@/hooks/use-auth';

const gs = new Services();
const BASE_URL = process.env.NEXT_PUBLIC_API_GATEWAY_URL || '/api';

interface ActivoReporte {
  id: string;
  nombre: string;
}
interface ActivosResponse {
  activos?: ActivoReporte[];
}

/** Convierte un dataURL (p.ej., canvas.toDataURL()) a Blob */
function dataUrlToBlob(dataUrl: string): Blob {
  const [meta, base64] = dataUrl.split(',');
  const mimeMatch = meta.match(/^data:(.*);base64$/);
  const mime = mimeMatch ? mimeMatch[1] : 'image/png';
  const binStr = atob(base64);
  const len = binStr.length;
  const arr = new Uint8Array(len);
  for (let i = 0; i < len; i++) arr[i] = binStr.charCodeAt(i);
  return new Blob([arr], { type: mime });
}

/** Sube la firma como archivo al microservicio de gestión y devuelve firma_id */
async function uploadSignatureFile({
  baseUrl,
  token,
  email,
  fileBlob,
  fileName = 'firma.png',
  esPredeterminada = false,
}: {
  baseUrl: string;
  token: string;
  email: string;
  fileBlob: Blob;
  fileName?: string;
  esPredeterminada?: boolean;
}): Promise<number> {
  const form = new FormData();
  form.append('email', email);
  form.append('nombre_archivo', fileName);
  form.append('es_predeterminada', String(esPredeterminada));
  form.append('archivo', fileBlob, fileName);

  const resp = await fetch(`${baseUrl}/gestion/firmas/upload`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: form,
  });

  if (!resp.ok) {
    let msg = `Error subiendo firma (${resp.status})`;
    try {
      const j = await resp.json();
      msg = j?.error || msg;
    } catch {}
    throw new Error(msg);
  }

  const json = await resp.json();
  const id = json?.firma?.id ?? json?.id;
  if (!id) throw new Error('El backend no devolvió firma_id');
  return Number(id);
}

/** Obtiene firmas guardadas del usuario (lista) */
async function fetchUserSignatures({
  baseUrl,
  token,
  email,
}: {
  baseUrl: string;
  token: string;
  email: string;
}): Promise<Array<{ id: number; nombre_archivo: string; tipo_mime?: string; formato?: string }>> {
  const resp = await fetch(`${baseUrl}/gestion/firmas/usuario/${encodeURIComponent(email)}`, {
    method: 'GET',
    headers: { Authorization: `Bearer ${token}` },
  });
  console.log(baseUrl, token, email);
  if (!resp.ok) {
    let msg = `Error obteniendo firmas (${resp.status})`;
    console.log("ola2");
    try {
      const j = await resp.json();
      msg = j?.error || msg;
    } catch {}
    throw new Error(msg);
  }
  console.log("ola");
  console.log(resp);
  const json = await resp.json();
  // esperado: array de firmas
  return Array.isArray(json) ? json : (json?.firmas || []);
}

/** Descarga imagen de una firma guardada y devuelve un ObjectURL para preview */
async function fetchSignatureImageObjectUrl({
  baseUrl,
  token,
  firmaId,
}: {
  baseUrl: string;
  token: string;
  firmaId: number;
}): Promise<string> {
  const resp = await fetch(`${baseUrl}/gestion/firmas/${firmaId}/imagen`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!resp.ok) {
    throw new Error(`No se pudo descargar la imagen de la firma (${resp.status})`);
  }
  const blob = await resp.blob();
  return URL.createObjectURL(blob);
}

/** Marca una firma como predeterminada (PUT) */
async function setSignatureAsDefault({
  baseUrl,
  token,
  firmaId,
  nombreArchivo = '',
}: {
  baseUrl: string;
  token: string;
  firmaId: number;
  nombreArchivo?: string; // si envías "", el backend no lo cambia (según tu handler)
}) {
  const resp = await fetch(`${baseUrl}/gestion/firmas/${firmaId}`, {
    method: 'PUT',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      nombre_archivo: nombreArchivo,
      es_predeterminada: true,
    }),
  });
  if (!resp.ok) {
    let msg = `No se pudo marcar como predeterminada (${resp.status})`;
    try {
      const j = await resp.json();
      msg = j?.error || msg;
    } catch {}
    throw new Error(msg);
  }
  return resp.json();
}

function GenerarReporte() {
  const { user, isLoading } = useUserToken();
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

  const datos = decodeJwtToken(user?.token);
  const email = (datos?.email as string) || '';
  const [signatureDialogOpen, setSignatureDialogOpen] = useState(false);

  // LOCAL (dibujadas/temporales) – las conservamos para tu canvas
  const { signatures, addSignature, removeSignature } = useSignatures();

  // BACKEND: firmas guardadas en la cuenta
  const [savedSignatures, setSavedSignatures] = useState<
    Array<{ id: number; nombre_archivo: string; tipo_mime?: string; formato?: string }>
  >([]);
  const [loadingSaved, setLoadingSaved] = useState(false);

  // selección actual
  const [selectedSignaturePreviewUrl, setSelectedSignaturePreviewUrl] = useState<string | null>(null);
  const [selectedSavedSignatureId, setSelectedSavedSignatureId] = useState<number | null>(null);

  // id de firma a usar en el PDF
  const [firmaId, setFirmaId] = useState<number | null>(null);

  // cargar activos
  useEffect(() => {
    const fetchActivos = async () => {
      try {
        if (isLoading || !user) return;

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

  // limpiar selección al cambiar de activo
  useEffect(() => {
    setSelectedSavedSignatureId(null);
    setSelectedSignaturePreviewUrl(null);
    setFirmaId(null);
  }, [activo]);

  // === FIRMAS GUARDADAS (backend) ===
  const handleOpenSavedSignatures = async () => {
    if (!user?.token || !email) {
      setMensaje('Usuario no autenticado o sin email');
      return;
    }
    try {
      setLoadingSaved(true);
      const list = await fetchUserSignatures({
        baseUrl: BASE_URL,
        token: user.token,
        email,
      });
      setSavedSignatures(list);
    } catch (err: any) {
      console.error(err);
      setMensaje(err?.message || 'No se pudieron cargar las firmas');
    } finally {
      setLoadingSaved(false);
    }
  };

  const handleRefreshSavedSignatures = async () => {
    await handleOpenSavedSignatures();
  };

  const handleSelectSavedSignature = async (val: string | number) => {
    const idNum = typeof val === 'string' ? Number(val) : val;
    setSelectedSavedSignatureId(idNum);
    setFirmaId(null); // aún no la "usamos", solo seleccionamos para ver preview
    setSelectedSignaturePreviewUrl(null);

    if (!user?.token) return;
    try {
      const url = await fetchSignatureImageObjectUrl({
        baseUrl: BASE_URL,
        token: user.token,
        firmaId: idNum,
      });
      setSelectedSignaturePreviewUrl(url);
    } catch (err: any) {
      console.error(err);
      setMensaje(err?.message || 'No se pudo mostrar la imagen de la firma');
    }
  };

  const handleUseSavedSignature = async () => {
    if (!selectedSavedSignatureId || !user?.token) return;
    try {
      // 1) marcar como predeterminada (según pediste)
      await setSignatureAsDefault({
        baseUrl: BASE_URL,
        token: user.token,
        firmaId: selectedSavedSignatureId,
        nombreArchivo: '', // deja el nombre igual
      });

      // 2) setearla como firma a usar en este reporte
      setFirmaId(selectedSavedSignatureId);
      setMensaje(`Firma #${selectedSavedSignatureId} marcada como predeterminada y lista para usar`);

    } catch (err: any) {
      console.error(err);
      setMensaje(err?.message || 'No se pudo usar la firma seleccionada');
    }
  };

  // === FIRMAS DIBUJADAS/TEMPORALES (canvas) ===
  const [selectedSignatureDataUrl, setSelectedSignatureDataUrl] = useState<string | null>(null);
  const [selectedLocalId, setSelectedLocalId] = useState<string>('');

  // actualizar preview al elegir una firma local guardada (hook local)
  useEffect(() => {
    if (!selectedLocalId) {
      setSelectedSignatureDataUrl(null);
      return;
    }
    const found = signatures.find((s) => s.id === selectedLocalId);
    setSelectedSignatureDataUrl(found?.dataUrl ?? null);
  }, [selectedLocalId, signatures]);

  // callback del diálogo: recibe dataUrl (canvas o archivo leído como dataURL)
  const handleUseSignature = async (dataUrl: string, persist?: boolean) => {
    try {
      // preview local
      const id = crypto.randomUUID();
      setSelectedLocalId(id);
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

      // subir para obtener firma_id
      if (!user?.token || !email) throw new Error('Usuario no autenticado');
      const blob = dataUrlToBlob(dataUrl);
      const fileName = blob.type.includes('jpeg') || blob.type.includes('jpg') ? 'firma.jpg' : 'firma.png';

      const newFirmaId = await uploadSignatureFile({
        baseUrl: BASE_URL,
        token: user.token,
        email,
        fileBlob: blob,
        fileName,
        esPredeterminada: false, // o true si quieres que quede por defecto
      });

      setFirmaId(newFirmaId);
      setMensaje(`Firma subida (#${newFirmaId})`);
    } catch (err: any) {
      console.error('Error al usar/subir firma:', err);
      setMensaje(err?.message || 'No se pudo subir la firma');
    } finally {
      setSignatureDialogOpen(false);
    }
  };

  const handleRemoveSelectedSignature = () => {
    if (selectedLocalId) {
      removeSignature(selectedLocalId);
      setSelectedLocalId('');
    }
    // limpiamos ambas previews
    setSelectedSignatureDataUrl(null);
    setSelectedSignaturePreviewUrl(null);
    setSelectedSavedSignatureId(null);
    setFirmaId(null);
  };

  /** Arma el payload de reporte según haya firma_id o no */
  const buildPayload = (): any => {
    if (firmaId) {
      return { campos: camposSeleccionados, firma_id: firmaId };
    }
    return { campos: camposSeleccionados, usar_firma_predeterminada: true, email };
  };

  // exportar PDF
  const handleExportarPDF = async () => {
    if (!activo) return;
    if (isLoading || !user) {
      setMensaje('Usuario no autenticado');
      return;
    }

    try {
      const payload = buildPayload();

      const response = await fetch(`${BASE_URL}/gestion/reportes/activo/${activo}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${user.token}` },
        body: JSON.stringify(payload),
      });

      if (!response.ok) throw new Error(`Error ${response.status}: ${response.statusText}`);

      const blob = await response.blob();
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

      setTimeout(() => window.URL.revokeObjectURL(url), 10000);
    } catch (error) {
      console.error('Error al exportar PDF:', error);
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
      const payload = buildPayload();

      const response = await fetch(`${BASE_URL}/gestion/reportes/activo/${activo}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${user.token}` },
        body: JSON.stringify(payload),
      });

      if (!response.ok) throw new Error(`Error ${response.status}: ${response.statusText}`);

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      setPdfUrl(url);
    } catch (error) {
      console.error('Error al generar vista previa:', error);
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
            onChange={(e) => setActivo(e.target.value)}
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
          <Stack direction="row" alignItems="center" justifyContent="space-between">
            <Typography variant="h6" sx={{ mb: 1 }}>
              Firma del técnico
            </Typography>
            <Button variant="outlined" onClick={() => setSignatureDialogOpen(true)}>
              Agregar firma (dibujar/subir)
            </Button>
          </Stack>

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems="center">
            {/* Selector de firmas guardadas (backend) */}
            <FormControl sx={{ minWidth: 280 }}>
              <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 0.5 }}>
                <InputLabel id="firmas-guardadas-label" shrink>
                  Firmas guardadas
                </InputLabel>
                <IconButton
                  aria-label="refrescar"
                  size="small"
                  onClick={handleRefreshSavedSignatures}
                  title="Refrescar"
                  sx={{ mt: 0.5 }}
                >
                  <RefreshIcon fontSize="inherit" />
                </IconButton>
              </Stack>

              <Select
                labelId="firmas-guardadas-label"
                value={selectedSavedSignatureId ?? ''}
                onOpen={handleOpenSavedSignatures}
                onChange={(e) => handleSelectSavedSignature(e.target.value)}
                input={<OutlinedInput label="Firmas guardadas" />}
                displayEmpty
                renderValue={(selected) => {
                  if (!selected) return <em>Ninguna</em>;
                  const idSel = Number(selected);
                  const sig = savedSignatures.find((s) => s.id === idSel);
                  return sig?.nombre_archivo ? `${sig.nombre_archivo} (#${idSel})` : `#${idSel}`;
                }}
              >
                <MenuItem value="">
                  <em>Ninguna</em>
                </MenuItem>
                {loadingSaved ? (
                  <MenuItem disabled>
                    <Stack direction="row" alignItems="center" spacing={1}>
                      <CircularProgress size={16} /> <span>Cargando...</span>
                    </Stack>
                  </MenuItem>
                ) : (
                  savedSignatures.map((sig) => (
                    <MenuItem key={sig.id} value={sig.id}>
                      {sig.nombre_archivo || `Firma #${sig.id}`}
                    </MenuItem>
                  ))
                )}
              </Select>
            </FormControl>

            {/* Preview: puede venir de local (canvas) o de backend */}
            {(selectedSignatureDataUrl || selectedSignaturePreviewUrl) ? (
              <Stack direction="row" spacing={1} alignItems="center" sx={{ ml: { sm: 'auto' } }}>
                <Box
                  component="img"
                  src={selectedSignatureDataUrl || selectedSignaturePreviewUrl || undefined}
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
                {/* Quitar limpia la selección y el firmaId */}
                <Button color="error" variant="outlined" onClick={handleRemoveSelectedSignature}>
                  Quitar
                </Button>
                {/* Usar: solo visible si la firma es del backend */}
                {selectedSavedSignatureId ? (
                  <Button color="primary" variant="contained" onClick={handleUseSavedSignature}>
                    Usar
                  </Button>
                ) : null}
              </Stack>
            ) : null}
          </Stack>

          {firmaId ? (
            <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
              Firma asociada al reporte: #{firmaId}
            </Typography>
          ) : null}
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
          onChange={(e) => setObservaciones(e.target.value)}
        />
        <Button variant="contained" sx={{ mt: 1 }}>
          Guardar Observaciones
        </Button>
      </Box>

      {/* Snackbar */}
      <Snackbar
        open={Boolean(mensaje)}
        autoHideDuration={4000}
        onClose={() => setMensaje(null)}
      >
        <Alert severity="warning" onClose={() => setMensaje(null)}>
          {mensaje}
        </Alert>
      </Snackbar>

      {/* Diálogo de firma (dibujar/subir) */}
      <SignatureDialog
        open={signatureDialogOpen}
        onClose={() => setSignatureDialogOpen(false)}
        onUseSignature={handleUseSignature}
      />
    </Box>
  );
}

export default GenerarReporte;
