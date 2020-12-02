package agent

import (
	"github.com/aaqaishtyaq/gotopdf/topdf"
	log "github.com/sirupsen/logrus"
)

// GeneratePdfRequest types
type GeneratePdfRequest struct {
	FileName string `json:"fileName,omitempty"`
	TargetURL string `json:"targetURL"`
	Headers map[string]string `json:"headers,omitempty"`
	Orientation string `json:"orientation"`
	PrintBackground bool `json:"printBackground"`
	MarginTop float64 `json:"marginTop"`
	MarginRight float64 `json:"marginRight"`
	MarginBottom float64 `json:"marginBottom"`
	MarginLeft float64 `json:"marginLeft"`
}

// DefaultGeneratePdfRequest generates default options for the request
func DefaultGeneratePdfRequest() GeneratePdfRequest {
	return GeneratePdfRequest {
		Orientation: "Portrait",
		PrintBackground: true,
		MarginTop: 0.4,
		MarginRight: 0.4,
		MarginBottom: 0.4,
		MarginLeft: 0.4,
	}
}

// GeneratePDF using topdf
func GeneratePDF(fileName string) {
	options := DefaultGeneratePdfRequest()
	
}