'use client'

import * as React from 'react';
import type { Contacto } from '@/components/dashboard/contactos/contactos-table';
import { ContactosClient } from '@/components/dashboard/contactos/contactos-client';

// export const metadata = { title: `Customers | Dashboard | ${config.site.name}` } satisfies Metadata;

import { useEffect, useState } from 'react';
import Services from '@/modules/Services'
const gs = new Services()

export default function Page(): React.JSX.Element {
  // const [contactos, setContactos] = useState<Contacto[]>([]);

  // useEffect(() => {
  //   const fetchContactos = async () => {
  //     try {
  //       const response = await gs.get("/gestion/tecnicos");
  //       console.log("Respuesta de /gestion/tecnicos:", response);
  //       // Mapea los campos del backend a los del frontend
  //       const contactosMapeados = Array.isArray(response)
  //         ? response.map((c: any) => ({
  //             id: c.id,
  //             name: c.nombre,
  //             email: c.email,
  //             phone: c.telefono,
  //             especialidad: c.especialidad,
  //             avatar: '/assets/avatar-8.png',
  //           }))
  //         : [];
  //       setContactos(contactosMapeados);
  //     } catch (err) {
  //       console.error("Error al obtener los técnicos", err);
  //       setContactos([]);
  //     }
  //   };

  //   fetchContactos();
  // }, []);

  const [contactos, setContactos] = useState<Contacto[]>([]);
  const [especialidades, setEspecialidades] = useState<string[]>([]);

  useEffect(() => {
    const fetchContactos = async () => {
      try {
        const response = await gs.get("/gestion/tecnicos");
        const contactosMapeados = Array.isArray(response)
          ? response.map((c: any) => ({
              id: c.id,
              name: c.nombre,
              email: c.email,
              phone: c.telefono,
              especialidad: c.especialidad,
              avatar: '/assets/avatar-8.png',
            }))
          : [];
        setContactos(contactosMapeados);

        // Extrae especialidades únicas
        const especialidadesUnicas = [
          ...new Set(contactosMapeados.map((c) => c.especialidad).filter(Boolean)),
        ];
        setEspecialidades(especialidadesUnicas);
      } catch (err) {
        setContactos([]);
        setEspecialidades([]);
      }
    };

    fetchContactos();
  }, []);

  return <ContactosClient contactos={contactos} especialidades={especialidades} />;
}
