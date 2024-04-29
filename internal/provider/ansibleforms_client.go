// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// HostURL - Default AnsibleForms URL
const HostURL string = "https://localhost:8443"

// Client -
type AnsibleFormsClient struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	Token string `json:"token"`
}

// NewClient -
func ansibleformsNewClient(ctx context.Context, host, username, password *string) (*AnsibleFormsClient, error) {
	c := AnsibleFormsClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
		// Default AnsibleForms URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	ar, err := c.SignIn(ctx)
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token

	return &c, nil
}

func (c *AnsibleFormsClient) doRequest(ctx context.Context, req *http.Request, authToken *string) ([]byte, error) {

	if authToken != nil {
		token := *authToken
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		if req.Method == "POST" {
			req.Header.Set("Content-type", "application/json")
		}

	} else {
		if c.Auth.Username != "" || c.Auth.Password != "" {
			req.Header.Set("Authorization", "Basic "+basicAuth(c.Auth.Username, c.Auth.Password))
		}

	}
	ctx = tflog.SetField(ctx, "requestURL", req.URL)

	tflog.Debug(ctx, "Launching AnsibleFormsClient.doRequest")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		ctx = tflog.SetField(ctx, "statusCode", res.StatusCode)
		ctx = tflog.SetField(ctx, "Body", body)
		tflog.Debug(ctx, "Response AnsibleFormsClient.doRequest is NOT OK")
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	ctx = tflog.SetField(ctx, "statusCode", res.StatusCode)
	ctx = tflog.SetField(ctx, "Body", body)
	tflog.Debug(ctx, "Response AnsibleFormsClient.doRequest is OK")
	return body, err
}

// SignIn - Get a new token for user
func (c *AnsibleFormsClient) SignIn(ctx context.Context) (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/auth/login", c.HostURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(c.Auth.Username, c.Auth.Password))

	body, err := c.doRequest(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// GetJob - Returns a specifc job
// func (c *AnsibleFormsClient) GetJob(ctx context.Context, jobID string) (*AnsibleFormsJobResponse, error) {
func (c *AnsibleFormsClient) GetJob(ctx context.Context, jobID string) (*AnsibleFormsGetJobResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/job/%s", c.HostURL, jobID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req, &c.Token)
	if err != nil {
		return nil, err
	}

	ctx = tflog.SetField(ctx, "BodyReturnJob", string(body))
	tflog.Debug(ctx, "Receive Body Job")

	job := AnsibleFormsGetJobResponse{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// CreateJob - Create new job
func (c *AnsibleFormsClient) CreateJob(ctx context.Context, createjob AnsibleFormsCreateJob) (*AnsibleFormsCreateJobResponse, error) {

	rb, err := json.Marshal(createjob)
	if err != nil {
		return nil, err
	}
	ctx = tflog.SetField(ctx, "BodyJob", string(rb))
	tflog.Debug(ctx, "Send Body Job")
	//myTempJson := "{\"formName\":\"AXA Share Create\",\"extravars\":{\"region\":\"myegion\",\"opco\":\"myopco\",\"svm_name\":\"mysvm_name\",\"state\":\"mystate\",\"exposure\":\"myexposure\",\"env\":\"myenv\",\"dataclass\":\"mydataclass\",\"share_name\":\"myshare_name\",\"accountid\":\"myaccountid\",\"size\":\"mysize\",\"protection_required\":\"myprotection_required\"}}"
	//req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/job", c.HostURL), strings.NewReader(string(myTempJson)))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/job", c.HostURL), strings.NewReader(string(rb)))

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(ctx, req, &c.Token)
	if err != nil {
		return nil, err
	}

	ctx = tflog.SetField(ctx, "BodyReturnJob", string(body))
	tflog.Debug(ctx, "Send Body Job")

	job := AnsibleFormsCreateJobResponse{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	ctx = tflog.SetField(ctx, "Body", job)
	tflog.Debug(ctx, "Response Body OUTPUT Job")
	return &job, nil
}
