import axios from 'axios'

// Preferimos rutas relativas para aprovechar el reverse proxy de Next (/api -> localhost:3500)
const BASE_URL = process.env.NEXT_PUBLIC_API_GATEWAY_URL || '/api'

const apiClient = axios.create({
    baseURL: BASE_URL,
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json',
    },
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

    async post(uri, data) {
        return await apiClient.post(uri, data)
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
