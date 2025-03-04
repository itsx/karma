package mapper

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prymitive/karma/internal/models"
)

var (
	alertMappers   = []AlertMapper{}
	silenceMappers = []SilenceMapper{}
	statusMappers  = []StatusMapper{}
)

// Mapper converts Alertmanager response body and maps to karma data structures
type Mapper interface {
	IsSupported(version string) bool
	AbsoluteURL(baseURI string) (string, error)
	QueryArgs() string
	IsOpenAPI() bool
}

// AlertMapper handles mapping of Alertmanager alert information to karma AlertGroup models
type AlertMapper interface {
	Mapper
	Decode(io.ReadCloser) ([]models.AlertGroup, error)
	Collect(string, map[string]string, time.Duration, http.RoundTripper) ([]models.AlertGroup, error)
}

// SilenceMapper handles mapping of Alertmanager silence information to karma Silence models
type SilenceMapper interface {
	Mapper
	Decode(io.ReadCloser) ([]models.Silence, error)
	Collect(string, map[string]string, time.Duration, http.RoundTripper) ([]models.Silence, error)
}

// StatusMapper handles mapping Alertmanager status information containing cluster config
type StatusMapper interface {
	Mapper
	Decode(io.ReadCloser) (models.AlertmanagerStatus, error)
	Collect(string, map[string]string, time.Duration, http.RoundTripper) (models.AlertmanagerStatus, error)
}

// RegisterAlertMapper allows to register mapper implementing alert data
// handling for specific Alertmanager versions
func RegisterAlertMapper(m AlertMapper) {
	alertMappers = append(alertMappers, m)
}

// GetAlertMapper returns mapper for given version
func GetAlertMapper(version string) (AlertMapper, error) {
	for _, m := range alertMappers {
		if m.IsSupported(version) {
			return m, nil
		}
	}
	return nil, fmt.Errorf("can't find alert mapper for Alertmanager %s", version)
}

// RegisterSilenceMapper allows to register mapper implementing silence data
// handling for specific Alertmanager versions
func RegisterSilenceMapper(m SilenceMapper) {
	silenceMappers = append(silenceMappers, m)
}

// GetSilenceMapper returns mapper for given version
func GetSilenceMapper(version string) (SilenceMapper, error) {
	for _, m := range silenceMappers {
		if m.IsSupported(version) {
			return m, nil
		}
	}
	return nil, fmt.Errorf("can't find silence mapper for Alertmanager %s", version)
}

// RegisterStatusMapper allows to register mapper implementing status data
// handling for specific Alertmanager versions
func RegisterStatusMapper(m StatusMapper) {
	statusMappers = append(statusMappers, m)
}

// GetStatusMapper returns mapper for given version
func GetStatusMapper(version string) (StatusMapper, error) {
	for _, m := range statusMappers {
		if m.IsSupported(version) {
			return m, nil
		}
	}
	return nil, fmt.Errorf("can't find status mapper for Alertmanager %s", version)
}
