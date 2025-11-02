'use client';

import * as React from 'react';
import { useUserToken } from '@/hooks/use-usertoken';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Paper from '@mui/material/Paper';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardActions from '@mui/material/CardActions';
import CircularProgress from '@mui/material/CircularProgress';
import Services from '@/modules/Services';
import {decodeJwtToken} from '@/hooks/use-auth'
import type { ApiPublicacion, ApiForoResponse, ApiComentario, ApiComentariosResponse} from './helper';
import type {DecodedJwt} from '@/types/token'

import { useAuthUser } from '@/contexts/user-context';

const gs = new Services();

function makeId(prefix: string): string{
  const base =
    typeof crypto !== 'undefined' && 'randomUUID' in crypto
      ? crypto.randomUUID()
      : `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;
  return `${prefix}-${base}`;
}

const FORO_API_BASE = '/foro-edificio';
const PUBLICACION_API = '/publicacion-foro-id';
const COMENTARIOS_API = '/comentarios-foro-id';

async function fetchForumPosts(
  buildingId: number,
  accessToken: string,
  page: number
): Promise<{ data: ApiPublicacion[] | null; totalItems: number; error: string | null }> {
  let uri = `${FORO_API_BASE}/${buildingId}`;
  if (page > 1) {
    uri += `?pagina=${page}`;
  }

  try {
    const res = await gs.authorizedGet(uri, accessToken) as ApiForoResponse & { error?: { data?: { mensaje?: string }, message?: string } };

    // Manejo de errores de la API
    if (res.error) {
      const msg = res.error.data?.mensaje || res.error?.message || 'Error al cargar el foro.';
      return { data: null, totalItems: 0, error: msg };
    }

    return {
      data: res.foro || [],
      totalItems: res.total_items || 0,
      error: null
    };
  } catch (err) {
    return { data: null, totalItems: 0, error: 'Error de red o conexión al servidor.' };
  }
}

export default function ForoPage(): React.JSX.Element {
  const [posts, setPosts] = React.useState<ApiPublicacion[]>([]);
  const [isPostsLoading, setIsPostsLoading] = React.useState(false);
  const [postsError, setPostsError] = React.useState<string | null>(null);
  const [isPublishing, setIsPublishing] = React.useState(false);
  const [_replyError, setReplyError] = React.useState<string | null>(null);

  //Paginacion
  const [currentPage, setCurrentPage] = React.useState(1);
  const [totalPages, setTotalPages] = React.useState(1);
  const [_itemsPerPage, setItemsPerPage] = React.useState(10);

  const [newPostContent, setNewPostContent] = React.useState('');
  const { user, isLoading, error } = useUserToken();
  const edificios = user?.edificio;
  const token = user?.token;

  const { user: userContext} = useAuthUser();
  const selectedEdificioId = userContext?.selected_edificio?.id;

  const payload: DecodedJwt | null = token ? decodeJwtToken(token) : null;

  const buildingIds = React.useMemo<number[]>(
    () => (Array.isArray(edificios) ? edificios.map((e) => e.id).filter((n) => Number.isFinite(n)) : []),
    [edificios]
  );

   React.useEffect((): void => {
    // Asegurarse de tener token y un ID de edificio válido para hacer la llamada
    const bId = typeof selectedEdificioId === 'number' ? selectedEdificioId : null;
    if (!bId || !token) {
      setPosts([]);
      setTotalPages(1);
      return;
    }

    // Función de carga real
    const loadPosts = async (): Promise<void> => {
      setIsPostsLoading(true);
      setPostsError(null);

      const { data, totalItems, error: fetchError } = await fetchForumPosts(bId, token, currentPage);

      if (fetchError) {
        setPostsError(fetchError);
        setPosts([]);
      } else {
        setPosts(data || []);
        const itemsPerPageFinal = 10; 
        setItemsPerPage(itemsPerPageFinal);
        setTotalPages(Math.ceil(totalItems / itemsPerPageFinal));
      }

      setIsPostsLoading(false);
    };

    void loadPosts();

  }, [selectedEdificioId, currentPage, token]);

  // Aqui iria el post de las publicaciones
  const handlePublish = async (): Promise<void> => {
    const content = newPostContent.trim();
    const bId = typeof selectedEdificioId === 'number' ? selectedEdificioId : null;

    // Validacion
    if (!content || !bId || !token) {
      setPostsError('Error: Faltan datos para publicar (contenido, edificio o token expirado).');
      return;
    }

    setIsPublishing(true);
    setPostsError(null);

    const API_PUBLISH_URI = `${PUBLICACION_API}/${bId}`;

    const requestBody = {
      // NOTA: tipo siempre será "falla agua". Esto a futuro podria ser seleccionable, 
      tipo: "falla agua", 
      descripcion: content,
    };

    try {
      // Llamado a la API
      const res = await gs.authorizedPost(API_PUBLISH_URI, requestBody, token) as ApiPublicacion & { error?: { data?: { mensaje?: string }, message?: string } };
      if (res.error) {
        const msg = res.error.data?.mensaje || res.error.message || 'Error al publicar la falla.';
        setPostsError(msg);
      } else {
        // Publicación exitosa:
        const { data: updatedPosts } = await fetchForumPosts(bId, token, 1);
        if (updatedPosts) {
          setPosts(updatedPosts);
        }

        // Limpiar el campo
        setNewPostContent('');
        setCurrentPage(1);
      }

    } catch (err) {
      setPostsError('Error de red al intentar publicar.');
    } finally {
      setIsPublishing(false);
    }
  };


  const handleReply = async (postId: number, replyText: string): Promise<void> => {
    const content = replyText.trim();
    const currentToken = token;

    if (!content || !currentToken || !payload) {
      setReplyError('Error: Comentario vacío o token expirado.');
      return;
    }

    setReplyError(null);

    const uri = `${COMENTARIOS_API}/${postId}`;
    const requestBody = {
    comentario: content,
    };

    try {
      const res = await gs.authorizedPost(uri, requestBody, currentToken) as ApiComentariosResponse;

      if (res.error) {
        const msg = res.error.data?.mensaje || 'Error al publicar el comentario.';
        setReplyError(msg);
        return;
      }

      const newComment: ApiComentario = {
          id_comentario: res.id_comentario || makeId('reply'),
          comentario: content,
          fecha_comentario: new Date().toISOString(),
          username: payload.username || 'Usuario Actual',
      };

      setPosts(prevPosts => 
          prevPosts.map(post => {
              if (post.id_falla === postId) {
                  const currentComments = post.comentarios ?? []; 
                  return {
                      ...post,
                      comentarios: [...currentComments, newComment], 
                  };
              }
              return post;
          })
      );

    } catch (err) {
      setReplyError('Error de red al intentar responder.');
    } 
  };

  if (isLoading) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box sx={{ py: 8, textAlign: 'center' }}>
        <Typography color="error">{error}</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ display: 'grid', gap: 3 }}>
      <Typography variant="h4">Foro del edificio</Typography>
      <Typography variant="body2" color="text.secondary">
        Comparte novedades sobre los activos de tus edificios.
      </Typography>

      {/* Composer */}
      <Paper variant="outlined" sx={{ p: 2, borderRadius: 2 }}>
        <Stack spacing={2}>
          <TextField
            multiline
            minRows={3}
            placeholder="¿Alguna novedad en algún activo de tus edificios?"
            value={newPostContent}
            onChange={(e) => {
              setNewPostContent(e.target.value);
            }}
          />
          <Box sx={{ textAlign: 'right' }}>
            <Button variant="contained"
             onClick={handlePublish} 
             disabled={!newPostContent.trim() || !buildingIds.length}>
              {isPublishing ? 'Publicando...' : 'Publicar'} {/* Feedback visual */}
            </Button>
          </Box>
        </Stack>
      </Paper>

      <Divider />

      {isPostsLoading ? (
        <Box sx={{ textAlign: 'center', py: 6 }}>
          <CircularProgress />
          <Typography>Cargando publicaciones...</Typography>
        </Box>
      ) : postsError ? (
        <Box sx={{ textAlign: 'center', py: 6 }}>
          <Typography color="error">{postsError}</Typography>
        </Box>
      ) : !posts.length ? (
        <Box sx={{ textAlign: 'center', py: 6 }}>
          <Typography variant="h6" color="text.secondary">
            Aún no hay publicaciones para este edificio.
          </Typography>
        </Box>
      ) : (
        <Stack spacing={2}>
          {posts.map((p, index) => (
            <PostCard 
              key={p.id_falla || `temp-${index}-${Date.now()}`}
              post={p}
              onReply={(text) => handleReply(p.id_falla, text)} 
            />
          ))}

          {/* Control de Paginación */}
          {totalPages > 1 && (
            <Stack direction="row" justifyContent="center" spacing={2} sx={{ mt: 3 }}>
              <Button
                onClick={() => {setCurrentPage((p) => p - 1)}}
                disabled={currentPage === 1}
              >
              Anterior
              </Button>
              <Typography variant="body2" sx={{ alignSelf: 'center' }}>
                Página {currentPage} de {totalPages}
              </Typography>
              <Button
                onClick={() => {setCurrentPage((p) => p + 1)}}
                disabled={currentPage >= totalPages}
              >
                Siguiente
              </Button>
            </Stack>
          )}
        </Stack>
      )}
    </Box>
  );
}

function PostCard({
  post,
  onReply
}: {
  post: ApiPublicacion;
  onReply: (replyText: string) => void;
}): React.JSX.Element {
  const [replyText, setReplyText] = React.useState('');

  const comentarios = post.comentarios ?? []; 

  return (
    <Card variant="outlined" sx={{ borderRadius: 2 }}>
      <CardContent>
        <Stack spacing={1.25}>
          <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
            {post.descripcion}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            Publicado por {post.username} el {new Date(post.fecha_publicacion).toLocaleString()} • Edificio #{post.id_edificio}
          </Typography>
        </Stack>
      </CardContent>

      {Boolean(comentarios.length) && (
        <CardContent sx={{ pt: 0 }}>
          <Divider sx={{ mb: 1.5 }} />
          <Stack spacing={1.25} sx={{ pl: { xs: 0, sm: 1.5 } }}>
          {comentarios.map((r) => (
            <Reply 
                key={r.id_comentario} 
                content={r.comentario} 
                createdAt={r.fecha_comentario} 
                username={r.username}
              />
            ))}
          </Stack>
        </CardContent>
      )}

      <CardActions sx={{ px: 2, pb: 2 }}>
        <Stack direction="row" spacing={1} sx={{ width: '100%' }}>
          <TextField
            size="small"
            fullWidth
            placeholder='Reply'
            value={replyText}
            onChange={(e) => {
              setReplyText(e.target.value);
            }}
          />
          <Button
            variant="outlined"
            onClick={() => {
              onReply(replyText);
              setReplyText('');
            }}
            disabled={!replyText.trim()}
          >
            Responder
          </Button>
        </Stack>
      </CardActions>
    </Card>
  );
}

function Reply({ content, createdAt, username }: { content: string; createdAt: string; username: string }): React.JSX.Element {
  return (
    <Paper variant="outlined" sx={{ p: 1.25, borderRadius: 2 }}>
      <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
        {content}
      </Typography>
      <Typography variant="caption" color="text.secondary">
        Por {username} el {new Date(createdAt).toLocaleString()}
      </Typography>
    </Paper>
  );
}
