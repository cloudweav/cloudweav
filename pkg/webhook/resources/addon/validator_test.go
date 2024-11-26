package addon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
)

func Test_validateUpdatedAddon(t *testing.T) {
	var testCases = []struct {
		name          string
		oldAddon      *cloudweavv1.Addon
		newAddon      *cloudweavv1.Addon
		expectedError bool
	}{
		{
			name: "user can enable addon",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       false,
					ValuesContent: "sample",
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			expectedError: false,
		},
		{
			name: "user can disable addon",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       false,
					ValuesContent: "sample",
				},
			},
			expectedError: false,
		},
		{
			name: "user can't change chart field",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1-changed",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			expectedError: true,
		},
		{
			name: "user can't change disabling addon",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       false,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonDisabling,
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1-changed",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonDisabling,
				},
			},
			expectedError: true,
		},
		{
			name: "user can disable deployed addon",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonDeployed,
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       false,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonDeployed,
				},
			},
			expectedError: false,
		},
		{
			name: "user can't disable enabling addon",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "disable-enabling-addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonEnabling,
					Conditions: []cloudweavv1.Condition{
						{
							Type:   cloudweavv1.AddonOperationInProgress,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "disable-enabling-addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       false,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonEnabling,
					Conditions: []cloudweavv1.Condition{
						{
							Type:   cloudweavv1.AddonOperationInProgress,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
			expectedError: true,
		},
		{
			name: "user can change addon annotations when addon is being enabled",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "change-enabling-addon1-annotation",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonDeployed,
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "change-enabling-addon1-annotation",
					Annotations: map[string]string{
						"cloudweavhci.io/addon-operation-timeout": "2",
					},
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonDeployed,
				},
			},
			expectedError: false,
		},
		{
			name: "user can disable deployfailed addon",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "disable-deployfailed-addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonEnabling,
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "disable-deployfailed-addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       false,
					ValuesContent: "sample",
				},
				Status: cloudweavv1.AddonStatus{
					Status: cloudweavv1.AddonEnabling,
					Conditions: []cloudweavv1.Condition{
						{
							Type:   cloudweavv1.AddonOperationFailed,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "virtual cluster addon with valid dns",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name:      vClusterAddonName,
					Namespace: vClusterAddonNamespace,
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "vcluster",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name:      vClusterAddonName,
					Namespace: vClusterAddonNamespace,
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "vcluster",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "hostname: rancher.172.19.108.3.sslip.io\nrancherVersion: v2.7.4\nbootstrapPassword: cloudweavAdmin\n",
				},
			},
			expectedError: false,
		},
		{
			name: "virtual cluster addon with ingress-expose address",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name:      vClusterAddonName,
					Namespace: vClusterAddonNamespace,
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "vcluster",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name:      vClusterAddonName,
					Namespace: vClusterAddonNamespace,
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "vcluster",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "hostname: 172.19.108.3\nrancherVersion: v2.7.4\nbootstrapPassword: cloudweavAdmin\n",
				},
			},
			expectedError: true,
		},
		{
			name: "virtual cluster addon with invalid fqdn",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name:      vClusterAddonName,
					Namespace: vClusterAddonNamespace,
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "vcluster",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name:      vClusterAddonName,
					Namespace: vClusterAddonNamespace,
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "vcluster",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "hostname: FakeAddress.com\nrancherVersion: v2.7.4\nbootstrapPassword: cloudweavAdmin\n",
				},
			},
			expectedError: true,
		},
		{
			name: "virtual cluster addon empty hostname",
			oldAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name:      vClusterAddonName,
					Namespace: vClusterAddonNamespace,
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "vcluster",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name:      vClusterAddonName,
					Namespace: vClusterAddonNamespace,
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "vcluster",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "hostname: \nrancherVersion: v2.7.4\nbootstrapPassword: cloudweavAdmin\n",
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		err := validateUpdatedAddon(tc.newAddon, tc.oldAddon)
		if tc.expectedError {
			assert.NotNil(t, err, tc.name)
		} else {
			assert.Nil(t, err, tc.name)
		}
	}
}

func Test_validateNewAddon(t *testing.T) {
	var testCases = []struct {
		name          string
		newAddon      *cloudweavv1.Addon
		addonList     []*cloudweavv1.Addon
		expectedError bool
	}{
		{
			name: "user can add new addon",
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			addonList:     []*cloudweavv1.Addon{},
			expectedError: false,
		},
		{
			name: "user cannot add same addon, no matter differences in version and repo fields",
			newAddon: &cloudweavv1.Addon{
				ObjectMeta: metav1.ObjectMeta{
					Name: "addon1",
				},
				Spec: cloudweavv1.AddonSpec{
					Repo:          "repo1",
					Chart:         "chart1",
					Version:       "version1",
					Enabled:       true,
					ValuesContent: "sample",
				},
			},
			addonList: []*cloudweavv1.Addon{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "addon1",
					},
					Spec: cloudweavv1.AddonSpec{
						Repo:          "repo1",
						Chart:         "chart1",
						Version:       "version1",
						Enabled:       true,
						ValuesContent: "sample",
					},
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		err := validateNewAddon(tc.newAddon, tc.addonList)
		if tc.expectedError {
			assert.NotNil(t, err, tc.name)
		} else {
			assert.Nil(t, err, tc.name)
		}
	}
}
