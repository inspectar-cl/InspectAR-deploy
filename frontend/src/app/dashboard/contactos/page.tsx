'use client'

import * as React from 'react';
import type { Contacto } from '@/components/dashboard/contactos/contactos-table';
import { ContactosClient } from '@/components/dashboard/contactos/contactos-client';

// export const metadata = { title: `Customers | Dashboard | ${config.site.name}` } satisfies Metadata;

import { useEffect, useState } from 'react';
import Services from '@/modules/Services'
const gs = new Services()

export default function Page(): React.JSX.Element {
  const [contactos, setContactos] = useState<Contacto[]>([]);

  useEffect(() => {
    const fetchContactos = async () => {
      try {
        const response = await gs.get("/gestion/tecnicos");
        console.log("Respuesta de /gestion/tecnicos:", response);
        setContactos(response ?? []);
      } catch (err) {
        console.error("Error al obtener los técnicos", err);
      }
    };

    fetchContactos();
  }, []);

  return <ContactosClient contactos={contactos} />
}
