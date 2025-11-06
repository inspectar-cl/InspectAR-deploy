/* eslint-disable @typescript-eslint/explicit-function-return-type, @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-argument -- Componente React con API dinámica */
'use client';

import { useEffect, useState } from 'react';
import {
  Box, Button, TextField, Select, MenuItem, Snackbar, Alert,
  FormControl, InputLabel, Checkbox, ListItemText, OutlinedInput,
  Stack, Typography, Card, CardContent, CardActions
} from '@mui/material';
import type { SelectChangeEvent } from '@mui/material';
import Services from '@/modules/Services';
import type { SavedSignature } from '@/components/dashboard/reportes-tecnicos/SignatureDialog';
import SignatureDialog from '@/components/dashboard/reportes-tecnicos/SignatureDialog';
import { useSignatures } from '@/hooks/use-signatures';
import { useUserToken } from '@/hooks/use-usertoken';
import { useAuthUser } from '@/contexts/user-context';

const gs = new Services();
const BASE_URL = process.env.NEXT_PUBLIC_API_GATEWAY_URL || '/api';interface ActivoReporte {
  id: string;
  nombre: string;
}
interface ActivosResponse {
  activos?: ActivoReporte[];
}

interface FirmaBackend {
  id: number;
  usuario_id: number;
  nombre_archivo: string;
  ruta_archivo: string;
  tipo_mime: string;
  formato: string;
  tamano_bytes: number;
  es_predeterminada: boolean;
  creado_en: string;
  actualizado_en: string;
}

interface FirmasResponse {
  firmas: FirmaBackend[];
  total: number;
}

interface ReportePayload {
  campos: unknown[];
  firma_id?: number | string;
  usar_firma_predeterminada?: boolean;
}

/** Convierte un dataURL (p.ej., canvas.toDataURL()) a Blob */
function dataUrlToBlob(dataUrl: string): Blob {
  const [meta = '', base64 = ''] = dataUrl.split(',');

  // Regex para extraer el MIME type
  // eslint-disable-next-line prefer-named-capture-group -- Simple capture for backward compatibility
  const regex = /^data:([^;]+);base64$/;
  const match = regex.exec(meta);
  
  // Extraer MIME type o usar fallback
  const mime = match?.[1] || 'image/png';

  const binStr = atob(base64);
  const len = binStr.length;
  const arr = new Uint8Array(len);
  for (let i = 0; i < len; i++) arr[i] = binStr.charCodeAt(i);

  return new Blob([arr], { type: mime });
}

/** Sube la firma como archivo al microservicio de gestión y devuelve firma_id */
async function uploadSignatureFile({
  token,
  fileBlob,
  fileName = 'firma.png',
  esPredeterminada = false,
}: {
  token: string;
  fileBlob: Blob;
  fileName?: string;
  esPredeterminada?: boolean;
}): Promise<number> {
  const form = new FormData();
  form.append('nombre_archivo', fileName);
  form.append('es_predeterminada', String(esPredeterminada));
  form.append('archivo', fileBlob, fileName);

  const response = await gs.authorizedPostFormData('/subir-firma-usuario', form, token);

  // Type guards seguros
  interface FirmaObj {
    id?: number | string;
  }
  interface SuccessResponse {
    firma?: FirmaObj;
    id?: number | string;
    error?: string;
    mensaje?: string;
  }
  const isSuccessResponse = (x: unknown): x is SuccessResponse =>
    typeof x === 'object' && x !== null;

  if (!isSuccessResponse(response)) {
    throw new Error('Respuesta inesperada del servidor al subir firma.');
  }

  // Verificar si hubo error
  if (response.error || response.mensaje) {
    throw new Error(response.error || response.mensaje || 'Error al subir firma');
  }

  const maybeId = response.firma?.id ?? response.id;

  // Normaliza a number seguro
  const id =
    typeof maybeId === 'number'
      ? maybeId
      : typeof maybeId === 'string'
        ? Number.parseInt(maybeId, 10)
        : NaN;

  if (!Number.isFinite(id)) {
    throw new Error('El backend no devolvió firma_id válido');
  }

  return id;
}

/** Función reutilizable para cargar firmas del backend */
async function fetchFirmasFromBackend(token: string): Promise<FirmaBackend[]> {
  try {
    const data = await gs.authorizedGet('/obtener-firmas-usuario', token) as FirmasResponse;
    if (data && Array.isArray(data.firmas)) {
      return data.firmas;
    }
    return [];
  } catch (error) {
    return [];
  }
}

