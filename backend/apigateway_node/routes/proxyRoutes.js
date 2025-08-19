const express = require('express');
const proxy = require('../controllers/proxyController');
const router = express.Router();

router.use('/sensores', proxy.sensores);
router.use('/activo', (req, res, next) => {
  console.log('originalUrl:', req.originalUrl, 'path:', req.path);
  next();
}, proxy.activo);
// router.use('/products', proxy.productService);
router.use('/lectura', proxy.lectura)

module.exports = router;
