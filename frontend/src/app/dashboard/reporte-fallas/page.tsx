'use client';

import * as React from 'react';
import { useUserToken } from '@/hooks/use-usertoken';
import { forumPostsMock, type ForumPost, type ForumReply } from '@/mocks/forum';
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
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import Services from '@/modules/Services';
import { ApiPublicacion, ApiForoResponse} from './helper';


const gs = new Services();

function makeId(prefix: string) {
  const base =
    typeof crypto !== 'undefined' && 'randomUUID' in crypto
      ? crypto.randomUUID()
      : `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;
  return `${prefix}-${base}`;
}

const FORO_API_BASE = '/foro-edificio';

/**
 * Carga las publicaciones del foro para un edificio y página específicos.
 * @param buildingId El ID del edificio a consultar.
 * @param accessToken El token JWT para la autorización.
 * @param page La página a cargar (opcional, por defecto es 1).
 */
async function fetchForumPosts(
  buildingId: number,
  accessToken: string,
  page: number = 1
): Promise<{ data: ApiPublicacion[] | null; totalItems: number; error: string | null }> {
  // Construir la URI con el ID del edificio y el parámetro de paginación
  let uri = `${FORO_API_BASE}/${buildingId}`;
  if (page > 1) {
    uri += `?pagina=${page}`;
  }

  try {
    // Utilizamos authorizedGet o authorizedPost si tu gs lo soporta.
    // Aquí usaremos authorizedGet, asumiendo que lo agregaste a Services como discutimos antes.
    const res = await gs.authorizedGet(uri, accessToken) as ApiForoResponse & { error?: any };

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
    console.error('Error fetching forum posts:', err);
    return { data: null, totalItems: 0, error: 'Error de red o conexión al servidor.' };
  }
}

export default function ForoPage(): React.JSX.Element {
  const { user, isLoading, error } = useUserToken();
  const edificios = user?.edificio;
  const token = user?.token;

  const buildingIds = React.useMemo<number[]>(
    () => (Array.isArray(edificios) ? edificios.map((e) => e.id).filter((n) => Number.isFinite(n)) : []),
    [edificios]
  );

  const [posts, setPosts] = React.useState<ApiPublicacion[]>([]);
  const [isPostsLoading, setIsPostsLoading] = React.useState(false);
  const [postsError, setPostsError] = React.useState<string | null>(null);

  //Paginacion
  const [currentPage, setCurrentPage] = React.useState(1);
  const [totalPages, setTotalPages] = React.useState(1);
  const [itemsPerPage, setItemsPerPage] = React.useState(10);
  //Seleccion de edificio
  const [selectedBuildingId, setSelectedBuildingId] = React.useState<number | ''>('');

  const [newPostContent, setNewPostContent] = React.useState('');

  //Seteo de primer edificio por defecto (de la lista de edificios tomo el primero)
  React.useEffect(() => {
    if (buildingIds.length > 0 && selectedBuildingId === '') {
      setSelectedBuildingId(buildingIds[0] ?? '');
    }
  }, [buildingIds, selectedBuildingId]);

   React.useEffect(() => {
    // Asegurarse de tener token y un ID de edificio válido para hacer la llamada
    const bId = typeof selectedBuildingId === 'number' ? selectedBuildingId : null;
    if (!bId || !token) {
      setPosts([]);
      setTotalPages(1);
      return;
    }

    // Función de carga real
    const loadPosts = async () => {
      setIsPostsLoading(true);
      setPostsError(null);

      const { data, totalItems, error } = await fetchForumPosts(bId, token, currentPage);

      if (error) {
        setPostsError(error);
        setPosts([]);
      } else {
        setPosts(data || []);
        // Calcula el total de páginas, asumiendo 10 elementos por página si la API no lo devuelve
        const itemsPerPageFinal = 10; 
        setItemsPerPage(itemsPerPageFinal);
        setTotalPages(Math.ceil(totalItems / itemsPerPageFinal));
      }

      setIsPostsLoading(false);
    };

    loadPosts();

  }, [selectedBuildingId, currentPage, token]);

  // Manejar el cambio de edificio, forzando la vuelta a la página 1
  const handleBuildingChange = (id: number | ''): void => {
    setSelectedBuildingId(id);
    setCurrentPage(1); // Siempre resetear a la primera página al cambiar de edificio
  };

  // Los handlers de publicación y respuesta ahora son de "mock" temporalmente
  // Aqui iria el post de las publicaciones
  const handlePublish = (): void => { /* ... (mock por ahora) ... */ };
  const handleReply = (postId: number, replyText: string): void => { 
    // Aquí se debería llamar a una API de POST para crear un comentario
    console.log(`Respuesta al post ${postId}: ${replyText}`);
    // Por ahora solo loguea, ya que tu backend no tiene un POST implementado aquí
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
          {/* aqui se selecciona el edificio si el usuario tiene mas, es un filtro*/}
          {buildingIds.length > 1 && (
            <FormControl size="small" sx={{ minWidth: 220 }}>
              <InputLabel id="edificio-select-label">Publicar en</InputLabel>
              <Select
                labelId="edificio-select-label"
                label="Publicar en"
                value={selectedBuildingId}
                onChange={(e) => setSelectedBuildingId(Number(e.target.value))}
              >
                {edificios?.map((e) => (
                  <MenuItem key={e.id} value={e.id}>
                    {e.nombre ? `${e.nombre} (#${e.id})` : `Edificio #${e.id}`}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          )}

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
            <Button variant="contained" onClick={handlePublish} disabled={!newPostContent.trim() || !buildingIds.length}>
              Publicar
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
          {posts.map((p) => (
            <PostCard 
              key={p.id_falla}
              post={p} 
              onReply={(text) => handleReply(p.id_falla, text)} 
            />
          ))}

          {/* Control de Paginación */}
          {totalPages > 1 && (
            <Stack direction="row" justifyContent="center" spacing={2} sx={{ mt: 3 }}>
              <Button
                onClick={() => setCurrentPage((p) => p - 1)}
                disabled={currentPage === 1}
              >
              Anterior
              </Button>
              <Typography variant="body2" sx={{ alignSelf: 'center' }}>
                Página {currentPage} de {totalPages}
              </Typography>
              <Button
                onClick={() => setCurrentPage((p) => p + 1)}
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

  return (
    <Card variant="outlined" sx={{ borderRadius: 2 }}>
      <CardContent>
        <Stack spacing={1.25}>
          <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
            {post.descripcion}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            Publicado por **{post.username}** el {new Date(post.fecha_publicacion).toLocaleString()} 
            • Edificio #{post.id_edificio} • Estado: {post.estado} • Tipo: {post.tipo}
          </Typography>
        </Stack>
      </CardContent>

      {!!post.comentarios.length && (
        <CardContent sx={{ pt: 0 }}>
          <Divider sx={{ mb: 1.5 }} />
          <Stack spacing={1.25} sx={{ pl: { xs: 0, sm: 1.5 } }}>
          {post.comentarios.map((r) => (
            <Reply 
                key={r.id_comentario} 
                content={r.comentario} 
                createdAt={r.fecha_comentario} 
                username={r.username} // Añadimos el username al Reply
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
        Por **{username}** el {new Date(createdAt).toLocaleString()}
      </Typography>
    </Paper>
  );
}
