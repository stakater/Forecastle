package annotations

const (
	// IngressClassAnnotation const used for checking ingress class
	IngressClassAnnotation = "kubernetes.io/ingress.class"
	// ForecastleIconAnnotation const used for forecastle icon
	ForecastleIconAnnotation = "forecastle.stakater.com/icon"
	// ForecastleExposeAnnotation const used for checking whether an ingress is exposed to forecastle
	ForecastleExposeAnnotation = "forecastle.stakater.com/expose"
	// ForecastleAppNameAnnotation const used for overriding the name of the app
	ForecastleAppNameAnnotation = "forecastle.stakater.com/appName"
	// ForecastleGroupAnnotation const used for overriding group
	ForecastleGroupAnnotation = "forecastle.stakater.com/group"
)
