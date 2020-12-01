package restapi

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const (
	SliConfigResourcePath = SettingsBasePath + "/sli"
)

type MetricConfiguration struct {
	Name        string  `json:"metricName"`
	Aggregation string  `json:"metricAggregation"`
	Threshold   float64 `json:"threshold"`
}

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

type SliEntity struct {
	Type          string `json:"sliType"`
	ApplicationID string `json:"applicationId"`
	ServiceID     string `json:"serviceId"`
	EndpointID    string `json:"endpointId"`
	BoundaryScope string `json:"boundaryScope"`
}

func (s SliEntity) Validate() error {
	if utils.IsBlank(s.Type) {
		return errors.New("sli type is missing")
	}
	if utils.IsBlank(s.BoundaryScope) {
		return errors.New("boundary scope is missing")
	}
	return nil
}

type SliConfig struct {
	ID                         string              `json:"id"`
	Name                       string              `json:"sliName"`
	InitialEvaluationTimestamp int                 `json:"initialEvaluationTimestamp"`
	MetricConfiguration        MetricConfiguration `json:"metricConfiguration"`
	SliEntity                  SliEntity           `json:"sliEntity"`
}

func (s SliConfig) GetID() string {
	return s.ID
}

func (s SliConfig) Validate() error {
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
