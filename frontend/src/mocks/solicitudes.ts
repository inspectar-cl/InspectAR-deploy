import type { SolicitudAPI } from '@/types/form-solicitud'; // Ajusta la ruta a tus tipos

export const MOCK_SOLICITUDES: SolicitudAPI[] = [
  {
    id: 1,
    estado: 'Pendiente',
    fechaCreacion: '2025-10-30T10:00:00Z',
    tipoSolicitud: 'Edificio',
    asunto: 'Instalar nueva unidad de A/C en el techo',
    detalles: 'La unidad actual está fallando y necesita reemplazo antes del verano.',
    datosEspecificos: {
      nombre: 'Edificio Central',
      direccion: 'Av. Principal 123',
      latitud: -33.44889,
      longitud: -70.66926,
    }
  },
  {
    id: 2,
    estado: 'EnProgreso',
    fechaCreacion: '2025-10-29T14:30:00Z',
    tipoSolicitud: 'Activo',
    asunto: 'Revisión de bomba de agua Sótano 2',
    detalles: 'Se reporta un ruido extraño al operar.',
    datosEspecificos: {
      tipoActivo: 'BombaDeAgua',
      edificioId: 1,
      ubicacion: 'Sótano 2, sala de máquinas',
      descripcion: 'Bomba principal del sistema de agua potable.',
      imagen: 'https://via.placeholder.com/150' // URL de imagen de ejemplo
    }
  },
  {
    id: 4,
    estado: 'Pendiente',
    fechaCreacion: '2025-10-31T11:00:00Z',
    tipoSolicitud: 'Técnico',
    asunto: 'Asignar técnico a revisión de panel eléctrico',
    detalles: 'Se necesita un especialista eléctrico para el panel del piso 5.',
    datosEspecificos: {
      nombre: 'Juan Pérez',
      especialidad: 'Eléctrico',
      correo: 'juan.perez@tecnico.com',
      telefono: '+56912345678',
      activosAsociados: [201, 202]
    }
  },
  {
    id: 5,
    estado: 'Pendiente',
    fechaCreacion: '2025-10-31T12:00:00Z',
    tipoSolicitud: 'Edificio',
    asunto: 'Fuga en cañería piso 3',
    detalles: 'Goteo constante en el baño de hombres.',
    datosEspecificos: {
      nombre: 'Edificio Anexo',
      direccion: 'Calle Falsa 456',
      latitud: -33.45000,
      longitud: -70.67000,
    }
  }
];