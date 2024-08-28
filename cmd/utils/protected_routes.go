package utils

import "github.com/vodkaslime/wildcard"

var protectedRoutes map[string][]string
var matcher *wildcard.Matcher

func SetProtectedRoutes() {
	matcher = wildcard.NewMatcher()
	protectedRoutes = make(map[string][]string)
	protectedRoutes["POST /users/register"] = []string{""}
	protectedRoutes["POST /users/login"] = []string{""}
	protectedRoutes["POST /accounts"] = []string{"customer", "banker", "admin"}
	protectedRoutes["GET /accounts/*"] = []string{"customer", "banker", "admin"}
	protectedRoutes["GET /accounts"] = []string{"banker", "admin"}
	protectedRoutes["POST /transfers"] = []string{"customer", "admin"}
}

func IsProtectedRoute(endpoint string) ([]string, error) {
	for i, v := range protectedRoutes {
		if result, err := matcher.Match(i, endpoint); err == nil && result {
			return v, nil
		} else if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
