package finto

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestRequest(m, p string, b io.Reader, t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(m, p, b)
	if err != nil {
		t.Errorf("create request failed: %s", err.Error())
	}

	rec := httptest.NewRecorder()

	return req, rec
}

func setupTestFintoContext() (fc *fintoContext) {
	ts := NewRoleSet(&MockAssumeRoleClient{})
	ts.SetRole("test-alias", "test-arn")
	ts.SetRole("another-alias", "another-arn")

	fc = &fintoContext{
		set:          ts,
		instanceRole: "test-alias",
	}
	fc.setInstanceRole("test-alias")

	return
}

type handlerTest struct {
	method, path string
	body         io.Reader
	responseCode int
	responseBody map[string]interface{}
}

func TestFintoHandlers(t *testing.T) {
	// TODO: This is disgusting
	me := MockExpiry.Add(-300)

	cases := []handlerTest{
		{
			"GET",
			"/roles",
			nil,
			http.StatusOK,
			map[string]interface{}{
				"roles": []interface{}{"another-alias", "test-alias"},
			},
		},
		{
			"GET",
			"/roles?status=active",
			nil,
			http.StatusOK,
			map[string]interface{}{
				"roles": []interface{}{
					"test-alias",
				},
			},
		},
		{
			"PUT",
			"/roles",
			bytes.NewBuffer([]byte(`{"alias":"another-alias"}`)),
			http.StatusOK,
			map[string]interface{}{
				"active_role": "another-alias",
			},
		},
		{
			"PUT",
			"/roles",
			bytes.NewBuffer([]byte(`{"alias":"missing-alias"}`)),
			http.StatusBadRequest,
			map[string]interface{}{
				"error": "unknown role: missing-alias",
			},
		},
		{
			"GET",
			"/roles/test-alias",
			nil,
			http.StatusOK,
			map[string]interface{}{
				"arn":          "test-arn",
				"session_name": "finto-test-alias",
			},
		},
		{
			"GET",
			"/roles/missing-alias",
			nil,
			http.StatusNotFound,
			map[string]interface{}{
				"error": "unknown role: missing-alias",
			},
		},
		{
			"GET",
			"/roles/test-alias/credentials",
			nil,
			http.StatusOK,
			map[string]interface{}{
				"Code":            "Success",
				"LastUpdated":     "2015-07-07T23:06:33Z",
				"Type":            "AWS-HMAC",
				"AccessKeyId":     "test-arn-finto-test-alias",
				"SecretAccessKey": "mock-key",
				"Token":           "mock-token",
				"Expiration":      me.Format("2006-01-02T15:04:05Z"),
			},
		},
		{
			"GET",
			"/latest/meta-data/iam/security-credentials/test-alias",
			nil,
			http.StatusOK,
			map[string]interface{}{
				"Code":            "Success",
				"LastUpdated":     "2015-07-07T23:06:33Z",
				"Type":            "AWS-HMAC",
				"AccessKeyId":     "test-arn-finto-test-alias",
				"SecretAccessKey": "mock-key",
				"Token":           "mock-token",
				"Expiration":      me.Format("2006-01-02T15:04:05Z"),
			},
		},
	}

	for _, test := range cases {
		fc := setupTestFintoContext()
		router := FintoRouter(fc)

		req, rec := setupTestRequest(test.method, test.path, test.body, t)
		router.ServeHTTP(rec, req)

		var resp interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)

		if assert.NoError(t, err, rec.Body.String()) {
			assert.Equal(t, test.responseCode, rec.Code, test.path)
			assert.Equal(t, test.responseBody, resp, test.path)
		}
	}
}
func TestMockInstanceRole(t *testing.T) {
	req, rec := setupTestRequest(
		"GET",
		"/latest/meta-data/iam/security-credentials/",
		nil,
		t,
	)
	fc := setupTestFintoContext()

	FintoRouter(fc).ServeHTTP(rec, req)

	assert.Equal(t, "test-alias", rec.Body.String())
}
