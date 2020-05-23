package annotations

const (
	// ForecastleIconAnnotation const used for forecastle icon
	ForecastleIconAnnotation = "forecastle.stakater.com/icon"
	// ForecastleExposeAnnotation const used for checking whether an ingress is exposed to forecastle
	ForecastleExposeAnnotation = "forecastle.stakater.com/expose"
	// ForecastleAppNameAnnotation const used for overriding the name of the app
	ForecastleAppNameAnnotation = "forecastle.stakater.com/appName"
	// ForecastleGroupAnnotation const used for overriding group
	ForecastleGroupAnnotation = "forecastle.stakater.com/group"
	// ForecastleInstanceAnnotation const used for defining which instance of forecastle to use
	ForecastleInstanceAnnotation = "forecastle.stakater.com/instance"
	// ForecastleNetworkRestrictedAnnotation const used for specifying whether the app is network restricted or not
	ForecastleNetworkRestrictedAnnotation = "forecastle.stakater.com/network-restricted"
	// ForecastleURLAnnotation const used for specifying the URL for the forecastle app
	ForecastleURLAnnotation = "forecastle.stakater.com/url"
	// ForecastlePropertiesAnnotation const used for specifying app properties
	ForecastlePropertiesAnnotation = "forecastle.stakater.com/properties"
)
