package helm

import (
	"testing"
)

func TestGetChartFile(t *testing.T) {
	file := []byte(`
apiVersion: v2
name: argo-cd
description: App of Apps for argocd
type: application
version: 0.0.1
dependencies:
- name: "argo-cd"
  version: 5.34.1
  repository: "https://argoproj.github.io/argo-helm"
- name: datadog-monitors
  version: 0.0.0-latest
  repository: oci://557130146574.dkr.ecr.us-east-1.amazonaws.com`)
	y := string(file[:])
	_, err := GetChartFile(y)
	if err != nil {
		t.Errorf("error to convert Chart.yaml: %v", err)
	}
}

func TestChartVersion(t *testing.T) {
	r := "https://argoproj.github.io/argo-helm"
	n := "argo-cd"
	_, err := ChartVersion(r, n)
	if err != nil {
		t.Errorf("error to get chart latest version: %v", err)
	}
}

func TestVersionCompare(t *testing.T) {
	v1 := [3]string{"3.2.1", "3.3.3", "5.4.3"}
	v2 := [3]string{"1.2.1", "3.5.3", "5.4.3"}
	var r [3]bool
	for i := range v1 {
		r[i] = VersionCompare(v1[i], v2[i])
	}
	if r[0] != true {
		t.Errorf(" response: %v is not a expected response true", r[0])
	} else if r[1] != false {
		t.Errorf(" response: %v is not a expected response true", r[0])
	} else if r[2] != false {
		t.Errorf(" response: %v is not a expected response true", r[0])
	}
}
