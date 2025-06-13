const { createProxyMiddleware } = require('http-proxy-middleware');
const services = require('../services/serviceMap');

const proxy = {
  sensores: createProxyMiddleware({
    target: services.sensores,
    changeOrigin: true,
    pathRewrite: { '^/api/sensores': '' },
    onError(err, req, res) {
      console.error('[ERROR] No se pudo conectar al servicio de sensores:', err.message);
      res.status(520).json({
        error: true,
        message: 'Servicio de sensores no disponible.'
      })
    }
  }),


  // If you need to add an onError handler, define it as a function property like this:
  // onError: function(err, req, res) {
  //   res.status(500).json({ error: 'Proxy error', details: err.message });
  // },



//   productService: createProxyMiddleware({
//     target: services.productService,
//     changeOrigin: true,
//     pathRewrite: { '^/api/products': '' }
//   })
};

module.exports = proxy;
