package client_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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
			fileContent := strings.NewReader("test string")
			resp, err := client.CreateFlow(conf, tc.flowId, tc.flowId, "packageId", fileContent)
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

func TestUpdateFlow(t *testing.T) {

	testResp := map[string]struct {
		Status int
		Body   string
	}{
		"resultOne": {
			Status: http.StatusOK,
			Body:   `Integration flow: PurchaseOrderTest updated`,
		},
		"notFound": {
			Status: http.StatusNotFound,
			Body:   `{"error":{"code":"Not Found","message":{"lang":"en","value":"Integration design time artifact not found."}}}`,
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
			var out bytes.Buffer
			err := client.UpdateFlow(&out, conf, tc.flowId, tc.flowId, "packageId", "", fileContent)
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

			if strings.Contains(out.String(), tc.flowId) == false {
				t.Errorf("response body should contain flow id")
			}

		})
	}
}

func TestDownloadFlow(t *testing.T) {
	testResp := map[string]struct {
		Status int
		Body   string
	}{
		"ok": {
			Status: http.StatusOK,
			Body:   `example-test-file-content`,
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
			name:     "ok",
			flowId:   "PurchaseOrder",
			expError: nil,
			resp:     testResp["ok"],
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

			var out bytes.Buffer
			var outputContent bytes.Buffer

			err := client.DownloadFlow(&out, conf, tc.flowId, "active", &outputContent)
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
			if strings.Contains(out.String(), `number of bytes`) == false {
				t.Errorf("response body should contain number of bytes")
			}
		})
	}
}

