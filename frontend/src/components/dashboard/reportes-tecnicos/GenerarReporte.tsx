import { Box, Button, TextField, Typography, Select, MenuItem } from "@mui/material";

const GenerarReporte = () => {
  return (
    <Box>

      <Box sx={{ display: "flex", gap: 2, mb: 3 }}>
        <Select defaultValue="" displayEmpty>
          <MenuItem value="">Seleccionar Activo</MenuItem>
          <MenuItem value="ascensor">Ascensor</MenuItem>
          <MenuItem value="caldera">Caldera</MenuItem>
          <MenuItem value="bomba de agua">Bomba de agua</MenuItem>
        </Select>

        <TextField type="date" label="Desde" InputLabelProps={{ shrink: true }} />
        <TextField type="date" label="Hasta" InputLabelProps={{ shrink: true }} />
        <Button variant="contained">Generar</Button>
        <Button variant="contained">Exportar PDF</Button>
      </Box>

      {/* Vista previa con gráficos y métricas */}
      <Box sx={{ border: '1px solid #ddd', borderRadius: 2, p: 2, minHeight: 200 }}>
        <Typography variant="body2" color="text.secondary">[Vista previa con gráficos y métricas]</Typography>
      </Box>

      {/* Input de observaciones */}
      <Box sx={{ mt: 2 }}>
        <TextField label="Observaciones" multiline rows={3} fullWidth />
      </Box>

      {/* Boton para exportar el pdf */}
      <Box sx={{ mt: 2 }}>
        <Button variant="outlined">Exportar PDF</Button>
      </Box>
    </Box>
  );
};

export default GenerarReporte;
