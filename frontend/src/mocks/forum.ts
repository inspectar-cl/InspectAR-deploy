// src/mocks/forum.ts
export interface ForumReply {
  id: string;          // id de la respuesta
  postId: string;      // id del tema al que responde
  content: string;     // contenido de la respuesta
  createdAt: string;   // date
};

export interface ForumPost {
  id: string;               // id del tema
  buildingId: number;       // id del edificio asociado
  content: string;          // contenido del tema
  createdAt: string;
  replies: ForumReply[];
};

export const forumPostsMock: ForumPost[] = [
  {
    id: 'p-1001',
    buildingId: 1,
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
    buildingId: 1,
    content: 'Pedro dijo: Caldera con ruido intermitente, quizá buje suelto.',
    createdAt: '2025-10-03T18:40:00.000Z',
    replies: []
  },
  {
    id: 'p-1003',
    buildingId: 2,
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
