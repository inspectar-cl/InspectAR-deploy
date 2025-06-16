# Imagen base oficial de Node.js
FROM node:20

# Establece el directorio de trabajo
WORKDIR /usr/src/app

# Copia los archivos de dependencias
COPY package*.json ./

# Instala las dependencias
RUN npm install

# Copia el resto del proyecto
COPY . .

# Expone el puerto de tu API Gateway (ajústalo si es diferente)
EXPOSE 3000

# Comando para arrancar tu API Gateway
CMD ["node", "server.js"]
