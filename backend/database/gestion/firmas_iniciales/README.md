# Firmas Iniciales para Testing

## Instrucciones

1. **Coloca tus 4 archivos JPG de firmas aquí** con los siguientes nombres:
   - `firma_usuario_1.jpg` (para el usuario analista)
   - `firma_usuario_2.jpg` (para el usuario tecnico)
   - `firma_usuario_3.jpg` (para el usuario residente)
   - `firma_usuario_4.jpg` (para el usuario admin)

2. Las firmas deben ser imágenes JPG de preferencia con:
   - Dimensiones recomendadas: 300x150 px (ancho x alto)
   - Fondo transparente o blanco
   - Tamaño menor a 500KB

3. Estas firmas serán copiadas al contenedor de gestion-service durante la inicialización

## Archivos requeridos:
- [ ] firma_usuario_1.jpg
- [ ] firma_usuario_2.jpg
- [ ] firma_usuario_3.jpg
- [ ] firma_usuario_4.jpg

## Alternativa: Usar firmas de ejemplo

Si no tienes firmas reales, puedes crear archivos de texto simples que serán convertidos a imágenes:

```bash
# Crear firmas de ejemplo con ImageMagick (si está instalado)
convert -size 300x150 xc:white -pointsize 24 -draw "text 50,80 'Analista Firma'" firma_usuario_1.jpg
convert -size 300x150 xc:white -pointsize 24 -draw "text 50,80 'Técnico Firma'" firma_usuario_2.jpg
convert -size 300x150 xc:white -pointsize 24 -draw "text 50,80 'Residente Firma'" firma_usuario_3.jpg
convert -size 300x150 xc:white -pointsize 24 -draw "text 50,80 'Admin Firma'" firma_usuario_4.jpg
```

O simplemente copia cualquier imagen JPG 4 veces con estos nombres.
