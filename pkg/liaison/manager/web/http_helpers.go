package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/liaisonio/liaison/pkg/liaison/manager/iam"
	"github.com/liaisonio/liaison/pkg/liaison/repo/model"
)

var errUnauthorized = errors.New("unauthorized")

// writeJSON writes a JSON body with the given HTTP status.
func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// bearerToken extracts the Authorization: Bearer <token> value, if present.
func bearerToken(r *http.Request) (string, bool) {
	v := r.Header.Get("Authorization")
	if v == "" {
		return "", false
	}
	parts := strings.SplitN(v, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

// authenticateHTTP resolves the user from the request. It prefers the context
// value set by AuthMiddleware (if it ran) and falls back to re-parsing the
// bearer token. Custom HandleFunc routes cannot rely on middleware context
// propagation, so the fallback path is not optional.
func (web *web) authenticateHTTP(r *http.Request) (*model.User, error) {
	if v := r.Context().Value("user"); v != nil {
		if u, ok := v.(*model.User); ok && u != nil {
			return u, nil
		}
	}
	token, ok := bearerToken(r)
	if !ok {
		return nil, errUnauthorized
	}
	if strings.HasPrefix(token, model.PATPlaintextPrefix) {
		user, err := web.iamService.GetUserByPAT(token, iam.ExtractClientIP(r))
		if err != nil {
			return nil, errUnauthorized
		}
		return user, nil
	}
	user, err := web.iamService.GetUserByToken(token)
	if err != nil {
		return nil, errUnauthorized
	}
	return user, nil
}

// writeUnauthorized writes a 401 response using the shared envelope shape.
func writeUnauthorized(w http.ResponseWriter) {
	writeJSON(w, http.StatusUnauthorized, map[string]any{
		"code":    http.StatusUnauthorized,
		"message": "unauthorized",
	})
}
