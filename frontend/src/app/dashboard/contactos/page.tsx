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
  const [especialidades, setEspecialidades] = useState<string[]>([]);

  useEffect(() => {
    const fetchContactos = async () => {
      try {
        const response = await gs.get("/gestion/tecnicos");
        console.log("response", response)
        const contactosMapeados = Array.isArray(response)
          ? response.map((c: any) => ({
              id: c.id,
              name: c.nombre,
              email: "j3291674@gmail.com",
              phone: c.telefono,
              especialidad: c.especialidad,
              avatar: '/assets/avatar-8.png',
            }))
          : [];

        // Para cada técnico, busca sus activos asociados
        const contactosConActivos = await Promise.all(
          contactosMapeados.map(async (tecnico) => {
            const activos = await gs.get(`/gestion/activos-de-tecnico/${tecnico.id}`);
            // Mapea los activos a la estructura esperada por Contacto
            const activosMapeados = Array.isArray(activos)
              ? activos.map((a: any) => ({
                  id_activo: a.id_activo ?? a.id ?? "",
                  nombre_activo: a.nombre_activo ?? a.nombre ?? "",
                }))
              : [];
            return { ...tecnico, activo: activosMapeados };
          })
        );

        setContactos(contactosConActivos);

        // Extrae especialidades únicas
        const especialidadesUnicas = [
          ...new Set(contactosConActivos.map((c) => c.especialidad).filter(Boolean)),
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
