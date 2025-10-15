'use client';

import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import Divider from '@mui/material/Divider';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import { useUserToken } from '@/hooks/use-usertoken';
import { decodeJwtToken } from '@/hooks/use-auth'

type Role = 'residente' | 'admin' | 'analista' | 'tecnico' | null;

const capitalizeRole = (role: Role): string => {
  if (role === null) {
    return 'Error: Sin rol';
  }

  const roleMap: Record<NonNullable<Role>, string> = {
    'residente': 'Residente',
    'admin': 'Admin',
    'analista': 'Analista',
    'tecnico': 'Técnico', 
  };

  return roleMap[role] ?? 'Desconocido';
};

export function AccountInfo(): React.JSX.Element {

  const { user } = useUserToken();
  const payload = decodeJwtToken(user?.token); 

  const username = payload?.username ?? '';
  const role = capitalizeRole(user?.role ?? null);

  const usuario = {
    name: username,
    avatar: '/assets/avatar.png',
    rol: role,
  } as const;

  return (
    <Card>
      <CardContent>
        <Stack spacing={2} sx={{ alignItems: 'center' }}>
          <div>
            <Avatar src={usuario.avatar} sx={{ height: '80px', width: '80px' }} />
          </div>
          <Stack spacing={1} sx={{ textAlign: 'center' }}>
            <Typography variant="h5">{usuario.name} Edificio</Typography>
            <Typography color="text.secondary" variant="body2">
              {usuario.rol}
            </Typography>
          </Stack>
        </Stack>
      </CardContent>
      <Divider />
      <CardActions>
        <Button fullWidth variant="text">
          Subir foto
        </Button>
      </CardActions>
    </Card>
  );
}
