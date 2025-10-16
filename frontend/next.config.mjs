/** @type {import('next').NextConfig} */
const config = {
    reactStrictMode: false,
    // Suprimir warnings de hidratación causados por extensiones del navegador
    onDemandEntries: {
        // period (in ms) where the server will keep pages in the buffer
        maxInactiveAge: 25 * 1000,
        // number of pages that should be kept simultaneously without being disposed
        pagesBufferLength: 2,
    },
    // Permite acceder al dev server desde dominios ngrok en desarrollo
    allowedDevOrigins: ['*.ngrok-free.app'],
    async rewrites() {
        return [
            {
                source: '/api/:path*',
                destination: 'http://localhost:3500/api/:path*',
            },
        ];
    },
};

export default config;
