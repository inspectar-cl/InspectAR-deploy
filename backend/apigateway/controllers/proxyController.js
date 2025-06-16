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

  activo: createProxyMiddleware({
    target: services.activo,
    changeOrigin: true,
    pathRewrite: (path, req) => {
      console.log('Path original:', path);
      return '/activo';
    },
    onError(err, req, res) {
      console.error('[ERROR] No se pudo conectar al servicio de activo:', err.message);
      res.status(520).json({
        error: true,
        message: 'Servicio de activo no disponible.'
      })
    }
  }),
};

module.exports = proxy;
