'use client';

import * as React from 'react';
import { usePathname, useRouter } from 'next/navigation';
import Alert from '@mui/material/Alert';

import { paths } from '@/paths';
import { logger } from '@/lib/default-logger';
import { useUser } from '@/hooks/use-user';

// 👇 importa tu mapa real de permisos
import { rolePermissions } from '@/rolePermissions'; // <-- ajusta la ruta si difiere
type Role = keyof typeof rolePermissions;

// Helpers locales (sin archivo extra)
function normalize(p: string) {
  try {
    const u = new URL(p, 'http://x');
    p = u.pathname;
  } catch { /* pathname simple */ }
  if (p.length > 1 && p.endsWith('/')) p = p.slice(0, -1);
  return p;
}
function isPathAllowed(role: Role, pathname: string): boolean {
  const allowed = rolePermissions[role] ?? [];
  const path = normalize(pathname);
  const dashRoot = normalize(paths.dashboard.inicio); // "/dashboard"

  return allowed.some((prefix) => {
    const pref = normalize(prefix);
    if (pref === dashRoot) {
      //el inicio del dashboard NO habilita subrutas
      return path === pref;
    }
    // para el resto, sí permitimos subrutas
    return path === pref || path.startsWith(pref + '/');
  });
}
function firstAllowed(role: Role): string {
  const allowed = rolePermissions[role] ?? [];
  return allowed[0] ?? paths.dashboard.inicio;
}

export interface AuthGuardProps {
  children: React.ReactNode;
}

export function AuthGuard({ children }: AuthGuardProps): React.JSX.Element | null {
  const router = useRouter();
  const pathname = usePathname();
  const { user, error, isLoading } = useUser();
  const [isChecking, setIsChecking] = React.useState<boolean>(true);

  const checkPermissions = async (): Promise<void> => {
    if (isLoading) return;

    if (error) {
      setIsChecking(false);
      return;
    }

    if (!user) {
      logger.debug('[AuthGuard]: User is not logged in, redirecting to sign in');
      router.replace(paths.auth.signIn);
      return;
    }

    // 🔒 Gateo por rol y ruta
    const role = ((user as any).role ?? 'residente') as Role;

    // solo si estamos dentro de /dashboard gateamos por permisos (opcional)
    if (pathname?.startsWith('/dashboard') && !isPathAllowed(role, pathname)) {
      const fallback = firstAllowed(role);
      logger.debug('[AuthGuard]: Path not allowed for role, redirecting to', fallback);
      router.replace(fallback);
      return;
    }

    setIsChecking(false);
  };

  React.useEffect(() => {
    checkPermissions().catch(() => { /* noop */ });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [user, error, isLoading, pathname]); // 👈 importante: depende de pathname

  if (isChecking) return null;

  if (error) return <Alert color="error">{error}</Alert>;

  return <>{children}</>;
}
