package client_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tobiaszgithub/cig/client"
	"github.com/tobiaszgithub/cig/config"
)

func TestInspectFlow(t *testing.T) {

	testResp := map[string]struct {
		Status int
		Body   string
	}{
		"resultOne": {
			Status: http.StatusOK,
			Body: `{
			"d": {
							"Id": "PurchaseOrder",
							"Version": "1.0.5",
							"PackageId": "POscenerio",
							"Name": "PurchaseOrder",
							"Description": "PO notifications",
							"Sender": "",
							"Receiver": ""
			}
		}`,
		},
		"notFound": {
			Status: http.StatusNotFound,
			Body:   `{"error":{"code":"Not Found","message":{"lang":"en","value":"Integration design time artifact not found"}}}`,
		},
	}

	conf := getTestConfiguration()

	testCases := []struct {
		name     string
		flowId   string
		expError error
		resp     struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:     "resultOne",
			flowId:   "PurchaseOrder",
			expError: nil,
			resp:     testResp["resultOne"],
		},
		{
			name:        "notFound",
			flowId:      "notExistingFlowId",
			expError:    client.ErrNotFound,
			resp:        testResp["notFound"],
			closeServer: false,
		},
		{
			name:        "InvalidURL",
			flowId:      "PurchaseOrder",
			expError:    client.ErrConnection,
			resp:        testResp["notFound"],
			closeServer: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.resp.Status)
					fmt.Fprintln(w, tc.resp.Body)
				})
			defer cleanup()
			if tc.closeServer {
				cleanup()
			}

			conf.ApiURL = url
			version := "active"
			resp, err := client.InspectFlow(conf, tc.flowId, version)
			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got no error.", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("Expected error %q, got %q.", tc.expError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got %q.", err)
			}
			if resp.D.ID == "" {
				t.Errorf("flow ID should not be initial")
			}
			if resp.D.ID != tc.flowId {
				t.Errorf("Expected flowId: %s, got: %s", tc.flowId, resp.D.ID)
			}
		})
	}
}

func getTestConfiguration() config.Configuration {
	conf := config.Configuration{}
	conf.ApiURL = ""
	conf.Key = "test"
	conf.Authorization.Type = "basic"
	conf.Authorization.Username = ""
	conf.Authorization.Password = ""

	return conf
}

func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)

	return ts.URL, func() {
		ts.Close()
	}
}

func TestDeployFlow(t *testing.T) {
	testResp := map[string]struct {
		Status int
		Body   string
	}{
		"resultOne": {
			Status: http.StatusOK,
			Body:   `327626af-8e45-4c56-4791-4a4858573396`,
		},
		"notFound": {
			Status: http.StatusNotFound,
			Body:   `{"error":{"code":"Not Found","message":{"lang":"en","value":"Integration design time artifact not found"}}}`,
		},
	}

	conf := getTestConfiguration()

	testCases := []struct {
		name     string
		flowId   string
		expError error
		resp     struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:     "resultOne",
			flowId:   "PurchaseOrder",
			expError: nil,
			resp:     testResp["resultOne"],
		},
		{
			name:        "notFound",
			flowId:      "notExistingFlowId",
			expError:    client.ErrNotFound,
			resp:        testResp["notFound"],
			closeServer: false,
		},
		{
			name:        "InvalidURL",
			flowId:      "PurchaseOrder",
			expError:    client.ErrConnection,
			resp:        testResp["notFound"],
			closeServer: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.resp.Status)
					fmt.Fprintln(w, tc.resp.Body)
				})
			defer cleanup()
			if tc.closeServer {
				cleanup()
			}

			conf.ApiURL = url
			resp, err := client.DeployFlow(conf, tc.flowId, "active")
			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got no error.", tc.expError)
				}
				if !errors.Is(err, tc.expError) {
					t.Errorf("Expected error %q, got %q.", tc.expError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got %q.", err)
			}
			if resp == "" {
				t.Errorf("flow ID should not be initial")
			}
			// if resp.D.ID != tc.flowId {
			// 	t.Errorf("Expected flowId: %s, got: %s", tc.flowId, resp.D.ID)
			// }
		})
	}
}
