import * as React from 'react';

interface Edificio {
  id: number;
  nombre: string;
  direccion: string;
  creado_en: string;
}

interface AuthInfo {
    role: 'residente' | 'root' | 'analista' | 'tecnico' | null;
    edificio: Edificio[];
    token: string;
}

export const useUserToken = (): {
    user: AuthInfo | null;
    isLoading: boolean;
    error: string | null;
} => {
    const [user, setUser] = React.useState<AuthInfo | null>(null);
    const [isLoading, setIsLoading] = React.useState(true);
    const [error, setError] = React.useState<string | null>(null);

    React.useEffect(() => {
        try {
            const tokenString = localStorage.getItem('custom-auth-token');
            if (tokenString) {
                const userInfo = JSON.parse(tokenString) as unknown as AuthInfo;
                
                setUser(userInfo);
            }
        } catch (e) {
            setError("Error al cargar la sesión. Por favor, inicia sesión de nuevo.");
        } finally {
            setIsLoading(false);
        }
    }, []);

    return { user, isLoading, error };
};