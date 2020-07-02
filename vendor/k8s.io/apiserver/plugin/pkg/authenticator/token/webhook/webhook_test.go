/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"k8s.io/api/authentication/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/client-go/tools/clientcmd/api/v1"
)

// Service mocks a remote authentication service.
type Service interface {
	// Review looks at the TokenReviewSpec and provides an authentication
	// response in the TokenReviewStatus.
	Review(*v1beta1.TokenReview)
	HTTPStatusCode() int
}

// NewTestServer wraps a Service as an httptest.Server.
func NewTestServer(s Service, cert, key, caCert []byte) (*httptest.Server, error) {
	const webhookPath = "/testserver"
	var tlsConfig *tls.Config
	if cert != nil {
		cert, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	if caCert != nil {
		rootCAs := x509.NewCertPool()
		rootCAs.AppendCertsFromPEM(caCert)
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}
		tlsConfig.ClientCAs = rootCAs
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	serveHTTP := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, fmt.Sprintf("unexpected method: %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != webhookPath {
			http.Error(w, fmt.Sprintf("unexpected path: %v", r.URL.Path), http.StatusNotFound)
			return
		}

		var review v1beta1.TokenReview
		bodyData, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(bodyData, &review); err != nil {
			http.Error(w, fmt.Sprintf("failed to decode body: %v", err), http.StatusBadRequest)
			return
		}
		// ensure we received the serialized tokenreview as expected
		if review.APIVersion != "authentication.k8s.io/v1beta1" {
			http.Error(w, fmt.Sprintf("wrong api version: %s", string(bodyData)), http.StatusBadRequest)
			return
		}
		// once we have a successful request, always call the review to record that we were called
		s.Review(&review)
		if s.HTTPStatusCode() < 200 || s.HTTPStatusCode() >= 300 {
			http.Error(w, "HTTP Error", s.HTTPStatusCode())
			return
		}
		type userInfo struct {
			Username string              `json:"username"`
			UID      string              `json:"uid"`
			Groups   []string            `json:"groups"`
			Extra    map[string][]string `json:"extra"`
		}
		type status struct {
			Authenticated bool     `json:"authenticated"`
			User          userInfo `json:"user"`
		}

		var extra map[string][]string
		if review.Status.User.Extra != nil {
			extra = map[string][]string{}
			for k, v := range review.Status.User.Extra {
				extra[k] = v
			}
		}

		resp := struct {
			Kind       string `json:"kind"`
			APIVersion string `json:"apiVersion"`
			Status     status `json:"status"`
		}{
			Kind:       "TokenReview",
			APIVersion: v1beta1.SchemeGroupVersion.String(),
			Status: status{
				review.Status.Authenticated,
				userInfo{
					Username: review.Status.User.Username,
					UID:      review.Status.User.UID,
					Groups:   review.Status.User.Groups,
					Extra:    extra,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(serveHTTP))
	server.TLS = tlsConfig
	server.StartTLS()

	// Adjust the path to point to our custom path
	serverURL, _ := url.Parse(server.URL)
	serverURL.Path = webhookPath
	server.URL = serverURL.String()

	return server, nil
}

// A service that can be set to say yes or no to authentication requests.
type mockService struct {
	allow      bool
	statusCode int
	called     int
}

func (m *mockService) Review(r *v1beta1.TokenReview) {
	m.called++
	r.Status.Authenticated = m.allow
	if m.allow {
		r.Status.User.Username = "realHooman@email.com"
	}
}
func (m *mockService) Allow()              { m.allow = true }
func (m *mockService) Deny()               { m.allow = false }
func (m *mockService) HTTPStatusCode() int { return m.statusCode }

// newTokenAuthenticator creates a temporary kubeconfig file from the provided
// arguments and attempts to load a new WebhookTokenAuthenticator from it.
func newTokenAuthenticator(serverURL string, clientCert, clientKey, ca []byte, cacheTime time.Duration) (*WebhookTokenAuthenticator, error) {
	tempfile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	p := tempfile.Name()
	defer os.Remove(p)
	config := v1.Config{
		Clusters: []v1.NamedCluster{
			{
				Cluster: v1.Cluster{Server: serverURL, CertificateAuthorityData: ca},
			},
		},
		AuthInfos: []v1.NamedAuthInfo{
			{
				AuthInfo: v1.AuthInfo{ClientCertificateData: clientCert, ClientKeyData: clientKey},
			},
		},
	}
	if err := json.NewEncoder(tempfile).Encode(config); err != nil {
		return nil, err
	}

	c, err := tokenReviewInterfaceFromKubeconfig(p)
	if err != nil {
		return nil, err
	}

	return newWithBackoff(c, cacheTime, 0)
}

func TestTLSConfig(t *testing.T) {
	tests := []struct {
		test                            string
		clientCert, clientKey, clientCA []byte
		serverCert, serverKey, serverCA []byte
		wantErr                         bool
	}{
		{
			test:       "TLS setup between client and server",
			clientCert: clientCert, clientKey: clientKey, clientCA: caCert,
			serverCert: serverCert, serverKey: serverKey, serverCA: caCert,
		},
		{
			test:       "Server does not require client auth",
			clientCA:   caCert,
			serverCert: serverCert, serverKey: serverKey,
		},
		{
			test:       "Server does not require client auth, client provides it",
			clientCert: clientCert, clientKey: clientKey, clientCA: caCert,
			serverCert: serverCert, serverKey: serverKey,
		},
		{
			test:       "Client does not trust server",
			clientCert: clientCert, clientKey: clientKey,
			serverCert: serverCert, serverKey: serverKey,
			wantErr: true,
		},
		{
			test:       "Server does not trust client",
			clientCert: clientCert, clientKey: clientKey, clientCA: caCert,
			serverCert: serverCert, serverKey: serverKey, serverCA: badCACert,
			wantErr: true,
		},
		{
			// Plugin does not support insecure configurations.
			test:    "Server is using insecure connection",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		// Use a closure so defer statements trigger between loop iterations.
		func() {
			service := new(mockService)
			service.statusCode = 200

			server, err := NewTestServer(service, tt.serverCert, tt.serverKey, tt.serverCA)
			if err != nil {
				t.Errorf("%s: failed to create server: %v", tt.test, err)
				return
			}
			defer server.Close()

			wh, err := newTokenAuthenticator(server.URL, tt.clientCert, tt.clientKey, tt.clientCA, 0)
			if err != nil {
				t.Errorf("%s: failed to create client: %v", tt.test, err)
				return
			}

			// Allow all and see if we get an error.
			service.Allow()
			_, authenticated, err := wh.AuthenticateToken("t0k3n")
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error making authorization request: %v", err)
				}
				return
			}
			if !authenticated {
				t.Errorf("%s: failed to authenticate token", tt.test)
				return
			}

			service.Deny()
			_, authenticated, err = wh.AuthenticateToken("t0k3n")
			if err != nil {
				t.Errorf("%s: unexpectedly failed AuthenticateToken", tt.test)
			}
			if authenticated {
				t.Errorf("%s: incorrectly authenticated token", tt.test)
			}
		}()
	}
}

// recorderService records all token review requests, and responds with the
// provided TokenReviewStatus.
type recorderService struct {
	lastRequest v1beta1.TokenReview
	response    v1beta1.TokenReviewStatus
}

func (rec *recorderService) Review(r *v1beta1.TokenReview) {
	rec.lastRequest = *r
	r.Status = rec.response
}

func (rec *recorderService) HTTPStatusCode() int { return 200 }

func TestWebhookTokenAuthenticator(t *testing.T) {
	serv := &recorderService{}

	s, err := NewTestServer(serv, serverCert, serverKey, caCert)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	wh, err := newTokenAuthenticator(s.URL, clientCert, clientKey, caCert, 0)
	if err != nil {
		t.Fatal(err)
	}

	expTypeMeta := metav1.TypeMeta{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
	}

	tests := []struct {
		serverResponse        v1beta1.TokenReviewStatus
		expectedAuthenticated bool
		expectedUser          *user.DefaultInfo
	}{
		// Successful response should pass through all user info.
		{
			serverResponse: v1beta1.TokenReviewStatus{
				Authenticated: true,
				User: v1beta1.UserInfo{
					Username: "somebody",
				},
			},
			expectedAuthenticated: true,
			expectedUser: &user.DefaultInfo{
				Name: "somebody",
			},
		},
		{
			serverResponse: v1beta1.TokenReviewStatus{
				Authenticated: true,
				User: v1beta1.UserInfo{
					Username: "person@place.com",
					UID:      "abcd-1234",
					Groups:   []string{"stuff-dev", "main-eng"},
					Extra:    map[string]v1beta1.ExtraValue{"foo": {"bar", "baz"}},
				},
			},
			expectedAuthenticated: true,
			expectedUser: &user.DefaultInfo{
				Name:   "person@place.com",
				UID:    "abcd-1234",
				Groups: []string{"stuff-dev", "main-eng"},
				Extra:  map[string][]string{"foo": {"bar", "baz"}},
			},
		},
		// Unauthenticated shouldn't even include extra provided info.
		{
			serverResponse: v1beta1.TokenReviewStatus{
				Authenticated: false,
				User: v1beta1.UserInfo{
					Username: "garbage",
					UID:      "abcd-1234",
					Groups:   []string{"not-actually-used"},
				},
			},
			expectedAuthenticated: false,
			expectedUser:          nil,
		},
		{
			serverResponse: v1beta1.TokenReviewStatus{
				Authenticated: false,
			},
			expectedAuthenticated: false,
			expectedUser:          nil,
		},
	}
	token := "my-s3cr3t-t0ken"
	for i, tt := range tests {
		serv.response = tt.serverResponse
		user, authenticated, err := wh.AuthenticateToken(token)
		if err != nil {
			t.Errorf("case %d: authentication failed: %v", i, err)
			continue
		}
		if serv.lastRequest.Spec.Token != token {
			t.Errorf("case %d: Server did not see correct token. Got %q, expected %q.",
				i, serv.lastRequest.Spec.Token, token)
		}
		if !reflect.DeepEqual(serv.lastRequest.TypeMeta, expTypeMeta) {
			t.Errorf("case %d: Server did not see correct TypeMeta. Got %v, expected %v",
				i, serv.lastRequest.TypeMeta, expTypeMeta)
		}
		if authenticated != tt.expectedAuthenticated {
			t.Errorf("case %d: Plugin returned incorrect authentication response. Got %t, expected %t.",
				i, authenticated, tt.expectedAuthenticated)
		}
		if user != nil && tt.expectedUser != nil && !reflect.DeepEqual(user, tt.expectedUser) {
			t.Errorf("case %d: Plugin returned incorrect user. Got %#v, expected %#v",
				i, user, tt.expectedUser)
		}
	}
}

type authenticationUserInfo v1beta1.UserInfo

func (a *authenticationUserInfo) GetName() string     { return a.Username }
func (a *authenticationUserInfo) GetUID() string      { return a.UID }
func (a *authenticationUserInfo) GetGroups() []string { return a.Groups }

func (a *authenticationUserInfo) GetExtra() map[string][]string {
	if a.Extra == nil {
		return nil
	}
	ret := map[string][]string{}
	for k, v := range a.Extra {
		ret[k] = []string(v)
	}

	return ret
}

// Ensure v1beta1.UserInfo contains the fields necessary to implement the
// user.Info interface.
var _ user.Info = (*authenticationUserInfo)(nil)

// TestWebhookCache verifies that error responses from the server are not
// cached, but successful responses are. It also ensures that the webhook
// call is retried on 429 and 500+ errors
func TestWebhookCacheAndRetry(t *testing.T) {
	serv := new(mockService)
	s, err := NewTestServer(serv, serverCert, serverKey, caCert)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// Create an authenticator that caches successful responses "forever" (100 days).
	wh, err := newTokenAuthenticator(s.URL, clientCert, clientKey, caCert, 2400*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		description string

		token string
		allow bool
		code  int

		expectError bool
		expectOk    bool
		expectCalls int
	}{
		{
			description: "t0k3n, 500 error, retries and fails",

			token: "t0k3n",
			allow: false,
			code:  500,

			expectError: true,
			expectOk:    false,
			expectCalls: 5,
		},
		{
			description: "t0k3n, 404 error, fails (but no retry)",

			token: "t0k3n",
			allow: false,
			code:  404,

			expectError: true,
			expectOk:    false,
			expectCalls: 1,
		},
		{
			description: "t0k3n, 200 response, allowed, succeeds with a single call",

			token: "t0k3n",
			allow: true,
			code:  200,

			expectError: false,
			expectOk:    true,
			expectCalls: 1,
		},
		{
			description: "t0k3n, 500 response, disallowed, but never called because previous 200 response was cached",

			token: "t0k3n",
			allow: false,
			code:  500,

			expectError: false,
			expectOk:    true,
			expectCalls: 0,
		},

		{
			description: "an0th3r_t0k3n, 500 response, disallowed, should be called again with retries",

			token: "an0th3r_t0k3n",
			allow: false,
			code:  500,

			expectError: true,
			expectOk:    false,
			expectCalls: 5,
		},
		{
			description: "an0th3r_t0k3n, 429 response, disallowed, should be called again with retries",

			token: "an0th3r_t0k3n",
			allow: false,
			code:  429,

			expectError: true,
			expectOk:    false,
			expectCalls: 5,
		},
		{
			description: "an0th3r_t0k3n, 200 response, allowed, succeeds with a single call",

			token: "an0th3r_t0k3n",
			allow: true,
			code:  200,

			expectError: false,
			expectOk:    true,
			expectCalls: 1,
		},
		{
			description: "an0th3r_t0k3n, 500 response, disallowed, but never called because previous 200 response was cached",

			token: "an0th3r_t0k3n",
			allow: false,
			code:  500,

			expectError: false,
			expectOk:    true,
			expectCalls: 0,
		},
	}

	for _, testcase := range testcases {
		func() {
			serv.allow = testcase.allow
			serv.statusCode = testcase.code
			serv.called = 0

			_, ok, err := wh.AuthenticateToken(testcase.token)
			hasError := err != nil
			if hasError != testcase.expectError {
				t.Log(testcase.description)
				t.Errorf("Webhook returned HTTP %d, expected error=%v, but got error %v", testcase.code, testcase.expectError, err)
			}
			if serv.called != testcase.expectCalls {
				t.Log(testcase.description)
				t.Errorf("Expected %d calls, got %d", testcase.expectCalls, serv.called)
			}
			if ok != testcase.expectOk {
				t.Log(testcase.description)
				t.Errorf("Expected ok=%v, got %v", testcase.expectOk, ok)
			}
		}()
	}
}
