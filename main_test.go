package main

import (
	"testing"

	"github.com/alexflint/go-arg"
	"github.com/stretchr/testify/assert"
)

func TestParamsParser_getQueryParams(t *testing.T) {
	tests := []struct {
		name    string
		params  []string
		want    *GKELogQueryParams
		wantErr error
	}{
		{
			name:   "parse success with gke context",
			params: []string{"-g", "gke_staging_location_cluster-1", "-t", "k8s_container", "-n", "staging", "-c", "staging-xxx-xxx-xxx-xxx"},
			want: &GKELogQueryParams{
				ProjectID:     "staging",
				Location:      "location",
				ClusterName:   "cluster-1",
				GkeContext:    "gke_staging_location_cluster-1",
				ResourceType:  "k8s_container",
				NamespaceName: "staging",
				ContainerName: "staging-xxx-xxx-xxx-xxx",
			},
			wantErr: nil,
		},
		{
			name:    "parse failed with invalid gke context",
			params:  []string{"-g", "gke_xxx", "-t", "k8s_container", "-n", "staging", "-P", "staging-xxx-xxx-xxx-xxx"},
			want:    nil,
			wantErr: errNotGKEContext,
		},
		{
			name:    "parse failed with context is not gke",
			params:  []string{"-g", "gek_xxx_xxx_xxx", "-t", "k8s_container", "-n", "staging", "-P", "staging-xxx-xxx-xxx-xxx"},
			want:    nil,
			wantErr: errNotGKEContext,
		},
		{
			name:    "parse failed with empty gke context",
			params:  []string{"-g", "", "-t", "k8s_container", "-n", "staging", "-P", "staging-xxx-xxx-xxx-xxx"},
			want:    nil,
			wantErr: errNotGKEContext,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ParamsParser{
				parseArgs: func(q *GKELogQueryParams) {
					parser, _ := arg.NewParser(arg.Config{}, q)
					_ = parser.Parse(tt.params)
				},
			}
			got, err := p.getQueryParams()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestQueryParams_getGKEQuery(t *testing.T) {
	tests := []struct {
		name        string
		queryParams *GKELogQueryParams
		want        string
		wantErr     error
	}{
		{
			name: "get gke query success",
			queryParams: &GKELogQueryParams{
				ProjectID:     "staging",
				Location:      "location",
				ClusterName:   "cluster-1",
				ResourceType:  "k8s_container",
				NamespaceName: "staging",
				ContainerName: "staging-xxx-xxx-xxx-xxx",
			},
			want:    "https://console.cloud.google.com/logs/query;query=resource.labels.cluster_name=cluster-1%0Aresource.labels.container_name=staging-xxx-xxx-xxx-xxx%0Aresource.labels.location=location%0Aresource.labels.namespace_name=staging%0Aresource.labels.project_id=staging%0Aresource.type=k8s_container?project=staging",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := tt.queryParams
			got, err := qp.getGKELogQuery()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
