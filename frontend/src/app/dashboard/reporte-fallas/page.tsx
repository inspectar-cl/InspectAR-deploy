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
import { v4 as uuid } from 'uuid';

function isoNow(): string {
  return new Date().toISOString();
}

export default function ForoPage(): React.JSX.Element {
  const { user, isLoading, error } = useUser();
  const edificioId = user?.edificio;

  const [posts, setPosts] = React.useState<ForumPost[]>([]);
  const [newPostContent, setNewPostContent] = React.useState('');

  React.useEffect(() => {
    if (!edificioId) return;
    // carga inicial: filtrar por edificio y ordenar por fecha desc
    const initial = forumPostsMock
      .filter(p => p.buildingId === edificioId)
      .sort((a, b) => b.createdAt.localeCompare(a.createdAt));
    setPosts(initial);
  }, [edificioId]);

  const handlePublish = (): void => {
    const content = newPostContent.trim();
    if (!content) return;
    // Componer contenido estilo "Nombre dijo: ..." si quieren forzarlo
    // Aquí dejamos lo que escriba el usuario tal cual:
    const newPost: ForumPost = {
      id: `p-${uuid()}`,
      buildingId: edificioId ?? 'NN',
      content,
      createdAt: isoNow(),
      replies: []
    };
    setPosts(prev => [newPost, ...prev]);
    setNewPostContent('');
  };

  const handleReply = (postId: string, replyText: string): void => {
    const content = replyText.trim();
    if (!content) return;
    const reply: ForumReply = {
      id: `r-${uuid()}`,
      postId,
      content,
      createdAt: isoNow()
    };
    setPosts(prev => prev.map(p => {
      if (p.id !== postId) return p;
      return { ...p, replies: [...p.replies, reply] };
    }));
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
        Comparte novedades sobre los activos de tu edificio.
      </Typography>

      {/* Composer */}
      <Paper variant="outlined" sx={{ p: 2, borderRadius: 2 }}>
        <Stack spacing={2}>
          <TextField
            multiline
            minRows={3}
            placeholder="¿Alguna novedad en algún activo de tu edificio?"
            value={newPostContent}
            onChange={(e) => {setNewPostContent(e.target.value)} }
          />
          <Box sx={{ textAlign: 'right' }}>
            <Button variant="contained" onClick={handlePublish} disabled={!newPostContent.trim()}>
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
            Aún no hay publicaciones en tu edificio
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

function PostCard({ post, onReply }: { post: ForumPost; onReply: (postId: string, replyText: string) => void; }): React.JSX.Element {
  const [replyText, setReplyText] = React.useState('');

  return (
    <Card variant="outlined" sx={{ borderRadius: 2 }}>
      <CardContent>
        <Stack spacing={1.25}>
          <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
            {post.content}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {new Date(post.createdAt).toLocaleString()}
          </Typography>
        </Stack>
      </CardContent>

      {/* Respuestas existentes del mock de momento*/}
      {Boolean(post.replies.length) && (
        <CardContent sx={{ pt: 0 }}>
          <Stack spacing={1.25} sx={{ pl: { xs: 0, sm: 1.5 } }}>
            {post.replies.map((r) => (
              <Reply key={r.id} content={r.content} createdAt={r.createdAt} />
            ))}
          </Stack>
        </CardContent>
      )}
        {/* Composer de respuesta, aqui podriamos obtener el nombre del comentario quizas */}
      <CardActions sx={{ px: 2, pb: 2 }}>
        <Stack direction="row" spacing={1} sx={{ width: '100%' }}>
          <TextField
            size="small"
            fullWidth
            placeholder='Reply to...'
            value={replyText}
            onChange={(e) => {setReplyText(e.target.value)} }
          />
          <Button
            variant="outlined"
            onClick={() => { onReply(post.id, replyText); setReplyText(''); }}
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
