# InspectAR

## Requisitos
- Docker y Docker Compose (para el backend)
- Node.js y npm instalados en WSL (para el frontend)

## Ejecución

### Backend
```bash
make run-b
```

### Frontend
Primero instala las dependencias (solo la primera vez):
```bash
cd frontend && npm install
```

Luego ejecuta:
```bash
make run-front
```

### Detener Backend
```bash
make stop-b
```

**Nota**: Si usas WSL, asegúrate de tener Node.js instalado directamente en WSL, no en Windows.