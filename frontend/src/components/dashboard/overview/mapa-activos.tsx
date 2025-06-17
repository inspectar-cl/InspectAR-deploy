'use client';

import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import Card from '@mui/material/Card';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import Box from '@mui/material/Box';
import { useState } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import { activosMock } from '@/mocks';

import ubicacion from '/assets/ubicacion.png'

const edificios = [
  { id: 'EA', nombre: 'Edificio A', lat: -33.45, lng: -70.66 },
  { id: 'EB', nombre: 'Edificio B', lat: -33.456, lng: -70.662 },
];

export function MapaActivos({ sx }: { sx?: any }) {
  const [edificioSeleccionado, setEdificioSeleccionado] = useState<string | null>(null);

  const activosPorEdificio = activosMock.filter(a => a.id_edificio === edificioSeleccionado);

  return (
    <Card sx={sx}>
        <Box display="flex" flexDirection={{ xs: 'column', md: 'row' }} gap={2} p={2}>
        <Box sx={{ width: '100%', height: 500 }}>
            <MapContainer center={[-33.45, -70.66]} zoom={15} style={{ width: '100%', height: '100%' }}>
            <TileLayer 
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
            {edificios.map((edificio) => (
                <Marker
                key={edificio.id}
                position={[edificio.lat, edificio.lng]}
                eventHandlers={{
                    click: () => setEdificioSeleccionado(edificio.id),
                }}
                icon={L.icon({
                    iconUrl: 'https://cdn-icons-png.flaticon.com/512/8/8214.png',
                    iconSize: [45, 41],
                    iconAnchor: [23, 21],
                })}
                >
                <Popup>{edificio.nombre}</Popup>
                </Marker>
            ))}
            </MapContainer>
        </Box>
        </Box>
    </Card>
  );
}