func TestCopyFlow(t *testing.T) {

	testResp := map[string]struct {
		Status int
		Body   string
	}{
		"inspectFlowStatusOK": {
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
		"createFlowStatusOK": {
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
					"Id": "PurchaseOrderCopy1",
					"Version": "1.0.0",
					"PackageId": "POscenerio",
					"Name": "PurchaseOrder Copy1",
					"Description": "",
					"Sender": "",
					"Receiver": ""
	}
}`,
		},

		"notFound": {
			Status: http.StatusNotFound,
			Body:   `{"error":{"code":"Not Found","message":{"lang":"en","value":"Integration design time artifact not found."}}}`,
		},
	}

	conf := getTestConfiguration()

	testCases := []struct {
		name          string
		srcFlowId     string
		destFlowId    string
		destPackageId string
		destFlowName  string
		expError      error
		resp          struct {
			Status int
			Body   string
		}
		inspectResp struct {
			Status int
			Body   string
		}
		downloadResp struct {
			Status int
			Body   string
		}
		createResp struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:          "flowCopyStatusOk",
			srcFlowId:     "PurchaseOrder",
			destFlowId:    "PurchaseOrderCopy1",
			destPackageId: "POscenerio",
			destFlowName:  "PurchaseOrder Copy1",
			expError:      nil,
			resp:          testResp["inspectFlowStatusOK"],
			inspectResp:   testResp[""],
			downloadResp:  testResp[""],
			createResp:    testResp["createFlowStatusOK"],
		},
		{
			name:          "notFound",
			srcFlowId:     "PurchaseOrderTest",
			destFlowId:    "PurchaseOrderCopy1",
			destPackageId: "notexistingPackageId",
			destFlowName:  "PurchaseOrder Copy1",
			expError:      client.ErrNotFound,
			resp:          testResp["notFound"],
			closeServer:   false,
		},
		{
			name:        "InvalidURL",
			srcFlowId:   "PurchaseOrder",
			destFlowId:  "PurchaseOrderCopy1",
			expError:    client.ErrConnection,
			resp:        testResp["notFound"],
			closeServer: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					if r.Method == "POST" {
						w.WriteHeader(tc.createResp.Status)
						fmt.Fprintln(w, tc.createResp.Body)
					} else {
						w.WriteHeader(tc.resp.Status)
						fmt.Fprintln(w, tc.resp.Body)
					}
				})
			defer cleanup()
			if tc.closeServer {
				cleanup()
			}

			conf.ApiURL = url

			//fileContent := strings.NewReader("")
			var out bytes.Buffer
			//err := client.UpdateFlow(&out, conf, tc.flowId, tc.flowId, "packageId", "", fileContent)
			err := client.CopyFlow(&out, conf, tc.srcFlowId, tc.destFlowId, tc.destFlowName, tc.destPackageId)
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
			outStr := out.String()
			if strings.Contains(outStr, tc.destFlowId) == false {
				t.Errorf("response body should contain created flow id")
			}

		})
	}
}

func TestTransportFlow(t *testing.T) {

	testResp := map[string]struct {
		Status int
		Body   string
	}{
		"inspectFlowStatusOK": {
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
		"downloadFlowStatusOK": {
			Status: http.StatusOK,
			Body:   `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98XNBAAAcAwAADYAAAAAAAAAAAAAAAAA3hQAAHNyYy9tYWluL3Jlc291cmNlcy9tYXBwaW5nL09EYXRhX3NvdXJjZV9QT19wcmVwcm9jLnhzbFBLAQIUABQACAgIAIN0/FReDtP6LzkAAAEMAgBIAAAAAAAAAAAAAAAAAA8aAABzcmMvbWFpbi9yZXNvdXJjZXMvbWFwcGluZy9NYXBwaW5nX09EYXRhX3NvdXJjZV9QT190b19TT0FQX3RhcmdldF9QTy54c2xQSwECFAAUAAgICACEdPxUkeRmrEAGAAB7GgAANgAAAAAAAAAAAAAAAAC0UwAAc3JjL21haW4vcmVzb3VyY2VzL21hcHBpbmcvU09BUF90YXJnZXRfUE9fcG9zdHByb2MueHNsUEsBAhQAFAAICAgAhHT8VO2wmzfIAQAAnwUAADcAAAAAAAAAAAAAAAAAWFoAAHNyYy9tYWluL3Jlc291cmNlcy9zY3JpcHQvc3ViX3NwZWNpYWxfY2hhcmFjdGVycy5ncm9vdnlQSwECFAAUAAgICACEdPxUBJUCR9QAAAC3AQAACAAAAAAAAAAAAAAAAACFXAAALnByb2plY3RQSwECFAAUAAgICACEdPxU8wdz8nkDAAAOCgAAFAAAAAAAAAAAAAAAAACPXQAATUVUQS1JTkYvTUFOSUZFU1QuTUZQSwUGAAAAAAkACQAvAwAASmEAAAAA`,
		},
		"createFlowStatusOK": {
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
					"Id": "PurchaseOrderCopy1",
					"Version": "1.0.0",
					"PackageId": "POscenerio",
					"Name": "PurchaseOrder Copy1",
					"Description": "",
					"Sender": "",
					"Receiver": ""
	}
}`,
		},

		"notFound": {
			Status: http.StatusNotFound,
			Body:   `{"error":{"code":"Not Found","message":{"lang":"en","value":"Integration design time artifact not found."}}}`,
		},
	}

	conf := getTestConfiguration()
	destConf := getTestConfiguration()

	testCases := []struct {
		name          string
		srcFlowId     string
		destFlowId    string
		destPackageId string
		destFlowName  string
		expError      error
		resp          struct {
			Status int
			Body   string
		}
		inspectResp struct {
			Status int
			Body   string
		}
		downloadResp struct {
			Status int
			Body   string
		}
		createResp struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:          "flowCopyStatusOk",
			srcFlowId:     "PurchaseOrder",
			destFlowId:    "PurchaseOrderCopy1",
			destPackageId: "POscenerio",
			destFlowName:  "PurchaseOrder Copy1",
			expError:      nil,
			resp:          testResp["inspectFlowStatusOK"],
			inspectResp:   testResp[""],
			downloadResp:  testResp["downloadFlowStatusOK"],
			createResp:    testResp["createFlowStatusOK"],
		},
		{
			name:          "notFound",
			srcFlowId:     "PurchaseOrderTest",
			destFlowId:    "PurchaseOrderCopy1",
			destPackageId: "notexistingPackageId",
			destFlowName:  "PurchaseOrder Copy1",
			expError:      client.ErrNotFound,
			resp:          testResp["notFound"],
			closeServer:   false,
		},
		{
			name:        "InvalidURL",
			srcFlowId:   "PurchaseOrder",
			destFlowId:  "PurchaseOrderCopy1",
			expError:    client.ErrConnection,
			resp:        testResp["notFound"],
			closeServer: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					if r.Method == "POST" {
						w.WriteHeader(tc.createResp.Status)
						fmt.Fprintln(w, tc.createResp.Body)
						return
					}

					urlPath := r.URL.Path

					if r.Method == "GET" && strings.Contains(urlPath, "$value") {
						w.WriteHeader(tc.downloadResp.Status)
						//dec, _ := base64.StdEncoding.DecodeString(tc.downloadResp.Body)
						decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(tc.downloadResp.Body))
						//fmt.Fprintln(w, decoder)
						io.Copy(w, decoder)
					} else {
						w.WriteHeader(tc.resp.Status)
						fmt.Fprintln(w, tc.resp.Body)
					}
				})
			defer cleanup()
			if tc.closeServer {
				cleanup()
			}

			conf.ApiURL = url
			destConf.ApiURL = url

			//fileContent := strings.NewReader("")
			var out bytes.Buffer
			err := client.TransportFlow(&out, conf, tc.srcFlowId, destConf, tc.destFlowId, tc.destFlowName, tc.destPackageId)
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
			outStr := out.String()
			if strings.Contains(outStr, tc.destFlowId) == false {
				t.Errorf("response body should contain transported flow id")
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
