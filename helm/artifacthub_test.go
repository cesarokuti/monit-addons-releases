package helm

import (
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestArtifactHub(t *testing.T) {
	v, err := ArtifactHub("oci-karpenter/karpenter")
	if err != nil {
		t.Errorf("not get artifacthub data %v", err)
	}
	_, err = semver.NewVersion(v)
	if err != nil {
		t.Errorf("%v is not a semver version: %v", v, err)
	}

}
