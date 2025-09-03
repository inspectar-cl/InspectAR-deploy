import * as React from 'react';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardMedia from '@mui/material/CardMedia';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import CardActions from '@mui/material/CardActions';
import SendIcon from '@mui/icons-material/Send';
import TextField from '@mui/material/TextField';
import { Box } from '@mui/system';
import { Divider } from '@mui/material';

import murcy from '../../../../assets/avatar-11.png';

export function ChatBotCard() {
  return (
    <Card sx={{ maxWidth: 345 , minHeight: 520}}>
        <CardMedia
          component="img"
          height="140"
          image="https://mui.com/static/images/cards/contemplative-reptile.jpg"
          alt="green iguana"
        />
        <CardContent sx={{ pt: 2, pb: 1 }}>
          <Typography gutterBottom variant="h5" component="div" align="center">
            Murcy
          </Typography>
          <Typography variant="body2" sx={{ color: 'text.secondary' }} align="center">
            Preguntale a Murcy alguna duda sobre la ficha tecnica del activo
          </Typography>
        </CardContent>

        <Box component="form" sx={{ m: 2}}>
            <TextField id="standard-basic" multiline 
            rows={5} maxRows={5} fullWidth disabled 
            value={"Hola esta es una prueba de lo que puede generar el chat: AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"}/>
        </Box>

        <Divider />

        <Box component="form" sx={{ m: 2 }}>
            <TextField id="standard-basic" label="Pregunta" variant="standard" multiline maxRows={3} fullWidth/>
        </Box>
        <CardActions>
            <Button variant="contained" size="small" startIcon={<SendIcon />} >
                Enviar
            </Button>
        </CardActions>
    </Card>
  );
}
