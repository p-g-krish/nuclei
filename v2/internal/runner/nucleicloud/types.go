package nucleicloud

import (
	"time"

	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/contextargs"
)

// AddScanRequest is a nuclei scan input item.
type AddScanRequest struct {
	// RawTargets is a list of raw target URLs for the scan.
	RawTargets []*contextargs.MetaInput `json:"raw_targets,omitempty"`
	// PublicTemplates is a list of public templates for the scan
	PublicTemplates []string `json:"public_templates,omitempty"`
	// PrivateTemplates is a map of template-name->contents that
	// are private to the user executing the scan. (TODO: TBD)
	PrivateTemplates map[string]string `json:"private_templates,omitempty"`
	IsTemporary      bool              `json:"is_temporary"`
}

type GetResultsResponse struct {
	Finished bool                     `json:"finished"`
	Items    []GetResultsResponseItem `json:"items"`
}

type GetScanRequest struct {
	Id         string    `json:"id"`
	Total      int32     `json:"total"`
	Current    int32     `json:"current"`
	Finished   bool      `json:"finished"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
	Targets    int32     `json:"targets"`
	Templates  int32     `json:"templates"`
	Matches    int64     `json:"matches"`
}

type GetResultsResponseItem struct {
	ID  int64  `json:"id"`
	Raw string `json:"raw"`
}

type DeleteScanResults struct {
	OK bool `json:"ok"`
}