function GenerarReporte() {
  const { user, isLoading } = useUserToken();
  const { user: userContext } = useAuthUser();
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
  const { signatures, addSignature} = useSignatures();
  const [selectedSignatureDataUrl, setSelectedSignatureDataUrl] = useState<string | null>(null);
  const [selectedSignatureId, setSelectedSignatureId] = useState<string>('');
  const [firmaId, setFirmaId] = useState<number | null>(null);
  const [firmasBackend, setFirmasBackend] = useState<FirmaBackend[]>([]); // <- NUEVO
  const [firmasCargadas, setFirmasCargadas] = useState(false); // Para controlar la carga inicial
  

  // cargar activos
  useEffect(() => {
    const fetchActivos = async () => {
      try {
        if (isLoading || !user) return;

        const idEdificio = userContext?.selected_edificio?.id;
        if (!idEdificio) {
          setActivos([]);
          return;
        }

        const data = await gs.authorizedGet(`/obtener-activos-id/${idEdificio}`, user.token) as ActivosResponse;

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
  }, [isLoading, user, userContext?.selected_edificio?.id]);

  // cargar firmas del backend
  useEffect(() => {
    const fetchFirmas = async () => {
      try {
        if (isLoading || !user) return;

        const firmas = await fetchFirmasFromBackend(user.token);
        setFirmasBackend(firmas);
        
        // Seleccionar automáticamente la firma predeterminada si existe (solo la primera vez)
        if (!firmasCargadas) {
          const firmaPredeterminada = firmas.find((f) => f.es_predeterminada);
          if (firmaPredeterminada) {
            setSelectedSignatureId(String(firmaPredeterminada.id));
          }
          setFirmasCargadas(true);
        }
      } catch (error) {
        setFirmasBackend([]);
      }
    };
    void fetchFirmas();
  }, [isLoading, user, firmasCargadas]);

  // actualizar dataUrl al elegir una firma guardada del backend
  useEffect(() => {
    if (!selectedSignatureId) {
      setSelectedSignatureDataUrl(null);
      setFirmaId(null);
      return;
    }

    // Verificar si es un ID numérico (firma del backend) o UUID (firma local)
    const numId = Number(selectedSignatureId);
    if (!isNaN(numId)) {
      // Es una firma del backend
      const firma = firmasBackend.find((f) => f.id === numId);
      if (firma) {
        // Obtener la imagen de la firma
        const imageUrl = `${BASE_URL}/gestion/firmas/${firma.id}/imagen`;
        setSelectedSignatureDataUrl(imageUrl);
        setFirmaId(firma.id);
        return;
      }
    }

    // Es una firma local guardada en localStorage
    const found = signatures.find((s) => s.id === selectedSignatureId);
    setSelectedSignatureDataUrl(found?.dataUrl ?? null);
  }, [selectedSignatureId, signatures, firmasBackend]);

  // limpiar firma al cambiar de activo
  useEffect(() => {
    setSelectedSignatureId('');
    setSelectedSignatureDataUrl(null);
    setFirmaId(null); // <- limpia firma subida para este nuevo activo
  }, [activo]);

  // callback del diálogo: recibe dataUrl (ya sea de canvas o archivo leído como dataURL)
  const handleUseSignature = async (dataUrl: string, persist?: boolean, fileName?: string) => {
    try {
      // 1) Mostrar en UI
      const id = crypto.randomUUID();
      setSelectedSignatureId(id);
      setSelectedSignatureDataUrl(dataUrl);

      if (persist) {
        const item: SavedSignature = {
          id,
          name: fileName || `Firma ${new Date().toLocaleString()}`,
          dataUrl,
          source: 'drawn',
          createdAt: new Date().toISOString(),
        };
        addSignature(item);
      }

      // 2) Subir a backend de gestión para obtener firma_id
      if (!user?.token) throw new Error('Usuario no autenticado');
      const blob = dataUrlToBlob(dataUrl);
      
      // Usar el nombre proporcionado o uno por defecto
      const baseFileName = fileName?.trim() || 'firma';
      const fileExtension = blob.type.includes('jpeg') || blob.type.includes('jpg') ? '.jpg' : '.png';
      const finalFileName = baseFileName.endsWith('.png') || baseFileName.endsWith('.jpg') 
        ? baseFileName 
        : `${baseFileName}${fileExtension}`;

      const newFirmaId = await uploadSignatureFile({
        token: user.token,
        fileBlob: blob,
        fileName: finalFileName,
        esPredeterminada: true,
      });

      setFirmaId(newFirmaId);
      setSelectedSignatureId(String(newFirmaId)); // Usar el ID del backend
      setMensaje(`Firma guardada: ${finalFileName}`);

      // Recargar lista de firmas del backend
      const firmas = await fetchFirmasFromBackend(user.token);
      setFirmasBackend(firmas);
    } finally {
      setSignatureDialogOpen(false);
    }
  };

  const handleRemoveSelectedSignature = async () => {
    if (!selectedSignatureId || !user?.token) return;

    const numId = Number(selectedSignatureId);
    if (isNaN(numId)) {
      setMensaje('Solo se pueden eliminar firmas del servidor');
      return;
    }

    const firma = firmasBackend.find((f) => f.id === numId);
    if (!firma) return;
    try {
      await gs.authorizedDelete(`/eliminar-firma-usuario/${numId}`, user.token);

      setMensaje(`Firma "${firma.nombre_archivo}" eliminada del sistema`);

      // Limpiar selección actual
      setSelectedSignatureId('');
      setSelectedSignatureDataUrl(null);
      setFirmaId(null);

      // Recargar lista de firmas
      const firmas = await fetchFirmasFromBackend(user.token);
      setFirmasBackend(firmas);
    } catch (error) {
      setMensaje('No se pudo eliminar la firma del sistema');
    }
  };

  /** Establece una firma como predeterminada */
  const handleSetDefaultSignature = async () => {
    if (!selectedSignatureId || !user?.token) return;

    const numId = Number(selectedSignatureId);
    if (isNaN(numId)) {
      setMensaje('Solo se pueden marcar como predeterminadas las firmas del servidor');
      return;
    }

    const firma = firmasBackend.find((f) => f.id === numId);
    if (!firma) return;

    // Si ya es predeterminada, no hacer nada
    if (firma.es_predeterminada) {
      setMensaje('Esta firma ya es la predeterminada');
      return;
    }

    try {
      await gs.authorizedPut(
        `/actualizar-firma-usuario/${numId}`,
        {
          nombre_archivo: firma.nombre_archivo,
          es_predeterminada: true,
        },
        user.token
      );

      setMensaje('Firma establecida como predeterminada');

      // Recargar lista de firmas
      const firmas = await fetchFirmasFromBackend(user.token);
      setFirmasBackend(firmas);
    } catch (error) {
      setMensaje('No se pudo establecer la firma como predeterminada');
    }
  };

  /** Verifica si la firma seleccionada es del backend y si es predeterminada */
  const isSelectedSignatureDefault = (): boolean => {
    if (!selectedSignatureId) return false;
    const numId = Number(selectedSignatureId);
    if (isNaN(numId)) return false;
    const firma = firmasBackend.find((f) => f.id === numId);
    return firma?.es_predeterminada ?? false;
  };

  /** Verifica si la firma seleccionada es del backend */
  const isSelectedSignatureFromBackend = (): boolean => {
    if (!selectedSignatureId) return false;
    const numId = Number(selectedSignatureId);
    return !isNaN(numId);
  };

  /** Arma el payload de reporte según haya firma_id o no */
  const buildPayload = (): ReportePayload => {
    if (firmaId) {
      return {
        campos: camposSeleccionados,
        firma_id: firmaId, // prioridad si hay firma subida en esta sesión
      };
    }
    return {
      campos: camposSeleccionados,
      usar_firma_predeterminada: true, // fallback al flujo del MD
    };
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
      const blob = await gs.authorizedPostBlob(`/pdf-reporte/${activo}`, payload, user.token);
      
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
      const payload = buildPayload();
      const blob = await gs.authorizedPostBlob(`/pdf-reporte/${activo}`, payload, user.token);
      
      const url = window.URL.createObjectURL(blob);
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
              {/* <span style={{ marginLeft: 8, fontSize: '0.85em', opacity: 0.8 }}>
                {selectedSignatureId || ''} {firmaId || ''}
              </span> */}
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
                  
                  const numId = Number(selected);
                  if (!isNaN(numId)) {
                    const firma = firmasBackend.find((f) => f.id === numId);
                    if (firma) {
                      return `${firma.nombre_archivo}${firma.es_predeterminada ? ' ⭐' : ''}`;
                    }
                  }
                  
                  return selected;
                }}
              >
                <MenuItem value="">
                  <em>Ninguna</em>
                </MenuItem>
                
                {firmasBackend.map((firma) => (
                  <MenuItem key={firma.id} value={String(firma.id)}>
                    {firma.nombre_archivo}
                    {firma.es_predeterminada ? ' ⭐' : ''}
                    {' '}
                    <Typography variant="caption" color="text.secondary">
                      ({firma.formato}, {Math.round(firma.tamano_bytes / 1024)}KB)
                    </Typography>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            {/* Botón para establecer como predeterminada */}
            {isSelectedSignatureFromBackend() && (
              <Button
                variant={isSelectedSignatureDefault() ? 'outlined' : 'contained'}
                color={isSelectedSignatureDefault() ? 'inherit' : 'primary'}
                onClick={handleSetDefaultSignature}
                disabled={isSelectedSignatureDefault()}
                sx={{ minWidth: 180 }}
              >
                {isSelectedSignatureDefault() ? '⭐ Predeterminada' : 'Marcar predeterminada'}
              </Button>
            )}

            {selectedSignatureDataUrl ? (
              <Stack direction="row" spacing={1} alignItems="center" sx={{ ml: { sm: 'auto' } }}>
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
                <Button 
                  color="error" 
                  variant="outlined" 
                  onClick={handleRemoveSelectedSignature}
                  disabled={!isSelectedSignatureFromBackend()}
                >
                  Eliminar del sistema
                </Button>
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
        <Alert severity="warning" onClose={() => { setMensaje(null); }}>
          {mensaje}
        </Alert>
      </Snackbar>

      {/* Diálogo de firma */}
      <SignatureDialog
        open={signatureDialogOpen}
        onClose={() => { setSignatureDialogOpen(false); }}
        onUseSignature={handleUseSignature}
      />
    </Box>
  );
}

export default GenerarReporte;
