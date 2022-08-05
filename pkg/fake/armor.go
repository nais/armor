package fake

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func NewSecurityPoliciesRESTClient(policy string, exists bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		println(r.Method)
		switch {
		// GET security policies and rules
		case r.Method == http.MethodGet:
			if strings.HasSuffix(r.URL.Path, "securityPolicies") {
				response, _ := ioutil.ReadFile(fmt.Sprintf("common/testdata/list-policies-reponse.json"))
				_, _ = w.Write(response)
			}
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("securityPolicies/%s", policy)) {
				response, _ := ioutil.ReadFile(fmt.Sprintf("common/testdata/get-policy-reponse.json"))
				_, _ = w.Write(response)
			}
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("getRule")) {
				response, _ := ioutil.ReadFile(fmt.Sprintf("common/testdata/get-rule-reponse.json"))
				_, _ = w.Write(response)
			}
			if strings.HasSuffix(r.URL.Path, "operations/") {
				response, _ := ioutil.ReadFile(fmt.Sprintf("common/testdata/operation-done-reponse.json"))
				_, _ = w.Write(response)
			}
		// PATCH security policies
		case r.Method == http.MethodPatch:
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("securityPolicies/%s", policy)) {
				response, _ := ioutil.ReadFile(fmt.Sprintf("common/testdata/operation-patch-reponse.json"))
				_, _ = w.Write(response)
			}
		// POST (Patch) security rule
		case r.Method == http.MethodPost:
			if strings.HasSuffix(r.URL.Path, fmt.Sprintf("patchRule")) {
				response, _ := ioutil.ReadFile(fmt.Sprintf("common/testdata/operation-patch-reponse.json"))
				_, _ = w.Write(response)
			}
		}
	}
}
