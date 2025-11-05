package proxy

type ProxyRoute struct {
    PathPrefix         string   // Ej: "/api/gestion"
    TargetEnvVar       string   // Ej: "GESTION_URL"
    PrependPath        string   // Ej: "" o "/algo"
    Protected          bool     // Si requiere autenticación
    RequiredScopes     []string // Root siempre incluido automáticamente. Otros roles: Analista, Tecnico, Residente, Admin
    RequiredValidation string   // Puede ser "edificio", "activo" o vacío si no se requiere validación
}

var ProxyRoutes = []ProxyRoute{
	// Formato de las rutas, para agregarlas más fácilmente:
	// {
	// 	PathPrefix:   "/api/usuarios",
	// 	TargetEnvVar: "USUARIOS_URL",
	// 	PrependPath:  "",
	// }, // Muy importante esta coma
    {
        PathPrefix:   "/api/gestion",
        TargetEnvVar: "GESTION_URL",
        PrependPath:  "",
    },
    {
        PathPrefix:   "/api/notificacion",
        TargetEnvVar: "NOTIFICATION_URL",
        PrependPath:  "",
    },
    {
        PathPrefix:   "/api/parser",
        TargetEnvVar: "PARSER_URL",
        PrependPath:  "",
    },
    {
        PathPrefix:   "/api/documentacion",
        TargetEnvVar: "DOCUMENTATION_URL",
        PrependPath:  "/api/v1",
    },
    {
        PathPrefix:  "/api/ML",
        TargetEnvVar: "ML_URL",
        PrependPath:  "",
    },
}

