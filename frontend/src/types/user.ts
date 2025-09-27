export interface User {
  id: string;
  name?: string;
  avatar?: string;
  email?: string;
  role?: string;
  edificio?: string;
  [key: string]: unknown;
}
