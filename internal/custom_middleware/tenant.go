package custom_middleware

import (
	"net/http"
	"strings"
	"test_go_api/internal/database"
)

// TenantDB defines the database operations needed by the tenant middleware
type TenantDB interface {
	// Add any database methods you need here
	// For example:
	// ValidateTenant(tenantID string) error
}

func NewTenantMiddleware(db database.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host
			if strings.Contains(host, "127.0.0.1") {
				http.Error(w, "Tenant ID is required in subdomain", http.StatusBadRequest)
				return
			}
			parts := strings.Split(host, ".")

			// If we have a subdomain (e.g., tenant.example.com)
			if len(parts) > 2 {
				tenantId := parts[0]
				if tenantId == "" {
					http.Error(w, "Tenant ID is required", http.StatusBadRequest)
					return
				}

				db.SetTenant(tenantId)
			} else {
				http.Error(w, "Tenant ID is required in subdomain", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
