package provider

// AnsibleFormsJob -
// type AnsibleFormsJob struct {
// 	ID    int         `json:"id,omitempty"`
// 	FormName string   `json:"form_name,omitempty"`
// 	JsonParams string `json:"json_params,omitempty"`
// 	Status string 	  `json:"status,omitempty"`
// }

type AnsibleFormsCreateJob struct {
	FormName string   `json:"formName"`
	ExtraVars AnsibleFormsExtraVars  `json:"extravars"`
	Credentials AnsibleFormsCredentials `json:"credentials"`
}

type AnsibleFormsExtraVars struct {
	Region	string `json:"region"`
	Opco	string `json:"opco"`
    SvmName string `json:"svm_name"`
	State string `json:"state"`
	Exposure string `json:"exposure"`
	Env string `json:"env"`
	Dataclass string `json:"dataclass"`
	ShareName string `json:"share_name"`
	AccountID string `json:"accountid"`
	Size string `json:"size"`
	ProtectionRequired string `json:"protection_required"`
}

type AnsibleFormsCredentials struct {     
    OntapCred string `json:"ontap_cred,omitempty"`
    BindCred string `json:"bind_cred,omitempty"`
}

//type AnsibleFormsCreateJobResponse struct {
//	Status string 	  `json:"status"`
//	Message string 	  `json:"message"`
//	Data  AnsibleFormsJobResponseData `json:"data,omitempty"`
//}

// type AnsibleFormsJobResponseData struct {
// 	Output 	AnsibleFormsJobResponseDataOutput	`json:"output,omitempty"`
// 	Error string 	  `json:"error,omitempty"`
// }

// type AnsibleFormsJobResponseDataOutput struct {
// 	ID		int `json:"id"`
// }

// type AnsibleFormsJobResponse struct {
// 	Status string 	  `json:"status"`
// 	Message string 	  `json:"message"`
// 	Error string 	  `json:"error,omitempty"`
// 	Data  AnsibleFormsJobData `json:"data"`
// }

// type AnsibleFormsJobData struct {
// 	ID		int `json:"id"`
// 	FormName string `json:"form,omitempty"`
// 	Target string `json:"target,omitempty"`
// 	Status string `json:"status,omitempty"`
// 	Start string `json:"start,omitempty"`
// 	End string `json:"end,omitempty"`
// 	User string `json:"user,omitempty"`
// 	UserType string `json:"user_type,omitempty"`
// 	JobType string `json:"job_type,omitempty"`
// 	ParentID string `json:"parent_id,omitempty"`
// 	Approval string `json:"approval,omitempty"`
// 	JsonParams string `json:"extravars,omitempty"`
// }

 type AnsibleFormsCreateJobResponse struct {
 	Status string 	  `json:"status"`
 	Message string 	  `json:"message"`
 	Data  AnsibleFormsCreateJobResponseData `json:"data"`
 }

type AnsibleFormsCreateJobResponseData struct {
 	Output 	AnsibleFormsJobResponseDataOutput	`json:"output,omitempty"`
 	Error string 	  `json:"error,omitempty"`
}

 type AnsibleFormsJobResponseDataOutput struct {
 	ID		int `json:"id"`
 	Form string `json:"form,omitempty"`
 	Target string `json:"target,omitempty"`
 	Status string `json:"status,omitempty"`
 	Start string `json:"start,omitempty"`
 	End string `json:"end,omitempty"`
 	User string `json:"user,omitempty"`
 	UserType string `json:"user_type,omitempty"`
 	JobType string `json:"job_type,omitempty"`
 	ParentID string `json:"parent_id,omitempty"`
 	Approval string `json:"approval,omitempty"`
 	JsonParams string `json:"extravars,omitempty"`
}

type AnsibleFormsGetJobResponse struct {
	Status string 	  `json:"status"`
	Message string 	  `json:"message"`
	Data  AnsibleFormsJobResponseDataOutput `json:"data"`
}


