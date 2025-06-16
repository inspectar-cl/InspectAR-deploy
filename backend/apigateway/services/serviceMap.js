console.log("URL sensores:", process.env.SENSORES_URL);
console.log("URL activo:", process.env.ACTIVE_URL);

module.exports = {
    sensores: process.env.SENSORES_URL,
    activo: process.env.ACTIVE_URL,
    // userService_2: process.env.USER_SERVICE_URL_2
}