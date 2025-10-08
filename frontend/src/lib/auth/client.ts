'use client';

import type { User } from '@/types/user';
import Services from '@/modules/Services';
import { accessToken } from 'mapbox-gl';

type Role = 'analista' | 'tecnico' | 'residente' | 'admin';

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
}

interface SessionUser extends User {
  edificio: Edificio[] | undefined;
  role: Role;
}

const gs = new Services();

const STORAGE_KEY = 'custom-auth-token';

function generateToken(): string {
  const arr = new Uint8Array(12);
  window.crypto.getRandomValues(arr);
  return Array.from(arr, (v) => v.toString(16).padStart(2, '0')).join('');
}

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
  return v === 'analista' || v === 'tecnico' || v === 'residente' || v === 'admin';
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
      return roleString as Role;
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
    const API_ROL = '/rol'

    const requestBody = {
      email: email,
      password: password,
      device_id: '1',
    };

    try {
        const loginRes = await gs.post(API_URL, requestBody) as { access_token?: string; edificios?: Edificio[] };

        const accessToken = loginRes.access_token;
        const edificios = loginRes.edificios;

        if (!accessToken || !edificios) {
          return {error: 'Error de autenticación: Email o contraseña invalida'};
        }

        const rolRes = await gs.authorizedGet(API_ROL, accessToken) as { scope?: string, error?: any };

        if (rolRes.error) {
          return { error: 'Error al obtener el rol del usuario.' };
        }

        const scope = rolRes.scope;
        const role = scope ? extractRole(scope) : undefined;

        if (!role) {
          return {error: 'Rol de ususairo invalido o no reconocido'};
        }

        const payload: StoredPayload = {
          token: accessToken,
          edificio: edificios.length > 0 ? edificios : undefined,
          role: role,
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
        return { data: null, error: 'Invalid session payload' };
      }

      const data: SessionUser = {
        ...baseUser,
        edificio: parsed.edificio,
        role: parsed.role,
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
}

export const authClient = new AuthClient();
