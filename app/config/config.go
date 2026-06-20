package config

import "os"

// ShopEnabled returns true when SHOP_ENABLED=true in the environment.
func ShopEnabled() bool {
	return os.Getenv("SHOP_ENABLED") == "true"
}

// ResultadosEnabled returns true when RESULTADOS_ENABLED=true in the environment.
func ResultadosEnabled() bool {
	return os.Getenv("RESULTADOS_ENABLED") == "true"
}

// MapsAPIKey returns the Google Maps JavaScript API key from MAPS_API_KEY.
func MapsAPIKey() string {
	return os.Getenv("MAPS_API_KEY")
}

// MapsAddress returns the gym address used to geocode and pin the map marker.
// Override with the MAPS_ADDRESS env var once the real address is known.
func MapsAddress() string {
	if a := os.Getenv("MAPS_ADDRESS"); a != "" {
		return a
	}
	return "Rua das Artes Marciais, 123, Centro, São Paulo, SP"
}
