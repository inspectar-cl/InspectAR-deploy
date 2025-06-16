export const paths = {
  home: '/',
  auth: { signIn: '/auth/sign-in', signUp: '/auth/sign-up', resetPassword: '/auth/reset-password' },
  dashboard: {
    inicio: '/dashboard',
    account: '/dashboard/account',
    alertas: '/dashboard/alertas',
    activos: '/dashboard/activos',
    activos2: '/dashboard/activos2',
    informe: '/errors/not-found',
    planos: '/errors/not-found',
    graficos: '/errors/not-found',
    settings: '/dashboard/settings',
    activoDetail: (id: string) => `/dashboard/activos2/${id}`,
  },
  errors: { notFound: '/errors/not-found' },
} as const;
