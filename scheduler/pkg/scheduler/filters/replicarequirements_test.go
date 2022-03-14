package filters

import (
	"testing"

	. "github.com/onsi/gomega"
	pb "github.com/seldonio/seldon-core/scheduler/apis/mlops/scheduler"
	"github.com/seldonio/seldon-core/scheduler/pkg/store"
)

func TestReplicaRequirementsFilter(t *testing.T) {
	g := NewGomegaWithT(t)

	getTestModelWithRequirements := func(requirements []string) *store.ModelVersion {
		return store.NewModelVersion(
			&pb.Model{ModelSpec: &pb.ModelSpec{Requirements: requirements}, DeploymentSpec: &pb.DeploymentSpec{Replicas: 1}},
			1,
			"server",
			map[int]store.ReplicaStatus{3: {State: store.Loading}},
			false,
			store.ModelProgressing)
	}

	getTestServerReplicaWithCaps := func(capabilities []string) *store.ServerReplica {
		return store.NewServerReplica("svc", 8080, 5001, 1, nil, capabilities, 100, 100, nil, 100)
	}

	type test struct {
		name     string
		model    *store.ModelVersion
		server   *store.ServerReplica
		expected bool
	}

	tests := []test{
		{name: "Match", model: getTestModelWithRequirements([]string{"sklearn"}), server: getTestServerReplicaWithCaps([]string{"sklearn"}), expected: true},
		{name: "Mismatch", model: getTestModelWithRequirements([]string{"sklearn"}), server: getTestServerReplicaWithCaps([]string{"xgboost"}), expected: false},
		{name: "PartialMatch", model: getTestModelWithRequirements([]string{"sklearn", "xgboost"}), server: getTestServerReplicaWithCaps([]string{"xgboost"}), expected: false},
		{name: "MultiMatch", model: getTestModelWithRequirements([]string{"sklearn", "xgboost"}), server: getTestServerReplicaWithCaps([]string{"xgboost", "sklearn", "tensorflow"}), expected: true},
		{name: "Duplicates", model: getTestModelWithRequirements([]string{"sklearn", "xgboost", "sklearn"}), server: getTestServerReplicaWithCaps([]string{"xgboost", "sklearn", "tensorflow"}), expected: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter := RequirementsReplicaFilter{}
			ok := filter.Filter(test.model, test.server)
			g.Expect(ok).To(Equal(test.expected))
		})
	}
}
