package client_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
			closeServer: true,
		},
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
			closeServer: true,
		},
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

			var out bytes.Buffer

			conf.ApiURL = url
			err := client.DeployFlow(&out, conf, tc.flowId, "active")
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
			if out.String() == "" {
				t.Errorf("response body should not be initial")
			}
			if strings.Contains(out.String(), `327626af-8e45-4c56-4791-4a4858573396`) == false {
				t.Errorf("response body should contain task id")
			}
			// if resp.D.ID != tc.flowId {
			// 	t.Errorf("Expected flowId: %s, got: %s", tc.flowId, resp.D.ID)
			// }
		})
	}
}

func TestGetFlowConfigs(t *testing.T) {

	testResp := map[string]struct {
		Status int
		Body   string
	}{
		"resultOK": {
			Status: http.StatusOK,
			Body: `{
        "d": {
                "results": [
                        {
                                "ParameterKey": "APIKey",
                                "ParameterValue": "2yhr5KY7wl2cKhqGHZjAMBQAKr6GFp1X",
                                "DataType": "xsd:string"
                        },
                        {
                                "ParameterKey": "ExProp1Value",
                                "ParameterValue": "testValue7",
                                "DataType": "xsd:string"
                        },
                        {
                                "ParameterKey": "bodySize",
                                "ParameterValue": "105",
                                "DataType": "xsd:integer"
                        }
                ]
        }
}`,
		},
		"notFound": {
			Status: http.StatusNotFound,
			Body: `{
        "d": {
                "results": null
        }
}`,
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
			name:     "resultOK",
			flowId:   "PurchaseOrder",
			expError: nil,
			resp:     testResp["resultOK"],
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
			closeServer: true,
		},
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
			resp, err := client.GetFlowConfigs(conf, tc.flowId, version)
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
			if len(resp.D.Results) != 3 {
				t.Errorf("Expected 3 parameters got %d", len(resp.D.Results))
			}
			if resp.D.Results[0].ParameterKey != "APIKey" {
				t.Errorf("Expected ParameterKey: %s, got: %s", "APIKey", resp.D.Results[0].ParameterKey)
			}
		})
	}
}

func TestCreateFlow(t *testing.T) {

	testResp := map[string]struct {
		Status int
		Body   string
	}{
		"resultOne": {
			Status: http.StatusOK,
			Body: `{
	"d": {
					"__metadata": {
									"id": "https://ec69d614trial.it-cpitrial03.cfapps.ap21.hana.ondemand.com:443/api/v1/IntegrationDesigntimeArtifacts(Id='PurchaseOrderTest',Version='1.0.0')",
									"uri": "https://ec69d614trial.it-cpitrial03.cfapps.ap21.hana.ondemand.com:443/api/v1/IntegrationDesigntimeArtifacts(Id='PurchaseOrderTest',Version='1.0.0')",
									"type": "com.sap.hci.api.IntegrationDesigntimeArtifact",
									"content_type": "application/octet-stream",
									"media_src": "https://ec69d614trial.it-cpitrial03.cfapps.ap21.hana.ondemand.com:443/api/v1/IntegrationDesigntimeArtifacts(Id='PurchaseOrderTest',Version='1.0.0')/$value",
									"edit_media": "https://ec69d614trial.it-cpitrial03.cfapps.ap21.hana.ondemand.com:443/api/v1/IntegrationDesigntimeArtifacts(Id='PurchaseOrderTest',Version='1.0.0')/$value"
					},
					"Id": "PurchaseOrderTest",
					"Version": "1.0.0",
					"PackageId": "POscenerio",
					"Name": "Purchase Order Test Create",
					"Description": "",
					"Sender": "",
					"Receiver": ""
	}
}`,
		},
		"notFound": {
			Status: http.StatusNotFound,
			Body:   `{"error":{"code":"Not Found","message":{"lang":"en","value":"Package ID POsceneri1 does not exist."}}}`,
		},
	}

	conf := getTestConfiguration()

	testCases := []struct {
		name      string
		flowId    string
		packageId string
		expError  error
		resp      struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:      "resultOne",
			flowId:    "PurchaseOrderTest",
			packageId: "POscenerio",
			expError:  nil,
			resp:      testResp["resultOne"],
		},
		{
			name:        "notFound",
			flowId:      "PurchaseOrderTest",
			packageId:   "notexistingPackageId",
			expError:    client.ErrNotFound,
			resp:        testResp["notFound"],
			closeServer: false,
		},
		{
			name:        "InvalidURL",
			flowId:      "PurchaseOrder",
			expError:    client.ErrConnection,
			resp:        testResp["notFound"],
			closeServer: true,
		},
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
			//version := "active"
			//var fileContent io.Reader
			fileContent := strings.NewReader("")
			resp, err := client.CreateFlow(conf, tc.flowId, tc.flowId, "packageId", "", fileContent)
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
