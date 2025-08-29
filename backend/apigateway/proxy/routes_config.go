package proxy

type ProxyRoute struct {
    PathPrefix   string // Ej: "/api/gestion"
    TargetEnvVar string // Ej: "GESTION_URL"
    PrependPath  string // Ej: "" o "/algo"
}

var ProxyRoutes = []ProxyRoute{
	// Formato de las rutas, para agregarlas más fácilmente:
	// {
	// 	PathPrefix:   "/api/usuarios",
	// 	TargetEnvVar: "USUARIOS_URL",
	// 	PrependPath:  "",
	// },
    {
        PathPrefix:   "/api/gestion",
        TargetEnvVar: "GESTION_URL",
        PrependPath:  "",
    },
}

