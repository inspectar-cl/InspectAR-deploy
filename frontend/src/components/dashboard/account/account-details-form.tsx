/* eslint-disable @typescript-eslint/explicit-function-return-type -- Componente React, tipos inferidos automáticamente */
'use client';

import * as React from 'react';
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import Divider from '@mui/material/Divider';
import Typography from '@mui/material/Typography';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import OutlinedInput from '@mui/material/OutlinedInput';
import Grid from '@mui/material/Grid';
import { useUserToken } from '@/hooks/use-usertoken';
import { decodeJwtToken } from '@/hooks/use-auth'

import Services from '@/modules/Services';

const gs = new Services();

export function AccountDetailsForm(): React.JSX.Element {
  const [receiveEmails, setReceiveEmails] = React.useState(true);
  const [_isSubmitting, setIsSubmitting] = React.useState(false);
  const { user } = useUserToken();
  const payload = decodeJwtToken(user?.token); 

  const username = payload?.username ?? '';
  const email = payload?.email ?? '';
  const role = user?.role;

  const isTecnico = role === 'tecnico';

  const handleCheckboxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setReceiveEmails(event.target.checked);
  };

  const _handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsSubmitting(true);
    
    const _data = new FormData(event.currentTarget);
    /*
    const updatedData = {
      nombre: data.get('nombre') as string,
      apellido: data.get('apellido') as string,
      email: data.get('email') as string,
      telefono: data.get('phone') as string,
      // Solo incluye el estado del checkbox si es técnico
      permitir_contacto_tecnico: isTecnico ? receiveEmails : undefined, 
    };
    */

    const API_PUBLISH_URI = `/actualizar-estado-contacto/1`;

    const requestBody = {
      autorizado: isTecnico ? receiveEmails : undefined, 
    };

    try {
      const response = await gs.authorizedPut(API_PUBLISH_URI, requestBody, user?.token) as { error?: string; data?: unknown };

      if (response.error) {
        /* Intentionally empty - future implementation planned */
      } else {
        /* Intentionally empty - future implementation planned */
      }
    } catch (err) {
      /* Intentionally empty - future implementation planned */
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
      }}
    >
      <Card>
        <CardHeader subheader="La información puede ser editada." title="Perfil" />
        <Divider />
        <CardContent>
          <Grid container spacing={3}>
            <Grid size={{ md:6, xs:12}}>
              <FormControl fullWidth required>
                <InputLabel>Nombre</InputLabel>
                <OutlinedInput defaultValue={username} label="Nombre" name="nombre" />
              </FormControl>
            </Grid>
            <Grid size={{ md:6, xs:12}}>
              <FormControl fullWidth required>
                <InputLabel>Apellido</InputLabel>
                <OutlinedInput defaultValue="Edificio" label="Apellido" name="apellido" />
              </FormControl>
            </Grid>
            <Grid size={{ md:6, xs:12}}>
              <FormControl fullWidth required>
                <InputLabel>Correo electrónico</InputLabel>
                <OutlinedInput defaultValue={email} label="Email address" name="email" />
              </FormControl>
            </Grid>
            <Grid size={{ md:6, xs:12}}>
              <FormControl fullWidth>
                <InputLabel>Número de telefono</InputLabel>
                <OutlinedInput label="Phone number" name="phone" type="tel" />
              </FormControl>
            </Grid>
            {isTecnico ? <Grid size={{ xs:12 }}>
              <FormControlLabel 
                control={
                  <Checkbox 
                    checked={receiveEmails} 
                    onChange={handleCheckboxChange} 
                    name="receiveEmails"
                  />
                }
                label={
                  <Typography variant="body1" color="text.primary">
                    Deseo ser contactado para otorgar ayuda Técnica en mis edificios registrados
                  </Typography>
                }
              />
            </Grid> : null}
          </Grid>
        </CardContent>
        <Divider />
        <CardActions sx={{ justifyContent: 'flex-end' }}>
          <Button variant="contained">Guardar detalles</Button>
        </CardActions>
      </Card>
    </form>
  );
}
