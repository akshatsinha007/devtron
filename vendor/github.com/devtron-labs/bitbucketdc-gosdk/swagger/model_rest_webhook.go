/*
 * Bitbucket Data Center
 *
 * This is the reference document for the Atlassian Bitbucket REST API. The REST API is for developers who want to:    - integrate Bitbucket with other applications;   - create scripts that interact with Bitbucket; or   - develop plugins that enhance the Bitbucket UI, using REST to interact with the backend.    You can read more about developing Bitbucket plugins in the [Bitbucket Developer Documentation](https://developer.atlassian.com/bitbucket/server/docs/latest/).
 *
 * API version: 8.19
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RestWebhook struct {
	Name                    string                  `json:"name,omitempty"`
	SslVerificationRequired bool                    `json:"sslVerificationRequired,omitempty"`
	Statistics              *interface{}            `json:"statistics,omitempty"`
	ScopeType               string                  `json:"scopeType,omitempty"`
	Configuration           *interface{}            `json:"configuration,omitempty"`
	Credentials             *RestWebhookCredentials `json:"credentials,omitempty"`
	Url                     string                  `json:"url,omitempty"`
	Events                  []string                `json:"events,omitempty"`
	Active                  bool                    `json:"active,omitempty"`
}
