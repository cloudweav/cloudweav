package upgrade

import (
	"testing"

	upgradeapiv1 "github.com/rancher/system-upgrade-controller/pkg/apis/upgrade.cattle.io/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/fake"
	"github.com/cloudweav/cloudweav/pkg/util/fakeclients"
)

func TestPodHandler_OnChanged(t *testing.T) {
	type input struct {
		key     string
		pod     *corev1.Pod
		plan    *upgradeapiv1.Plan
		upgrade *cloudweavv1.Upgrade
	}
	type output struct {
		upgrade *cloudweavv1.Upgrade
		err     error
	}
	var testCases = []struct {
		name     string
		given    input
		expected output
	}{
		{
			name: "upgrade repo vm ready",
			given: input{
				key: "upgrade-repo-vm-pod",
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "upgrade-repo-vm-pod",
						Namespace: upgradeNamespace,
						Labels: map[string]string{
							cloudweavUpgradeLabel:          testUpgradeName,
							cloudweavUpgradeComponentLabel: upgradeComponentRepo,
						},
					},
					Status: corev1.PodStatus{
						ContainerStatuses: []corev1.ContainerStatus{
							{
								Ready: true,
							},
						},
					},
				},
				plan:    newTestPlanBuilder().Build(),
				upgrade: newTestUpgradeBuilder().WithLabel(upgradeStateLabel, StatePreparingRepo).Build(),
			},
			expected: output{
				upgrade: newTestUpgradeBuilder().WithLabel(upgradeStateLabel, StateRepoPrepared).RepoProvisionedCondition(corev1.ConditionTrue, "", "").Build(),
				err:     nil,
			},
		},
	}
	for _, tc := range testCases {
		var clientset = fake.NewSimpleClientset(tc.given.plan, tc.given.upgrade)
		var handler = &podHandler{
			namespace:     cloudweavSystemNamespace,
			planCache:     fakeclients.PlanCache(clientset.UpgradeV1().Plans),
			upgradeClient: fakeclients.UpgradeClient(clientset.CloudweavhciV1beta1().Upgrades),
			upgradeCache:  fakeclients.UpgradeCache(clientset.CloudweavhciV1beta1().Upgrades),
		}
		var actual output
		var getErr error
		_, actual.err = handler.OnChanged(tc.given.key, tc.given.pod)
		actual.upgrade, getErr = handler.upgradeCache.Get(handler.namespace, tc.given.upgrade.Name)
		assert.Nil(t, getErr)

		emptyConditionsTime(tc.expected.upgrade.Status.Conditions)
		emptyConditionsTime(actual.upgrade.Status.Conditions)

		assert.Equal(t, tc.expected, actual, "case %q", tc.name)
	}
}
