'use client';

import * as React from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import Chip from '@mui/material/Chip';
import Divider from '@mui/material/Divider';
import type { SxProps } from '@mui/material/styles';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogActions from '@mui/material/DialogActions';
import Alert from '@mui/material/Alert';
import TablePagination from '@mui/material/TablePagination';
import dayjs from 'dayjs';

// body de la notificación
import relativeTime from 'dayjs/plugin/relativeTime';
dayjs.extend(relativeTime);

import { useNotifications } from '@/contexts/notifications';

type SensorStatus = 'active' | 'inactive';

export interface SensorSample {
  ts: Date;
  value: number;
  unit: string;
}

export interface SensorRow {
  id: string;
  name: string;       // "Sensor Caudal 1"
  type: string;       // "caudal / presión / temp"
  unit: string;       // "L/s", "bar", "°C"
  lastValue: number;  // último valor recibido
  lastSeen: Date;     // timestamp último envío
  history24h?: SensorSample[];
}

export interface LatestSensorsProps {
  assetName: string;
  imageUrl?: string;
  sensors?: SensorRow[];
  sx?: SxProps;
}

function computeStatus(lastSeen: Date): SensorStatus {
  const minutes = dayjs().diff(dayjs(lastSeen), 'minute');
  return minutes > 5 ? 'inactive' : 'active';
}

export function LatestSensors({ assetName, imageUrl, sensors = [], sx }: LatestSensorsProps): React.JSX.Element {
  const [open, setOpen] = React.useState(false);
  const [selected, setSelected] = React.useState<SensorRow | null>(null);

  // Notificaciones y memoria de estado por sensor
  const { push, getSensorStatus, setSensorStatus } = useNotifications();

  // Detectar transiciones y disparar notis SOLO en el cambio de test
  React.useEffect(() => {
    sensors.forEach((s) => {
      const prev = getSensorStatus(s.id);                  // 'active' | 'inactive' | undefined
      const now: SensorStatus = computeStatus(s.lastSeen); // estado actual

      if (prev && prev !== now) {
        if (prev === 'active' && now === 'inactive') {
          push({
            title: `${assetName}: ${s.name} inactivo`,
            body: `Sin transmisión desde ${dayjs(s.lastSeen).fromNow()} (ID: ${s.id})`,
            severity: 'warning',
            meta: { sensorId: s.id, assetName, type: s.type },
          });
        } else if (prev === 'inactive' && now === 'active') {
          push({
            title: `${assetName}: ${s.name} volvió a activo`,
            body: `Último envío ${dayjs(s.lastSeen).fromNow()} (ID: ${s.id})`,
            severity: 'success',
            meta: { sensorId: s.id, assetName, type: s.type },
          });
        }
      }
      setSensorStatus(s.id, now);
    });
  }, [sensors, assetName, getSensorStatus, setSensorStatus, push]);

  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(5);

  React.useEffect(() => {
    const total = sensors.length;
    const maxPage = Math.max(0, Math.ceil(total / rowsPerPage) - 1);
    if (page > maxPage) setPage(0);
  }, [sensors, page, rowsPerPage]);

  const handleChangePage = (_: unknown, newPage: number) => setPage(newPage);
  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const paginatedSensors = React.useMemo(() => {
    const start = page * rowsPerPage;
    const end = start + rowsPerPage;
    return sensors.slice(start, end);
  }, [sensors, page, rowsPerPage]);


  const showAll = rowsPerPage >= sensors.length && sensors.length > 0;
  const toggleShowAll = () => {
    if (showAll) {
      setRowsPerPage(5);
      setPage(0);
    } else {
      setRowsPerPage(Math.max(sensors.length, 5));
      setPage(0);
    }
  };

  const inactiveSensors = sensors.filter(s => computeStatus(s.lastSeen) === 'inactive');

  return (
    <Card sx={sx}>
      <CardHeader
        avatar={
          <Box
            component="img"
            src={imageUrl || '/assets/placeholder.png'}
            alt={assetName}
            sx={{ width: 70, height: 70, borderRadius: 1, objectFit: 'cover' }}
          />
        }
        title={assetName}
        subheader="Sensores asociados"
      />
      <Divider />

      {inactiveSensors.length > 0 && (
        <>
          <Box sx={{ px: 2, pt: 2 }}>
            <Alert severity="warning">
              {inactiveSensors.length === 1
                ? `1 sensor inactivo: ${inactiveSensors[0].name} (ID ${inactiveSensors[0].id})`
                : `${inactiveSensors.length} sensores inactivos`}
            </Alert>
          </Box>
          <Divider sx={{ mt: 2 }} />
        </>
      )}

      <Box sx={{ overflowX: 'auto' }}>
        <Table sx={{ minWidth: 900 }}>
          <TableHead>
            <TableRow>
              <TableCell>Sensor</TableCell>
              <TableCell>ID</TableCell>
              <TableCell>Tipo</TableCell>
              <TableCell align="right">Último valor</TableCell>
              <TableCell>Última transmisión</TableCell>
              <TableCell>Estado</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedSensors.map((s) => {
              const status = computeStatus(s.lastSeen);
              const color = status === 'active' ? 'success' : 'error';
              const label = status === 'active' ? 'Activo' : 'Inactivo';

              return (
                <TableRow
                  hover
                  key={s.id}
                  sx={{ cursor: 'pointer' }}
                  onClick={() => {
                    setSelected(s);
                    setOpen(true);
                  }}
                >
                  <TableCell>{s.name}</TableCell>
                  <TableCell>{s.id}</TableCell>
                  <TableCell sx={{ textTransform: 'capitalize' }}>{s.type}</TableCell>
                  <TableCell align="right">
                    {s.lastValue} {s.unit}
                  </TableCell>
                  <TableCell>{dayjs(s.lastSeen).format('MMM D, YYYY HH:mm')}</TableCell>
                  <TableCell>
                    <Chip color={color as any} label={label} size="small" />
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </Box>

      <Divider />
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', px: 1 }}>
        <Button color="inherit" size="small" variant="text" onClick={toggleShowAll}>
          {showAll ? 'Ver 5' : 'Ver todos'}
        </Button>

        <TablePagination
          component="div"
          count={sensors.length}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          rowsPerPageOptions={[5, 10, 25]}
          labelRowsPerPage="Filas por página"
        />
      </Box>

      {/* Detalle del sensor seleccionado */}
      <Dialog fullWidth maxWidth="md" open={open} onClose={() => setOpen(false)}>
        <DialogTitle>
          {selected ? `${selected.name} — últimas 24 horas` : 'Sensor'}
        </DialogTitle>
        <DialogContent dividers>
          {selected?.history24h && selected.history24h.length > 0 ? (
            <Box sx={{ overflowX: 'auto' }}>
              <Table sx={{ minWidth: 700 }}>
                <TableHead>
                  <TableRow>
                    <TableCell>Hora</TableCell>
                    <TableCell align="right">Valor</TableCell>
                    <TableCell>Unidad</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {selected.history24h.map((p, idx) => (
                    <TableRow key={idx}>
                      <TableCell>{dayjs(p.ts).format('MMM D, HH:mm')}</TableCell>
                      <TableCell align="right">{p.value}</TableCell>
                      <TableCell>{p.unit}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </Box>
          ) : (
            <Alert severity="info">
              Este es un mock de UI: añade <code>history24h</code> al sensor para listar muestras.
            </Alert>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cerrar</Button>
        </DialogActions>
      </Dialog>
    </Card>
  );
}
