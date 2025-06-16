const express = require('express');
const proxy = require('../controllers/proxyController');
const router = express.Router();

router.use('/sensores', proxy.sensores);
router.use('/activo', (req, res, next) => {
  console.log('Ruta recibida en proxy:', req.originalUrl);
  next();
}, proxy.activo);
// router.use('/products', proxy.productService);

module.exports = router;
