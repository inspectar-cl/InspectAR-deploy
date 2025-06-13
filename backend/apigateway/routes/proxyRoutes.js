const express = require('express');
const proxy = require('../controllers/proxyController');
const router = express.Router();

router.use('/sensores', proxy.sensores);
// router.use('/products', proxy.productService);

module.exports = router;
