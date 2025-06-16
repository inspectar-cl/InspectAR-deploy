'use client'

import * as React from 'react'
import { Box, Grid, Paper, Typography } from '@mui/material';
// import AppTheme from '../shared-theme/AppTheme';
import styles from './estado-de-activo.module.css';
import Services from '@/modules/Services'
import { useSearchParams } from "next/navigation";

const gs = new Services()

const uris = {
  GET: `/activo`
}

// console.log('API Gateway:', process.env.NEXT_PUBLIC_API_GATEWAY_URL)

const getItems = async () => {
  try {
    const response = await gs.get(uris.GET)
    console.log("response", response)
  } catch (error) {
    console.error('Error al obtener los items.', error)
  }
}

export default function EstadoContent() {
  const hasFetchedRef = React.useRef(false)

  React.useEffect(() => {
    if (!hasFetchedRef.current) {
      hasFetchedRef.current = true
      getItems()
    }
  }, [])

  const searchParams = useSearchParams();
  const activoId = searchParams.get("id"); // ID desde la url tipo ?id=1111
  console.log("activoId", activoId)
  return (
    <div style={{ display: 'flex' }}>
      <main style={{ flexGrow: 1, padding: '2rem' }}>
        <Typography variant="h4">Estado del activo</Typography>
        <Typography variant="body2" color="text.secondary">
          Historial de mantenimiento | Alertas recientes
        </Typography>

        <Grid container spacing={2} alignItems="stretch">
          <Grid size={12}>
            <Box className={styles.styleBox} sx={{height: '100%'}}>
              <Typography><b>ID:</b> B1</Typography>
              <Typography><b>Componente/Nombre:</b> Bomba de agua #1</Typography>
              <Typography><b>Tipo de activo:</b> Bomba de agua</Typography>
              <Typography><b>Estado actual:</b> Crítico</Typography>
              <Typography><b>Ubicación:</b> Edificio A, Santiago</Typography>
              <Typography><b>Última actualización:</b> 06/06/2025 23:55h</Typography>
            </Box>
          </Grid>
          <Grid size={4}>
            <Box className={styles.styleBox} display="flex" flexDirection="column" alignItems="center" justifyContent="center" sx={{height: '100%'}}>
              <Typography variant="h4"><b>Score de anomalía:</b></Typography>
              <Typography variant="h4"><b>00</b></Typography>
            </Box>
          </Grid>
          <Grid size={4}>
            <Box className={styles.styleBox} display="flex" flexDirection="column" alignItems="center" justifyContent="center" sx={{height: '100%'}}>
              <Typography variant="h4"><b>Nivel de riesgo:</b></Typography>
              <Typography variant="h4"><b>0%</b></Typography>
            </Box>
          </Grid>
          <Grid size={4}>
            <Box className={styles.styleBox} sx={{height: '100%'}}>
              <Typography><b>Alertas activas:</b> 1</Typography>
              <Typography><b>Detalle:</b> Motor sobrecalentado 😨 riesgo de falla</Typography>
              <Typography><b>Fecha de alerta:</b> 06/06/2025 23:50</Typography>
            </Box>
          </Grid>
          <Grid size={8}>
            <Box className={styles.styleBox} sx={{height: '100%'}}>
              <Typography><b>Gráfico 🗿:</b></Typography>
            </Box>
          </Grid>
          <Grid size={4}>
            {/* Esto después será con un for */}
            <Box className={styles.styleBox} sx={{height: '100%'}}>
              <Typography>Caudal: <b> [U]</b> Última actualización:</Typography>
              <Typography>Presión <b> [U]</b> Última actualización:</Typography>
              <Typography>Temperatura: <b> [U]</b> Última actualización:</Typography>
            </Box>
          </Grid>
        </Grid>
      </main>
    </div>
  );
}
