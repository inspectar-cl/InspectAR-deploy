import * as React from 'react';
import type { Metadata } from 'next';
import { config } from '@/config';
import type { Contacto } from '@/components/dashboard/contactos/contactos-table';
import { ContactosClient } from '@/components/dashboard/contactos/contactos-client';

export const metadata = { title: `Customers | Dashboard | ${config.site.name}` } satisfies Metadata;

const contactos = [
  {
    id: 'USR-010',
    id_edificio: 'EDF-001',
    name: 'Alcides Antonio',
    avatar: '/assets/avatar-10.png',
    email: 'alcides.antonio@devias.io',
    phone: '908-691-3242',
    activo: [
      { id_activo: 'ACT-001', nombre_activo: 'Bomba de Agua 1' }
    ],
    especialidad: 'Bomba de Agua',
  },
  {
    id: 'USR-009',
    id_edificio: 'EDF-001',
    name: 'Marcus Finn',
    avatar: '/assets/avatar-9.png',
    email: 'marcus.finn@devias.io',
    phone: '415-907-2647',
    activo: [
      { id_activo: 'ACT-003', nombre_activo: 'Bomba de Agua 2' },
      { id_activo: 'ACT-004', nombre_activo: 'Bomba de Agua 3' },
      { id_activo: 'ACT-002', nombre_activo: 'Ascensor 1' },
      { id_activo: 'ACT-001', nombre_activo: 'Ascensor 1' }
    ],
    especialidad: 'Bomba de Agua',
  },
  {
    id: 'USR-008',
    id_edificio: 'EDF-001',
    name: 'Jie Yan',
    avatar: '/assets/avatar-8.png',
    email: 'jie.yan.song@devias.io',
    phone: '770-635-2682',
    activo: [
      { id_activo: 'ACT-005', nombre_activo: 'Ascensor 2' },
      { id_activo: 'ACT-006', nombre_activo: 'Ascensor 3' }
    ],
    especialidad: 'Ascensor',
  },
  {
    id: 'USR-007',
    id_edificio: 'EDF-001',
    name: 'Nasimiyu Danai',
    avatar: '/assets/avatar-7.png',
    email: 'nasimiyu.danai@devias.io',
    phone: '801-301-7894',
    activo: [
      { id_activo: 'ACT-007', nombre_activo: 'Panel Electrico 1' },
      { id_activo: 'ACT-008', nombre_activo: 'Panel Electrico 2' }
    ],
    especialidad: 'Sistema Eléctrico',
  },
  {
    id: 'USR-006',
    id_edificio: 'EDF-001',
    name: 'Iulia Albu',
    avatar: '/assets/avatar-6.png',
    email: 'iulia.albu@devias.io',
    phone: '313-812-8947',
    activo: [
      { id_activo: 'ACT-009', nombre_activo: 'Bomba de Agua 5' },
      { id_activo: 'ACT-010', nombre_activo: 'Panel Electrico 3' }
    ],
    especialidad: 'Sistema Eléctrico',
  },
  {
    id: 'USR-005',
    id_edificio: 'EDF-001',
    name: 'Fran Perez',
    avatar: '/assets/avatar-5.png',
    email: 'fran.perez@devias.io',
    phone: '712-351-5711',
    activo: [
      { id_activo: 'ACT-011', nombre_activo: 'Ascensor 5' },
      { id_activo: 'ACT-012', nombre_activo: 'Panel Electrico 4' }
    ],
    especialidad: 'Ascensor',
  },

  {
    id: 'USR-004',
    id_edificio: 'EDF-003',
    name: 'Penjani Inyene',
    avatar: '/assets/avatar-4.png',
    email: 'penjani.inyene@devias.io',
    phone: '858-602-3409',
    activo: [
      { id_activo: 'ACT-013', nombre_activo: 'Ascensor 6' },
      { id_activo: 'ACT-014', nombre_activo: 'Panel Electrico 5' }
    ],
    especialidad: 'Sistema Eléctrico',
  },
  {
    id: 'USR-003',
    id_edificio: 'EDF-003',
    name: 'Carson Darrin',
    avatar: '/assets/avatar-3.png',
    email: 'carson.darrin@devias.io',
    phone: '304-428-3097',
    activo: [
      { id_activo: 'ACT-015', nombre_activo: 'Bomba de Agua 6' },
      { id_activo: 'ACT-016', nombre_activo: 'Panel Electrico 6' }
    ],
    especialidad: 'Bomba de Agua',
  },
  {
    id: 'USR-002',
    id_edificio: 'EDF-002',
    name: 'Siegbert Gottfried',
    avatar: '/assets/avatar-2.png',
    email: 'siegbert.gottfried@devias.io',
    phone: '702-661-1654',
    activo: [
      { id_activo: 'ACT-017', nombre_activo: 'Bomba de Agua 7' },
      { id_activo: 'ACT-018', nombre_activo: 'Ascensor 6' }
    ],
    especialidad: 'Ascensor',
  },
  {
    id: 'USR-001',
    id_edificio: 'EDF-002',
    name: 'Miron Vitold',
    avatar: '/assets/avatar-1.png',
    email: 'miron.vitold@devias.io',
    phone: '972-333-4106',
    activo: [
      { id_activo: 'ACT-019', nombre_activo: 'Bomba de Agua 7' },
      { id_activo: 'ACT-020', nombre_activo: 'Panel Electrico 7' }
    ],
    especialidad: 'Bomba de Agua',
  },
] satisfies Contacto[];

export default function Page(): React.JSX.Element {
  return <ContactosClient contactos={contactos} />
}
