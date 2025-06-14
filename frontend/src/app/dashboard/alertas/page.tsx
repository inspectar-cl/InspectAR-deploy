import React from 'react';

import AlertasTable from '@/components/dashboard/alertas/AlertasTable';

export default function ActivosPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Alertas</h1>
      <AlertasTable />
    </div>
  )
}
