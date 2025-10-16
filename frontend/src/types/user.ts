interface EdificioAPI {
  id: number;
  nombre: string;
  direccion: string;
  creado_en: string;
}

export interface User {
  id: string;
  name?: string;
  avatar?: string;
  email?: string;
  role?: string;
  edificio?: EdificioAPI[];
  [key: string]: unknown;
}
