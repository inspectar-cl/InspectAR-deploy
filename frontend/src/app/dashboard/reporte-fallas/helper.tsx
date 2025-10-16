export interface ApiComentario {
  comentario: string;
  fecha_comentario: string;
  id_comentario: number | string;
  username: string;
}

export interface ApiComentariosResponse {
  id_comentario: number | string;
  error?: { data: { mensaje: string } | undefined; message: string; };
}

export interface ApiPublicacion {
  comentarios: ApiComentario[];
  descripcion: string; // Esto será el "content" del post
  estado: string;
  fecha_publicacion: string; // ISO String (createdAt)
  id_edificio: number;
  id_falla: number; // Esto será el ID del post
  tipo: string;
  username: string;
}

export interface ApiForoResponse {
  edificio_id: number;
  foro: ApiPublicacion[];
  total_items: number;
  items_per_page?: number;
  pagina?: number;
}

export interface Edificio {
  id: number;
  nombre: string;
  direccion: string;
  creado_en: string;
}