package fake

import (
	"fmt"
	"google.golang.org/api/option"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

func newSecurityPoliciesRESTClient(policy string, exists bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		// GET security policies and rules
		case r.Method == http.MethodGet:
			if strings.HasSuffix(r.URL.Path, "securityPolicies") {
				response, _ := os.ReadFile("testdata/list-policies-reponse.json")
				_, _ = w.Write(response)
			}
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("securityPolicies/%s", policy)) {
				response, _ := os.ReadFile("testdata/get-policy-reponse.json")
				_, _ = w.Write(response)
			}
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("getRule")) {
				response, _ := os.ReadFile("testdata/get-rule-reponse.json")
				_, _ = w.Write(response)
			}
			if strings.HasSuffix(r.URL.Path, "operations/") {
				response, _ := os.ReadFile("testdata/operation-done-reponse.json")
				_, _ = w.Write(response)
			}
		// PATCH security policies
		case r.Method == http.MethodPatch:
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("securityPolicies/%s", policy)) {
				response, _ := os.ReadFile("testdata/operation-patch-reponse.json")
				_, _ = w.Write(response)
			}
		// POST (Patch) security rule
		case r.Method == http.MethodPost:
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("patchRule")) {
				response, _ := os.ReadFile("testdata/operation-patch-reponse.json")
				_, _ = w.Write(response)
			}
		}
	}
}

func newBackendServicesRESTClient(policy string, exists bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		// POST security policy to backend service
		case r.Method == http.MethodPost:
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("patchRule")) {
				response, _ := os.ReadFile("testdata/operation-patch-reponse.json")
				_, _ = w.Write(response)
			}
		}
	}
}

func SecurityApiServer(existingPolicy string, exists bool) ([]option.ClientOption, error) {
	testServer := httptest.NewServer(newSecurityPoliciesRESTClient(existingPolicy, exists))
	opts := []option.ClientOption{
		option.WithEndpoint(testServer.URL),
		option.WithoutAuthentication(),
	}
	return opts, nil
}

func ServiceApiServer(existingPolicy string, exists bool) ([]option.ClientOption, error) {
	testServer := httptest.NewServer(newBackendServicesRESTClient(existingPolicy, exists))
	opts := []option.ClientOption{
		option.WithEndpoint(testServer.URL),
		option.WithoutAuthentication(),
	}
	return opts, nil
}
