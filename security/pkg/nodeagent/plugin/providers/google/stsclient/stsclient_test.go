// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stsclient

import (
	"context"
	"testing"

	"istio.io/istio/security/pkg/stsservice/tokenmanager/google/mock"
)

func TestGetFederatedToken(t *testing.T) {
	GKEClusterURL = mock.FakeGKEClusterURL
	r := NewPlugin()

	ms, err := mock.StartNewServer(t, mock.Config{Port: 0})
	if err != nil {
		t.Fatalf("failed to start a mock server: %v", err)
	}
	SecureTokenEndpoint = ms.URL + "/v1/identitybindingtoken"
	defer func() {
		if err := ms.Stop(); err != nil {
			t.Logf("failed to stop mock server: %v", err)
		}
		SecureTokenEndpoint = "https://securetoken.googleapis.com/v1/identitybindingtoken"
	}()

	token, _, _, err := r.ExchangeToken(context.Background(), mock.FakeTrustDomain, mock.FakeSubjectToken)
	if err != nil {
		t.Fatalf("failed to call exchange token %v", err)
	}
	if token != mock.FakeFederatedToken {
		t.Errorf("Access token got %q, expected %q", token, mock.FakeFederatedToken)
	}
}

func TestRealToken(t *testing.T) {
	GKEClusterURL = "https://container.googleapis.com/v1/projects/sven-asm-vms/locations/us-west3-a/clusters/vm-cluster"
	SecureTokenEndpoint = "https://securetoken.googleapis.com/v1/identitybindingtoken"
	r := NewPlugin()

	trustDomain := "sven-asm-vms.svc.id.goog"
	jwt := "eyJhbGciOiJSUzI1NiIsImtpZCI6Im5EcWZDbFYwSHN6dU8tTHUwNnFBYzRtLTlrTGNNVDRKNGNPcUFzd1A1a3MifQ.eyJhdWQiOlsic3Zlbi1hc20tdm1zLnN2Yy5pZC5nb29nIl0sImV4cCI6MTU5MDY4NjQxMSwiaWF0IjoxNTkwNjgyODExLCJpc3MiOiJodHRwczovL2NvbnRhaW5lci5nb29nbGVhcGlzLmNvbS92MS9wcm9qZWN0cy9zdmVuLWFzbS12bXMvbG9jYXRpb25zL3VzLXdlc3QzLWEvY2x1c3RlcnMvdm0tY2x1c3RlciIsImt1YmVybmV0ZXMuaW8iOnsibmFtZXNwYWNlIjoiaGlwc3RlciIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJkZWZhdWx0IiwidWlkIjoiNGMxOTlhMWYtNTIwNC00ZmEzLTk0ZjAtZTIwMDg1ZmJmZWUzIn19LCJuYmYiOjE1OTA2ODI4MTEsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpoaXBzdGVyOmRlZmF1bHQifQ.dW2Jb70r0uG-mMZjZ5Vje74mHPERIzYjkv9ZX4Q2LHbcxnA0CidM13oJzcw01msdPpvO_t5A_nmkXAqZccGSPcm9wROGYLvuP3FthA2ktlKoCWN6LNP_Vjoxw3qjh1i6TkOuT8GvxmVz_bABTFe9WOT5qN0n6VHq8mjYSedyAmkMwRhPJmarte8G1znQzywZ-uWZNBsa3YfLHQvpYQdFp0oYgPPGLKo6LkdrRJOV0KM2jx9hmL39KcnDZTeSI50e9GYWSwwnO_3eDRUeBIbPNZSEwLu9XD4Od0UoTejFONzSctae8tY3QIPYz1tBVDG0zIgGx1Ugece6VH7pcVRmzQ"

	token, _, code, err := r.ExchangeToken(context.Background(), trustDomain, jwt)
	if err != nil {
		t.Fatalf("failed on exchange %v", err)
	}
	if len(token) == 0 {
		t.Fatalf("bad token %v", token)
	}
	if code != 200 {
		t.Fatalf("got response code %v", code)
	}
}
