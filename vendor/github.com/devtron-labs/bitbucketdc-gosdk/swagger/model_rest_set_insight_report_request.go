/*
 * Bitbucket Data Center
 *
 * This is the reference document for the Atlassian Bitbucket REST API. The REST API is for developers who want to:    - integrate Bitbucket with other applications;   - create scripts that interact with Bitbucket; or   - develop plugins that enhance the Bitbucket UI, using REST to interact with the backend.    You can read more about developing Bitbucket plugins in the [Bitbucket Developer Documentation](https://developer.atlassian.com/bitbucket/server/docs/latest/).
 *
 * API version: 8.19
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RestSetInsightReportRequest struct {
	CoverageProviderKey string                  `json:"coverageProviderKey,omitempty"`
	CreatedDate         int64                   `json:"createdDate,omitempty"`
	Data                []RestInsightReportData `json:"data"`
	Details             string                  `json:"details,omitempty"`
	Link                string                  `json:"link,omitempty"`
	LogoUrl             string                  `json:"logoUrl,omitempty"`
	Reporter            string                  `json:"reporter,omitempty"`
	Result              string                  `json:"result,omitempty"`
	Title               string                  `json:"title"`
}
