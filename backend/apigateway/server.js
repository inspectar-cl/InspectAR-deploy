const dotenv = require("dotenv"); 
require('dotenv').config();
const express = require("express");
const cors = require("cors");
const morgan = require("morgan");
const bodyParser = require("body-parser");

const app = express();
app.use(morgan("dev"));
app.use(express.json());
app.use(cors({
  // origin: process.env.FRONT_URL // en caso que queramos aceptar solo peticiones desde la ip del front. -Vixo 14/06
}));

const proxyRoutes = require('./routes/proxyRoutes')
app.use('/api', proxyRoutes)
app.use((err, req, res, next) => {
  console.error('[ERROR]', err.message);
  if (!res.headersSent) {
    res.status(520).json({
      error: true,
      message: 'Servicio no disponible.'
    });
  }
});

app.use((req, res, next) => {
  res.status(404).json({
    error: true,
    message: 'Ruta no encontrada.'
  });
});

// Start server
const port = process.env.PORT; // || 3000;
app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});
