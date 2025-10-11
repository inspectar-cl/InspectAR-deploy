import axios from 'axios'

// Preferimos rutas relativas para aprovechar el reverse proxy de Next (/api -> localhost:3500)
const BASE_URL = process.env.NEXT_PUBLIC_API_GATEWAY_URL || '/api'

const apiClient = axios.create({
    baseURL: BASE_URL,
    timeout: 30000,
    headers: {
        "ngrok-skip-browser-warning": "true",
    },
    // No seteamos Content-Type globalmente, lo hace axios según el body
})

// TODO: Agregar los demás métodos (-Vixo 14/06).

export default class Services {
    async get(uri, params) {
        return await apiClient.get(uri, params ? {params} : {})
            .then(res => res.data)
            .catch(error => ({
                mensaje: "error inesperado",
                error: error.response
            }))
    }

    async authorizedGet(uri, accessToken) {
        const config = {
        headers: {
            'Authorization': `Bearer ${accessToken}`
            }
        };

        return await apiClient.get(uri, config)
            .then(res => res.data)
            .catch(error => ({
                mensaje: "error inesperado",
                error: error.response
        }))
    }

    async authorizedPost(uri, data, accessToken) {
    // Configuración de headers, incluyendo la autorización
        const config = {
            headers: {
                'Authorization': `Bearer ${accessToken}`,
                'Content-Type': 'application/json' // Aseguramos el tipo JSON para el body
            }
        };

        return await apiClient.post(uri, data, config)
            .then(res => res.data)
            .catch(error => ({
                mensaje: "error inesperado",
                error: error.response
        }))
    }

    async post(uri, data) {
        // Si es FormData, no seteamos Content-Type, axios lo hace solo
        const isFormData = (typeof FormData !== 'undefined') && data instanceof FormData;
        const config = isFormData ? {} : { headers: { 'Content-Type': 'application/json' } };
        return await apiClient.post(uri, data, config)
            .then(res => res.data)
            .catch(error => ({
                mensaje: "error inesperado",
                error: error.response
            }))
    }

    async put(uri, data) {
        return await apiClient.put(uri, data)
            .then(res => res.data)
            .catch(error => ({
                mensaje: "error inesperado",
                error: error.response
            }))
    }

    async delete(uri) {
        return await apiClient.delete(uri)
            .then(res => res.data)
            .catch(error => ({
                mensaje: "error inesperado",
                error: error.response
            }))
    }
}
