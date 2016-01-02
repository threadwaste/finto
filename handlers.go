package finto

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Contains application context.
type fintoContext struct {
	set          *RoleSet
	instanceRole string
}

func InitFintoContext(rs *RoleSet) *fintoContext {
	return &fintoContext{
		set: rs,
	}
}

func (fc *fintoContext) setInstanceRole(role string) error {
	_, err := fc.set.Role(role)
	if err != nil {
		return err
	}

	fc.instanceRole = role
	return nil
}

// VarsHandlerFunc accepts mux route variables as an argument.
type VarsHandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

func (f VarsHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	f(w, r, vars)
}

// List available roles.
func rolesList(fc *fintoContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var roles []string

		if r.FormValue("status") == "active" {
			roles = []string{fc.instanceRole}
		} else {
			roles = fc.set.Roles()
		}

		jsonResponse(w, map[string][]string{"roles": roles})
	})
}

// Show a role's configuration.
func rolesShow(fc *fintoContext) http.Handler {
	return VarsHandlerFunc(func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
		role, err := fc.set.Role(vars["alias"])
		if err != nil {
			errorResponse(w, err.Error(), http.StatusNotFound)
			return
		}

		jsonResponse(w, map[string]string{
			"arn":          role.Arn(),
			"session_name": role.SessionName(),
		})
	})
}

// Set role to be served as the instance profile role.
func rolesSetActive(fc *fintoContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type activateRequest struct {
			Alias string `json:"alias"`
		}

		var req activateRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := fc.setInstanceRole(req.Alias); err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonResponse(w, map[string]string{"active_role": fc.instanceRole})
	})
}

// Mock the EC2 security-credentials meta-data endpoint.
func mockProfile(fc *fintoContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fc.instanceRole))
	})
}

// Mock the EC2 instance profile role meta-data endpoint.
func mockProfileCreds(fc *fintoContext) http.Handler {
	return VarsHandlerFunc(func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
		role, err := fc.set.Role(vars["alias"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		creds, err := role.Credentials()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(map[string]string{
			"Code":            "Success",
			"LastUpdated":     "2015-07-07T23:06:33Z",
			"Type":            "AWS-HMAC",
			"AccessKeyId":     creds.AccessKeyId,
			"SecretAccessKey": creds.SecretAccessKey,
			"Token":           creds.SessionToken,
			"Expiration":      creds.Expiration.Format("2006-01-02T15:04:05Z"),
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var o bytes.Buffer
		json.Indent(&o, b, "", "  ")

		o.WriteTo(w)
	})
}

func jsonResponse(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(body)
}

func errorResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	jsonResponse(w, map[string]string{"error": message})
}
