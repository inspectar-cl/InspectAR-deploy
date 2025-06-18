export const paths = {
  home: '/',
  auth: { signIn: '/auth/sign-in', signUp: '/auth/sign-up', resetPassword: '/auth/reset-password' },
  dashboard: {
    inicio: '/dashboard',
    account: '/dashboard/account',
    alertas: '/dashboard/alertas',
    activos: '/dashboard/activos',
    informe: '/errors/not-found',
    planos: '/errors/not-found',
    graficos: '/dashboard/graficos',
    settings: '/dashboard/settings',
    activoDetail: (id: string) => `/dashboard/activos/${id}`,
  },
  errors: { notFound: '/errors/not-found' },
} as const;
