import React from 'react';

import ActivosTable from '@/components/dashboard/activos/ActivosTable';

export default function ActivosPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Lista de Activos</h1>
      <ActivosTable />
    </div>
  )
}
