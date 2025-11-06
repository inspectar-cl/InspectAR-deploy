'use client';

import * as React from 'react';

import type { User } from '@/types/user';
import { authClient } from '@/lib/auth/client';
import { logger } from '@/lib/default-logger';

type Role = 'analista' | 'tecnico' | 'residente' | 'root' | 'administrador';

interface Edificio {
  id: number;
  nombre: string;
  direccion: string;
  creado_en: string;
}

interface SessionUser extends User {
  edificio: Edificio[] | undefined;
  role: Role;
  selected_edificio: Edificio;
}

export interface UserContextValue {
  user: SessionUser | null;
  error: string | null;
  isLoading: boolean;
  checkSession?: () => Promise<void>;
  changeSelectedEdificio: (newEdificio: Edificio) => Promise<void>;
}

export const UserContext = React.createContext<UserContextValue | undefined>(undefined);

export interface UserProviderProps {
  children: React.ReactNode;
}

export function UserProvider({ children }: UserProviderProps): React.JSX.Element {
  const [state, setState] = React.useState<{ user: SessionUser | null; error: string | null; isLoading: boolean }>({
    user: null,
    error: null,
    isLoading: true,
  });

  const checkSession = React.useCallback(async (): Promise<void> => {
    try {
      const { data, error } = await authClient.getUser(); 

      if (error) {
        logger.error(error);
        setState((prev) => ({ ...prev, user: null, error: 'Something went wrong', isLoading: false }));
        return;
      }

      setState((prev) => ({ ...prev, user: data ?? null, error: null, isLoading: false }));
    } catch (err) {
      logger.error(err);
      setState((prev) => ({ ...prev, user: null, error: 'Something went wrong', isLoading: false }));
    }
  }, []);

  const changeSelectedEdificio = React.useCallback(async (newEdificio: Edificio): Promise<void> => {
    const { error } = await authClient.setSelectedEdificio(newEdificio); 
    
    if (error) {
        logger.error('Error al cambiar edificio:', error);
        return;
    }

    setState((prevState) => {
        if (!prevState.user) return prevState;
        return {
            ...prevState,
            user: {
                ...prevState.user,
                selected_edificio: newEdificio,
            },
        };
    });
  }, []);

  React.useEffect(() => {
    checkSession().catch((err: unknown) => {
      logger.error(err);
      // noop
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps -- Expected
  }, []);

  const contextValue = React.useMemo(() => ({
    ...state,
    checkSession,
    changeSelectedEdificio,
  }), [state, checkSession, changeSelectedEdificio]);

  return <UserContext.Provider value={contextValue}>{children}</UserContext.Provider>;
}

export function useAuthUser(): UserContextValue {
  const context = React.useContext(UserContext);

  if (context === undefined) {
    throw new Error('useAuthUser must be used within a UserProvider');
  }

  return context;
}

export const UserConsumer = UserContext.Consumer;