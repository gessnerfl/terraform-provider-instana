package restapi

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const (
	//SliConfigResourcePath path to sli config resource of Instana RESTful API
	SliConfigResourcePath = SettingsBasePath + "/sli"
)

//MetricConfiguration represents the nested object metric configuration of the sli config REST resource at Instana
type MetricConfiguration struct {
	Name        string  `json:"metricName"`
	Aggregation string  `json:"metricAggregation"`
	Threshold   float64 `json:"threshold"`
}

//Validate implemention of the interface InstanaDataObject for MetricConfiguration
func (m MetricConfiguration) Validate() error {
	if utils.IsBlank(m.Name) {
		return errors.New("metric name is missing")
	}
	if utils.IsBlank(m.Aggregation) {
		return errors.New("metric aggregation is missing")
	}
	if m.Threshold <= 0.0 {
		return errors.New("threshold must be higher than 0.0")
	}
	return nil
}

//SliEntity represents the nested object sli entity of the sli config REST resource at Instana
type SliEntity struct {
	Type          string `json:"sliType"`
	ApplicationID string `json:"applicationId"`
	ServiceID     string `json:"serviceId"`
	EndpointID    string `json:"endpointId"`
	BoundaryScope string `json:"boundaryScope"`
}

//Validate implemention of the interface InstanaDataObject for SliEntity
func (s SliEntity) Validate() error {
	if utils.IsBlank(s.Type) {
		return errors.New("sli type is missing")
	}
	if utils.IsBlank(s.BoundaryScope) {
		return errors.New("boundary scope is missing")
	}
	return nil
}

//SliConfig represents the REST resource of sli configuration at Instana
type SliConfig struct {
	ID                         string              `json:"id"`
	Name                       string              `json:"sliName"`
	InitialEvaluationTimestamp int                 `json:"initialEvaluationTimestamp"`
	MetricConfiguration        MetricConfiguration `json:"metricConfiguration"`
	SliEntity                  SliEntity           `json:"sliEntity"`
}

//GetIDForResourcePath implemention of the interface InstanaDataObject
func (s *SliConfig) GetIDForResourcePath() string {
	return s.ID
}

//Validate implemention of the interface InstanaDataObject for SliConfig
func (s *SliConfig) Validate() error {
	if utils.IsBlank(s.ID) {
		return errors.New("id is missing")
	}
	if utils.IsBlank(s.Name) {
		return errors.New("sli name is missing")
	}

	if err := s.MetricConfiguration.Validate(); err != nil {
		return err
	}

	if err := s.SliEntity.Validate(); err != nil {
		return err
	}

	return nil
}
