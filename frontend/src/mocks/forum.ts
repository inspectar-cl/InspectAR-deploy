// src/mocks/forum.ts
export interface ForumReply {
  id: string;          // id de la respuesta
  postId: string;      // id del tema al que responde
  content: string;     // puedes incluir "Autor dijo: ..." aquí mismo
  createdAt: string;   // ISO date
};

export interface ForumPost {
  id: string;               // id del tema
  buildingId: string;       // "EA" | "EB" | ...
  content: string;          // "Juanito dijo: ...", etc.
  createdAt: string;
  replies: ForumReply[];
};

export const forumPostsMock: ForumPost[] = [
  {
    id: 'p-1001',
    buildingId: 'EA',
    content: 'Juanito dijo: La luz del ascensor parpadea a ratos, posible falla eléctrica.',
    createdAt: '2025-10-04T11:05:00.000Z',
    replies: [
      {
        id: 'r-2001',
        postId: 'p-1001',
        content: 'María reply to "Juanito": Ayer también lo noté en el piso 3.',
        createdAt: '2025-10-04T12:20:00.000Z',
      }
    ]
  },
  {
    id: 'p-1002',
    buildingId: 'EA',
    content: 'Pedro dijo: Caldera con ruido intermitente, quizá buje suelto.',
    createdAt: '2025-10-03T18:40:00.000Z',
    replies: []
  },
  {
    id: 'p-1003',
    buildingId: 'EB',
    content: 'Ana dijo: Bomba de agua marcó presión alta anoche.',
    createdAt: '2025-10-03T09:15:00.000Z',
    replies: [
      {
        id: 'r-2002',
        postId: 'p-1003',
        content: 'Luis reply to "Ana": Se normalizó por la mañana.',
        createdAt: '2025-10-03T10:05:00.000Z',
      }
    ]
  }
];
