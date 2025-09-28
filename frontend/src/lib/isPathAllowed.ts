import { rolePermissions } from '@/rolePermissions';
type Role = keyof typeof rolePermissions;

/** Normaliza quitando slash final y query */
function normalize(p: string) {
  try {
    const u = new URL(p, 'http://x');
    p = u.pathname;
  } catch { /* pathname simple */ }
  if (p.length > 1 && p.endsWith('/')) p = p.slice(0, -1);
  return p;
}

/** Devuelve true si el path es igual a un permitido o cuelga de él (prefijo + '/') */
export function isPathAllowed(role: Role, pathname: string): boolean {
  const allowed = rolePermissions[role] ?? [];
  const path = normalize(pathname);
  return allowed.some((prefix) => {
    const pref = normalize(prefix);
    return path === pref || path.startsWith(pref + '/');
  });
}

export function firstAllowed(role: Role): string {
  const allowed = rolePermissions[role] ?? ['/'];
  return allowed[0] ?? '/';
}
