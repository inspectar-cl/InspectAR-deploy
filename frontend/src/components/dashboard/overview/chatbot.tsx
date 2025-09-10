import * as React from 'react';
import { Box, Typography, TextField, Button, Card, CardContent, Divider, Fab} from '@mui/material';
import CardMedia from '@mui/material/CardMedia';
import CardActions from '@mui/material/CardActions';
import SendIcon from '@mui/icons-material/Send';
import Slide from '@mui/material/Slide';
import ChatIcon from '@mui/icons-material/Chat';

import Services from '@/modules/Services'

const gs = new Services()

interface Message {
  pregunta: string;
  respuesta: string;
}

export function ChatBotCard ({ id }: { id: string | number | null}): React.JSX.Element {
  const [open, setOpen] = React.useState(false);
  const [inputR, setInput] = React.useState('');
  const [loading, setLoading] = React.useState(false);

  const [messages, setMessages] = React.useState<Message[]>([]);

  const inspy = '../../../../assets/inspai.png';

  const handleSend = async (): Promise<void> => {
    if (!inputR.trim()) return;

    setLoading(true);
    const API_URL = `/documentacion/documentos/${id}/consultar`;

    try {
      const response = await gs.post(API_URL, { pregunta: inputR }) as { mensaje?: string; respuesta?: string };

      const newMessage = {
        pregunta: inputR,
        respuesta:
          response?.mensaje === 'error inesperado'
            ? 'Error al obtener respuesta de Inspy'
            : response.respuesta || 'Sin respuesta de Inspy',
      };

      setMessages((prev) => [...prev, newMessage]);
      setInput('');
    } catch (error) {
      // Error handling for API call failure
      setMessages((prev) => [
        ...prev,
        { pregunta: inputR, respuesta: 'Error al obtener respuesta de Inspy' },
      ]);
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
        onClick={() => { setOpen((prev) => !prev); }}
        >
        <ChatIcon />
      </Fab>

      {/* Contenedor del Chat flotante */}
      {open ? <Slide direction="left" in={open} mountOnEnter unmountOnExit timeout={{enter:700}}>
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
                  image= {inspy}
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

                <Box component="form" sx={{
                  m: 2,
                  minHeight: 200,
                  maxHeight: 250,
                  overflowY: "auto",
                  border: "1px solid #6e6e6eff",
                  borderRadius: 1,
                  p: 1,
                }}>
                    {messages.length === 0 && (
                      <Typography variant="body2" color="text.secondary">
                        Hola, soy Inspy ¡Estoy para ayudarte!
                      </Typography>
                    )}

                    {messages.map((msg, index) => {
                      const isLast = index === messages.length - 1;
                      const typed = msg.respuesta;

                      return (
                        <Box key={msg.pregunta} sx={{ mb: 2 }}>
                          <Typography variant="subtitle2" color="primary" marginBottom="5px">
                            Tu pregunta: {msg.pregunta}
                          </Typography>
                          <Typography variant="body2">
                            {isLast && loading ? 'Generando respuesta...' : isLast ? typed : msg.respuesta}
                          </Typography>
                          <Divider sx={{ my: 1 }} />
                        </Box>
                      );
                    })}
                </Box>

                <Divider />

                <Box component="form" sx={{ m: 2 }}>
                    <TextField
                      label="Pregunta"
                      variant="standard"
                      multiline
                      minRows={3}
                      maxRows={3}
                      fullWidth
                      value={inputR}
                      onChange={(e) => { setInput(e.target.value); }}
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
        </Slide> : null}
    </>
  );
}
