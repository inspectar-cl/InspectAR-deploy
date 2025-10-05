'use client';

import type { User } from '@/types/user';

type Role = 'analista' | 'tecnico' | 'residente' | 'admin';

interface StoredPayload {
  token: string;
  edificio: string;
  role: Role;
}

interface SessionUser extends User {
  edificio: string;
  role: Role;
}

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

function isStoredPayload(v: unknown): v is StoredPayload {
  if (typeof v !== 'object' || v === null) return false;
  const obj = v as Record<string, unknown>;
  return (
    typeof obj.token === 'string' &&
    typeof obj.edificio === 'string' &&
    isRole(obj.role)
  );
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

    if (email !== 'sofia@devias.io' || password !== 'Secret1') {
      return { error: 'Invalid credentials' };
    }

    const payload: StoredPayload = {
      token: generateToken(),
      edificio: 'EA',
      role: 'residente', // cambiar
    };

    localStorage.setItem(STORAGE_KEY, JSON.stringify(payload));
    return {};
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
