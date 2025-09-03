import * as React from 'react';
import { Box, Typography, TextField, Button, Card, CardContent, Divider, Fab} from '@mui/material';
import CardMedia from '@mui/material/CardMedia';
import CardActions from '@mui/material/CardActions';
import SendIcon from '@mui/icons-material/Send';
import Slide from '@mui/material/Slide';
import { useTheme } from '@mui/material/styles';
import ChatIcon from '@mui/icons-material/Chat';


export function ChatBotCard() {
  const [open, setOpen] = React.useState(false);
  const [inputR, setInput] = React.useState('');
  const [respuesta, setRespuesta] = React.useState('');
  const [loading, setLoading] = React.useState(false);

  const theme = useTheme();

  const inspai = '../../../../assets/inspai.png';

  const API_URL = '/api/chat';

  const handleSend = async () => {
    if (!inputR.trim()) return;

    setLoading(true);
    try {
      const res = await fetch(API_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pregunta: inputR }), // la API recibe la pregunta
      });

      const data = await res.json();
      setRespuesta(data.respuesta); // la API devuelve la respuesta
      setInput('');
    } catch (error) {
      console.error('Error llamando a la API:', error);
      setRespuesta('Error al obtener respuesta del bot');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      {/* Boton flotante */}
      <Fab
        color="primary"
        sx={{
          position: 'fixed',
          bottom: 16,
          right: 16,
          zIndex: 2000,
        }}
        onClick={() => setOpen((prev) => !prev)}
        >
        <ChatIcon />
      </Fab>

      {/* Contenedor del Chat flotante */}
      {open && (
        <Slide direction="left" in={open} mountOnEnter unmountOnExit timeout={{enter:700}}>
          <Box
            sx={{
              position: 'fixed',
              bottom: 80,
              right: 16,
              width: 360,
              zIndex: 2000,
            }}
          >
            <Card sx={{ maxWidth: 345 , minHeight: 520}}>
                <CardMedia
                  component="img"
                  height="140"
                  image= {inspai}
                  alt="InspecAR PET"
                />
                <CardContent sx={{ pt: 2, pb: 1 }}>
                  <Typography gutterBottom variant="h5" component="div" align="center">
                    Inspy
                  </Typography>
                  <Typography variant="body2" sx={{ color: 'text.secondary' }} align="center">
                    Preguntale a Inspy alguna duda sobre la ficha tecnica del activo
                  </Typography>
                </CardContent>

                <Box component="form" sx={{ m: 2}}>
                    <TextField
                      multiline
                      rows={5}
                      maxRows={5}
                      fullWidth
                      focused
                      InputProps={{ readOnly: true}}
                      value={
                        loading
                          ? 'Generando respuesta...'
                          : respuesta || 'Aquí aparecerá la respuesta del bot'
                      }
                    />
                </Box>

                <Divider />

                <Box component="form" sx={{ m: 2 }}>
                    <TextField
                      label="Pregunta"
                      variant="standard"
                      multiline
                      rows={3}
                      maxRows={3}
                      fullWidth
                      value={inputR}
                      onChange={(e) => setInput(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && handleSend()}
                    />
                </Box>
                <CardActions>
                    <Button
                      variant="contained"
                      size="small"
                      startIcon={<SendIcon />}
                      onClick={handleSend}
                      disabled={loading}
                    >
                      Enviar
                    </Button>
                </CardActions>
            </Card>
          </Box>
        </Slide>
        )}
    </>
  );
}
