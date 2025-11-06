'use client';

import type { User } from '@/types/user';
import Services from '@/modules/Services';
import {decodeJwtToken} from '@/hooks/use-auth'

type Role = 'analista' | 'tecnico' | 'residente' | 'root';

interface Edificio {
  id: number;
  nombre: string;
  direccion: string;
  creado_en: string;
}

interface StoredPayload {
  token: string;
  edificio: Edificio[] | undefined;
  role: Role;
  selected_edificio: Edificio;
}

interface SessionUser extends User {
  edificio: Edificio[] | undefined;
  role: Role;
  selected_edificio: Edificio;
}

const gs = new Services();

const STORAGE_KEY = 'custom-auth-token';

const baseUser: User = {
  id: 'USR-000',
  avatar: '/assets/avatar.png',
  firstName: 'Sofia',
  lastName: 'Rivers',
  email: 'sofia@devias.io',
};

export interface SignUpParams {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}

export interface SignInWithOAuthParams {
  provider: 'google' | 'discord';
}

export interface SignInWithPasswordParams {
  email: string;
  password: string;
}

export interface ResetPasswordParams {
  email: string;
}

function isRole(v: unknown): v is Role {
  return v === 'analista' || v === 'tecnico' || v === 'residente' || v === 'root';
}

function isEdificio(v: unknown): v is Edificio {
  if (typeof v !== 'object' || v === null) return false;
  const obj = v as Record<string, unknown>;
  // Una comprobación mínima para un Edificio
  return (
    typeof obj.id === 'number' &&
    typeof obj.nombre === 'string'
  );
}

function isStoredPayload(v: unknown): v is StoredPayload {
  if (typeof v !== 'object' || v === null) return false;
  const obj = v as Record<string, unknown>;
  
  // 1. Validar el token y el rol (sin cambios)
  if (!(typeof obj.token === 'string' && isRole(obj.role))) {
    return false;
  }

  // 2. Validar el campo 'edificio' (el array)
  const edificios = obj.edificio;
  if (edificios === undefined) {
    // Es válido si es 'undefined'
    return true;
  }

  if (Array.isArray(edificios)) {
    // Si es un array, todos sus elementos deben ser Edificio
    return edificios.every(isEdificio);
  }

  // Si no es undefined y no es un array, es inválido
  return false;
}

function extractRole(scope: string): Role | undefined {
  // Ejemplo de scope: "user-type:Analista"
  const parts = scope.split(':');
  if (parts.length === 2 && parts[0] === 'user-type') {
    // Capitaliza la primera letra y conviértela a minúsculas
    const roleString = parts[1].toLowerCase();
    // Verifica si es un rol válido
    if (isRole(roleString)) {
      return roleString;
    }
  }
  return undefined;
}

class AuthClient {
  async signUp(_: SignUpParams): Promise<{ error?: string }> {
    return { error: 'Sign up not implemented' };
  }

  async signInWithOAuth(_: SignInWithOAuthParams): Promise<{ error?: string }> {
    return { error: 'Social authentication not implemented' };
  }

  async signInWithPassword(params: SignInWithPasswordParams): Promise<{ error?: string }> {
    const { email, password } = params;
    const API_URL = '/login'

    const requestBody = {
      email,
      password,
      device_id: '1',
    };

    try {
        const loginRes = await gs.post(API_URL, requestBody) as { access_token?: string; edificios?: Edificio[] };

        const accessToken = loginRes.access_token;
        const edificios = loginRes.edificios;

        if (!accessToken || !edificios) {
          return {error: 'Error de autenticación: Email o contraseña invalida'};
        }

        const decodedPayload = decodeJwtToken(accessToken);

        if (!decodedPayload) {
            return { error: 'Token JWT inválido o malformado.' };
        }

        const rawRole = decodedPayload.scope; 

        let role: Role | undefined;
        if (typeof rawRole === 'string') {
            role = extractRole(rawRole) || (isRole(rawRole.toLowerCase()) ? rawRole.toLowerCase() as Role : undefined);
        }

        if (!role) {
          return {error: 'Rol de ususairo invalido o no reconocido en el token.'};
        }

        const payload: StoredPayload = {
          token: accessToken,
          edificio: edificios.length > 0 ? edificios : undefined,
          selected_edificio: edificios[0],
          role,
        };

        localStorage.setItem(STORAGE_KEY, JSON.stringify(payload));
        return {};
      } catch(error) {
        return  { error: 'Credenciales invalidas' };
      }
  }

  async resetPassword(_: ResetPasswordParams): Promise<{ error?: string }> {
    return { error: 'No implementado aun' };
  }

  async updatePassword(_: ResetPasswordParams): Promise<{ error?: string }> {
    return { error: 'No implementado aun' };
  }

  async getUser(): Promise<{ data?: SessionUser | null; error?: string }> {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return { data: null };

    try {
      const parsed: unknown = JSON.parse(raw);
      if (!isStoredPayload(parsed)) {
        localStorage.removeItem(STORAGE_KEY);
        return { data: null, error: 'Payload de la sesion invalido' };
      }

      const data: SessionUser = {
        ...baseUser,
        edificio: parsed.edificio,
        role: parsed.role,
        selected_edificio: parsed.selected_edificio,
      };

      return { data };
    } catch {
      localStorage.removeItem(STORAGE_KEY);
      return { data: null, error: 'Invalid session payload' };
    }
  }

  async signOut(): Promise<{ error?: string }> {
    localStorage.removeItem(STORAGE_KEY);
    return {};
  }

  async setSelectedEdificio(newEdificio: Edificio): Promise<{ error?: string }> {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return { error: 'No hay sesión activa para actualizar.' };
    }

    try {
      const parsed: unknown = JSON.parse(raw);
      
      if (!isStoredPayload(parsed)) {
        return { error: 'Payload de sesión inválido.' };
      }
      
      const currentPayload = parsed;

      // Verifica que el newEdificio sea parte del array 'edificio'
      const isAllowed = currentPayload.edificio?.some(e => e.id === newEdificio.id);
      if (!isAllowed) {
         return { error: 'El edificio seleccionado no está en la lista permitida.' };
      }

      currentPayload.selected_edificio = newEdificio;

      localStorage.setItem(STORAGE_KEY, JSON.stringify(currentPayload));

      return {};
    } catch (error) {
      return { error: 'Error al procesar la actualización del edificio.' };
    }
  }
}

export const authClient = new AuthClient();
