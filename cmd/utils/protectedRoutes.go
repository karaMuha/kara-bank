package utils

var protectedRoutes map[string]bool

func SetProtectedRoutes() {
	protectedRoutes = make(map[string]bool, 4)
	protectedRoutes["POST users"] = false
	protectedRoutes["POST accounts"] = true
	protectedRoutes["GET accounts"] = true
	protectedRoutes["POST transfer"] = true
}

func IsProtectedRoute(endpoint string) bool {
	return protectedRoutes[endpoint]
}
