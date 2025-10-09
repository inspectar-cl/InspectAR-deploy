'use client';

import * as React from 'react';
import { useUser } from '@/hooks/use-user';
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

type EdificioAPI = {
  id: number;
  nombre: string;
  direccion: string;
  creado_en: string;
};

function isoNow(): string {
  return new Date().toISOString();
}
function makeId(prefix: string) {
  const base =
    typeof crypto !== 'undefined' && 'randomUUID' in crypto
      ? crypto.randomUUID()
      : `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;
  return `${prefix}-${base}`;
}

export default function ForoPage(): React.JSX.Element {
  const { user, isLoading, error } = useUser();
  const edificios = (user as any)?.edificioId as EdificioAPI[] | undefined;

  const buildingIds = React.useMemo<number[]>(
    () => (Array.isArray(edificios) ? edificios.map((e) => e.id).filter((n) => Number.isFinite(n)) : []),
    [edificios]
  );

  const [posts, setPosts] = React.useState<ForumPost[]>([]);
  const [newPostContent, setNewPostContent] = React.useState('');
  const [selectedBuildingId, setSelectedBuildingId] = React.useState<number | ''>('');

  React.useEffect(() => {
    setSelectedBuildingId(buildingIds[0] ?? '');
  }, [buildingIds]);

  //Cargar posts de los edificios del usuario
  React.useEffect(() => {
    if (!buildingIds.length) {
      setPosts([]);
      return;
    }
    const idSet = new Set(buildingIds);
    const initial = forumPostsMock
      .filter((p) => idSet.has(p.buildingId))
      .sort((a, b) => b.createdAt.localeCompare(a.createdAt));
    setPosts(initial);
  }, [buildingIds]);

  const handlePublish = (): void => {
    const content = newPostContent.trim();
    const bId = typeof selectedBuildingId === 'number' ? selectedBuildingId : buildingIds[0];

    if (!content || !Number.isFinite(bId)) return;

    const newPost: ForumPost = {
      id: makeId('p'),
      buildingId: bId,
      content,
      createdAt: isoNow(),
      replies: []
    };

    setPosts((prev) => [newPost, ...prev]);
    setNewPostContent('');
  };

  const handleReply = (postId: string, replyText: string): void => {
    const content = replyText.trim();
    if (!content) return;
    const reply: ForumReply = {
      id: makeId('r'),
      postId,
      content,
      createdAt: isoNow()
    };
    setPosts((prev) =>
      prev.map((p) => {
        if (p.id !== postId) return p;
        return { ...p, replies: [...p.replies, reply] };
      })
    );
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

      {/* Listado de temas */}
      {!posts.length ? (
        <Box sx={{ textAlign: 'center', py: 6 }}>
          <Typography variant="h6" color="text.secondary">
            Aún no hay publicaciones en tus edificios
          </Typography>
        </Box>
      ) : (
        <Stack spacing={2}>
          {posts.map((p) => (
            <PostCard key={p.id} post={p} onReply={handleReply} />
          ))}
        </Stack>
      )}
    </Box>
  );
}

function PostCard({
  post,
  onReply
}: {
  post: ForumPost;
  onReply: (postId: string, replyText: string) => void;
}): React.JSX.Element {
  const [replyText, setReplyText] = React.useState('');

  return (
    <Card variant="outlined" sx={{ borderRadius: 2 }}>
      <CardContent>
        <Stack spacing={1.25}>
          <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
            {post.content}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {new Date(post.createdAt).toLocaleString()} • Edificio #{post.buildingId}
          </Typography>
        </Stack>
      </CardContent>

      {!!post.replies.length && (
        <CardContent sx={{ pt: 0 }}>
          <Stack spacing={1.25} sx={{ pl: { xs: 0, sm: 1.5 } }}>
            {post.replies.map((r) => (
              <Reply key={r.id} content={r.content} createdAt={r.createdAt} />
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
              onReply(post.id, replyText);
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

function Reply({ content, createdAt }: { content: string; createdAt: string }): React.JSX.Element {
  return (
    <Paper variant="outlined" sx={{ p: 1.25, borderRadius: 2 }}>
      <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
        {content}
      </Typography>
      <Typography variant="caption" color="text.secondary">
        {new Date(createdAt).toLocaleString()}
      </Typography>
    </Paper>
  );
}
