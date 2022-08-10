package handler

import (
	"encoding/json"
	"fmt"
	"google.golang.org/genproto/googleapis/cloud/compute/v1"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func filterResult(ruleType, version string, resource []*compute.WafExpressionSet) (filteredResponse []*compute.WafExpressionSet) {
	if ruleType == "" {
		filteredResponse = resource
	} else {
		for _, expression := range resource {
			// v33 is the latest version of preconfigured rules
			if strings.Contains(expression.GetId(), fmt.Sprintf("%s-%s", ruleType, version)) {
				filteredResponse = append(filteredResponse, expression)
			}
		}
	}
	return filteredResponse
}

func response(w http.ResponseWriter, response interface{}) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("encode %v", err), http.StatusInternalServerError)
		return
	}
}

func parse(input ...string) (bool, string) {
	// This will only match sequences of one or more sequences
	// of alphanumeric characters separated by a single -
	regex := "^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$"
	for _, v := range input {
		if v == "" {
			continue
		}
		if !regexp.MustCompile(regex).MatchString(v) {
			return false, v
		}
	}
	return true, ""
}

func parseInt(i string) (int32, error) {
	p, err := strconv.ParseInt(i, 10, 32)
	if err != nil {
		return int32(0), err
	}
	return int32(p), nil
}
