package interfaces

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/mapstructure"

	"terraform-provider-ansible-forms/internal/restclient"
	"terraform-provider-ansible-forms/internal/utils"
)

// JobGetDataSourceModel describes the data source model.
type JobGetDataSourceModel struct {
	Status  string `mapstructure:"status"`
	Message string `mapstructure:"message"`
	Data    struct {
		ID            int64  `mapstructure:"id"`
		Form          string `mapstructure:"form"`
		Target        string `mapstructure:"target"`
		Status        string `mapstructure:"status"`
		Start         string `mapstructure:"start"`
		End           string `mapstructure:"end"`
		User          string `mapstructure:"user"`
		UserType      string `mapstructure:"user_type"`
		JobType       string `mapstructure:"job_type"`
		Extravars     string `mapstructure:"extravars"`
		Credentials   string `mapstructure:"credentials"`
		Notifications string `mapstructure:"notifications"`
		NoOfRecords   int64  `mapstructure:"no_of_records"`
		Counter       int64  `mapstructure:"counter"`
		Output        string `mapstructure:"output"`
	} `mapstructure:"data"`
}

// GetJobById ...
func GetJobById(errorHandler *utils.ErrorHandler, r restclient.RestClient, id string) (*JobGetDataSourceModel, error) {
	statusCode, response, err := r.GetNilOrOneRecord("job/"+id, nil, nil)
	if err != nil {
		return nil, errorHandler.MakeAndReportError("error reading job info", fmt.Sprintf("error on GET job/: %s, statusCode %d", err, statusCode))
	}

	var apiResp *JobGetDataSourceModel
	if err = mapstructure.Decode(response, &apiResp); err != nil {
		return nil, errorHandler.MakeAndReportError("failed to decode response from GET job", fmt.Sprintf("error: %s, statusCode %d, response %#v", err, statusCode, response))
	}
	tflog.Debug(errorHandler.Ctx, fmt.Sprintf("read job info: %#v", apiResp.Data))

	return apiResp, nil
}
