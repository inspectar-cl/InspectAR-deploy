package handlers

type EndpointConfig struct {
	Path        string         // Ej: "/api/dashboard"
	Method      string         // Ej: "GET", "POST"
	Handler     string         // Nombre del handler
	Description string         // Descripción del endpoint
	Queries     []ServiceQuery // Consultas a realizar
}

type ServiceQuery struct {
	Name         string            // Nombre identificador (ej: "activos", "sensores")
	TargetEnvVar string            // Variable de entorno del servicio
	Path         string            // Path del endpoint (ej: "/activo", "/sensores")
	Method       string            // Método HTTP
	Headers      map[string]string // Headers adicionales si son necesarios
	Transform    string            // Nombre de la función de transformación (opcional)
	DependsOn    string            // Nombre de la query de la que depende (para queries secuenciales)
	PathTemplate string            // Template del path con variables (ej: "/activo/{activo_id}/sensores")
}

// Configuración de todos los endpoints compuestos
var CompositeEndpoints = []EndpointConfig{
	{
		Path:        "/api/calderas-estado",
		Method:      "GET",
		Handler:     "GetCalderasEstado",
		Description: "Obtiene todas las calderas y el estado de sus sensores",
		Queries: []ServiceQuery{
			{
				Name:         "calderas",
				TargetEnvVar: "GESTION_URL",
				Path:         "/activos/tipo/caldera",
				Method:       "GET",
				Transform:    "TransformCalderas",
			},
			// La segunda query se maneja completamente dentro de TransformSensoresEstado
			// No necesitamos definirla aquí porque se ejecuta dinámicamente
		},
	},
	// Aquí puedes agregar más endpoints fácilmente
	// {
	// 	Path:        "/api/dashboard",
	// 	Method:      "GET",
	// 	Handler:     "GetDashboard",
	// 	Description: "Dashboard principal con datos unificados",
	// 	Queries: []ServiceQuery{
	// 		{
	// 			Name:         "activos",
	// 			TargetEnvVar: "GESTION_URL",
	// 			Path:         "/gestion/activos",
	// 			Method:       "GET",
	// 		},
	// 		{
	// 			Name:         "alertas",
	// 			TargetEnvVar: "NOTIFICATION_URL",
	// 			Path:         "/notificacion/alertas",
	// 			Method:       "GET",
	// 		},
	// 	},
	// },
}
