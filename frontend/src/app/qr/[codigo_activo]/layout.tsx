import * as React from 'react';

interface LayoutProps {
  children: React.ReactNode;
}

export default function QRLayout({ children }: LayoutProps): React.JSX.Element {
  // Este layout NO requiere autenticación - es público
  return <>{children}</>;
}
