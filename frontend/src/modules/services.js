import axios from 'axios'

console.log(process.env.NEXT_PUBLIC_API_GATEWAY_URL)

const apiClient = axios.create ({
    baseURL: process.env.NEXT_PUBLIC_API_GATEWAY_URL,
    timeout: 5000,
    headers: {
        'Content-Type' : 'application/json'
    }
})

// TODO: Agregar los demás métodos (-Vixo 14/06).

export default class Services {
    async get(uri, item) {
        if (item) {
            return await apiClient.get(uri, {params: item})
                .then(res => res.data)
                .catch(error => ({
                    mensaje: "error inesperado",
                    error: error.response
                }))
        } else {
            return await apiClient.get(uri)
                .then(res => res.data)
                .catch(error => ({
                    mensaje: "error inesperado",
                    error: error.response
                }))
        }
    }
}
