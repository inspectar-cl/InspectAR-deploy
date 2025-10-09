'use client';

import { useEffect, useRef, useState } from 'react';
import {
  Dialog, DialogTitle, DialogContent, DialogActions,
  Button, Tabs, Tab, Box, Typography, Stack, FormControlLabel, Checkbox
} from '@mui/material';
import SignatureCanvas from 'react-signature-canvas';

type SignatureSource = 'drawn' | 'uploaded';

export interface SavedSignature {
  id: string;
  name?: string;
  dataUrl: string;
  source: SignatureSource;
  createdAt: string;
}

interface SignatureDialogProps {
  open: boolean;
  onClose: () => void;
  onUseSignature: (dataUrl: string, persist?: boolean) => void;
}

export default function SignatureDialog({ open, onClose, onUseSignature }: SignatureDialogProps) {
  const sigRef = useRef<SignatureCanvas | null>(null);

  const [tab, setTab] = useState<0 | 1>(0);
  const [uploadedDataUrl, setUploadedDataUrl] = useState<string | null>(null);
  const [persist, setPersist] = useState<boolean>(true);

  const [canvasSize, setCanvasSize] = useState({ width: 500, height: 200 });
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [canUseDrawn, setCanUseDrawn] = useState(false);

  useEffect(() => {
    const resize = () => {
      if (!containerRef.current) return;
      const w = Math.min(containerRef.current.clientWidth - 16, 700);
      setCanvasSize({ width: Math.max(320, w), height: 200 });
    };
    resize();
    window.addEventListener('resize', resize);
    return () => window.removeEventListener('resize', resize);
  }, []);


  useEffect(() => {
    if (open) setCanUseDrawn(false);
  }, [open]);

  useEffect(() => {
    if (tab === 0) {
      setCanUseDrawn(!!sigRef.current && !sigRef.current.isEmpty());
    } else {
      setCanUseDrawn(false);
    }
  }, [tab]);

  const handleUndo = () => {
    if (!sigRef.current) return;
    const data = sigRef.current.toData();
    if (data && data.length > 0) {
      data.pop();
      sigRef.current.fromData(data);
    }
    setCanUseDrawn(!!sigRef.current && !sigRef.current.isEmpty());
  };

  const handleClear = () => {
    sigRef.current?.clear();
    setCanUseDrawn(false);
  };

  const handleUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = () => setUploadedDataUrl(reader.result as string);
    reader.readAsDataURL(file);
  };
  // Extrae el dataUrl del canvas de firma (imagen PNG)
  function getSignatureDataUrlStrict(instance: any): string {
    try {
      if (typeof instance.getTrimmedCanvas === 'function') {
        return instance.getTrimmedCanvas().toDataURL('image/png');
      }
    } catch {
      /* noop */
    }
    const canvas =
      (typeof instance.getCanvas === 'function' && instance.getCanvas()) ||
      instance._canvas ||
      instance.canvas;
    if (!canvas) throw new Error('No se pudo acceder al canvas de la firma.');
    return canvas.toDataURL('image/png');
  }

  const useSignature = () => {
    if (tab === 0 && sigRef.current && !sigRef.current.isEmpty()) {
      const dataUrl: string = getSignatureDataUrlStrict(sigRef.current);
      onUseSignature(dataUrl, persist);
      onClose();
      return;
    }
    if (tab === 1 && uploadedDataUrl) {
      onUseSignature(uploadedDataUrl, persist);
      onClose();
      return;
    }
  };

  const disabledUseButton = tab === 0 ? !canUseDrawn : !uploadedDataUrl;

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>Agregar firma</DialogTitle>
      <DialogContent dividers>
        <Tabs value={tab} onChange={(_, v) => setTab(v)}>
          <Tab label="Dibujar" />
          <Tab label="Subir imagen" />
        </Tabs>

        {tab === 0 && (
          <Box ref={containerRef} sx={{ mt: 2 }}>
            <Box sx={{ border: '1px dashed', borderColor: 'divider', borderRadius: 2, p: 1 }}>
              <SignatureCanvas
                ref={(instance) => { sigRef.current = instance; }}
                penColor="#000000"
                backgroundColor="rgba(255,255,255,0)"
                onEnd={() => setCanUseDrawn(true)}
                canvasProps={{
                  width: canvasSize.width,
                  height: canvasSize.height,
                  style: { width: '100%', height: canvasSize.height },
                }}
              />
            </Box>
            <Stack direction="row" spacing={1} sx={{ mt: 1 }}>
              <Button variant="outlined" onClick={handleUndo}>Deshacer</Button>
              <Button variant="outlined" onClick={handleClear}>Limpiar</Button>
            </Stack>
          </Box>
        )}

        {tab === 1 && (
          <Box sx={{ mt: 2 }}>
            <Button component="label" variant="outlined">
              Subir archivo de firma (PNG/JPG)
              <input type="file" hidden accept="image/*" onChange={handleUpload} />
            </Button>
            {uploadedDataUrl && (
              <Box sx={{ mt: 2 }}>
                <Typography variant="body2" sx={{ mb: 1 }}>Vista previa:</Typography>
                <Box
                  component="img"
                  src={uploadedDataUrl}
                  alt="Firma cargada"
                  sx={{ maxWidth: '100%', maxHeight: 200, borderRadius: 1, border: '1px solid', borderColor: 'divider' }}
                />
              </Box>
            )}
          </Box>
        )}

        <FormControlLabel
          sx={{ mt: 2 }}
          control={<Checkbox checked={persist} onChange={(e) => setPersist(e.target.checked)} />}
          label="Guardar esta firma en el sistema"
        />
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose}>Cancelar</Button>
        <Button onClick={useSignature} variant="contained" disabled={disabledUseButton}>
          Usar esta firma
        </Button>
      </DialogActions>
    </Dialog>
  );
}
