package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// upsertFirewallRequest is the body of PUT /api/v1/proxies/{id}/firewall.
type upsertFirewallRequest struct {
	AllowedCIDRs []string `json:"allowed_cidrs"`
}

// handleFirewallHTTP dispatches GET/PUT/DELETE on /api/v1/proxies/{id}/firewall.
// Registered via HandleFunc, so it authenticates itself.
func (web *web) handleFirewallHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		web.handleGetProxyFirewallHTTP(w, r)
	case http.MethodPut:
		web.handleUpsertProxyFirewallHTTP(w, r)
	case http.MethodDelete:
		web.handleDeleteProxyFirewallHTTP(w, r)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"code":    http.StatusMethodNotAllowed,
			"message": "method not allowed",
		})
	}
}

func (web *web) handleGetProxyFirewallHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := web.authenticateHTTP(r)
	if err != nil {
		writeUnauthorized(w)
		return
	}
	proxyID, err := parseFirewallProxyID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "message": "invalid proxy id"})
		return
	}
	ctx := context.WithValue(r.Context(), "user_id", user.ID)
	data, err := web.controlPlane.GetProxyFirewall(ctx, proxyID)
	if err != nil {
		status := firewallHTTPStatus(err)
		writeJSON(w, status, map[string]any{"code": status, "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"code": 200, "message": "success", "data": data})
}

func (web *web) handleUpsertProxyFirewallHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := web.authenticateHTTP(r)
	if err != nil {
		writeUnauthorized(w)
		return
	}
	proxyID, err := parseFirewallProxyID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "message": "invalid proxy id"})
		return
	}
	var req upsertFirewallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "message": "invalid request body"})
		return
	}
	if req.AllowedCIDRs == nil {
		req.AllowedCIDRs = []string{}
	}
	ctx := context.WithValue(r.Context(), "user_id", user.ID)
	data, err := web.controlPlane.UpsertProxyFirewall(ctx, proxyID, req.AllowedCIDRs)
	if err != nil {
		status := firewallHTTPStatus(err)
		writeJSON(w, status, map[string]any{"code": status, "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"code": 200, "message": "success", "data": data})
}

func (web *web) handleDeleteProxyFirewallHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := web.authenticateHTTP(r)
	if err != nil {
		writeUnauthorized(w)
		return
	}
	proxyID, err := parseFirewallProxyID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "message": "invalid proxy id"})
		return
	}
	ctx := context.WithValue(r.Context(), "user_id", user.ID)
	if err := web.controlPlane.DeleteProxyFirewall(ctx, proxyID); err != nil {
		status := firewallHTTPStatus(err)
		writeJSON(w, status, map[string]any{"code": status, "message": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"code": 200, "message": "success"})
}

// parseFirewallProxyID extracts {id} from /api/v1/proxies/{id}/firewall.
func parseFirewallProxyID(r *http.Request) (uint, error) {
	// Strip the fixed prefix/suffix and parse what's left.
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/proxies/")
	path = strings.TrimSuffix(path, "/firewall")
	id, err := strconv.ParseUint(path, 10, 32)
	if err != nil || id == 0 {
		return 0, errors.New("invalid proxy id")
	}
	return uint(id), nil
}

func firewallHTTPStatus(err error) int {
	// Only ErrForbidden is mapped to 403 for now; everything else is a bad
	// request (invalid CIDR, missing proxy, etc.).
	if errors.Is(err, errUnauthorized) {
		return http.StatusUnauthorized
	}
	return http.StatusBadRequest
}
