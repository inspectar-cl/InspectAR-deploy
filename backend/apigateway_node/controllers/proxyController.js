const { createProxyMiddleware } = require('http-proxy-middleware');
const services = require('../services/serviceMap');

const proxy = {
  sensores: createProxyMiddleware({
    target: services.sensores,
    changeOrigin: true,
    pathRewrite: { '^/api/sensores': '' }
  }),

  activo: createProxyMiddleware({
    target: services.activo,
    changeOrigin: true,
    pathRewrite: (path, req) => {
      if (req.originalUrl.startsWith('/api/activo')) {
        return '/activo' + req.originalUrl.slice('/api/activo'.length);
      }
      return path;
    }
  }),

  lectura: createProxyMiddleware({
    target: services.lectura,
    changeOrigin: true,
    pathRewrite: (path, req) => {
      // Quita solo el prefijo /api
      if (req.originalUrl.startsWith('/api/lectura')) {
        return req.originalUrl.replace(/^\/api/, '');
      }
      return path;
    }
  }),
};

module.exports = proxy;
