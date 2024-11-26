package api_test

import (
	"fmt"
	"strings"

	ctlappsv1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
)

var _ = Describe("create a supportbundle request and verify taints on daemonset", func() {
	var sbk *cloudweavv1.SupportBundle
	var sbc ctlcloudweavv1.SupportBundleController
	var dsc ctlappsv1.DaemonSetCache
	var sc ctlcloudweavv1.SettingController
	var requiredToleration *corev1.Toleration

	BeforeEach(func() {
		sbk = &cloudweavv1.SupportBundle{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sample",
				Namespace: "cloudweav-system",
			},
			Spec: cloudweavv1.SupportBundleSpec{
				IssueURL:    "fake issue",
				Description: "fake description",
			},
		}

		requiredToleration = &corev1.Toleration{
			Operator: corev1.TolerationOpExists,
		}

		scaled := cloudweav.Scaled()
		sc = scaled.CloudweavFactory.Cloudweavhci().V1beta1().Setting()
		sbc = scaled.CloudweavFactory.Cloudweavhci().V1beta1().SupportBundle()
		dsc = scaled.AppsFactory.Apps().V1().DaemonSet().Cache()

		Eventually(func() error {
			_, err := sbc.Create(sbk)
			return err
		}, "30s", "5s").ShouldNot(HaveOccurred())
	})

	It("check a daemonset exists with correct annotations", func() {
		By("checking support-bundle-image setting is populated", func() {
			Eventually(func() error {
				sbi, err := sc.Get("support-bundle-image", metav1.GetOptions{})
				if err != nil {
					return err
				}
				sbi.Value = `{"repository":"rancher/support-bundle-kit","tag":"master-head","imagePullPolicy":"IfNotPresent"}`
				_, err = sc.Update(sbi)
				return err
			}, "30s", "5s").ShouldNot(HaveOccurred())
		})

		By("checking ds is created", func() {
			Eventually(func() error {
				dsList, err := dsc.List("cloudweav-system", labels.NewSelector())
				if err != nil {
					return err
				}

				for _, v := range dsList {
					if strings.Contains(v.Name, "supportbundle") && v.Spec.Template.Spec.Tolerations[0].MatchToleration(requiredToleration) {
						return nil
					}
				}

				return fmt.Errorf("waiting for ds to be created or toleration to exist")
			}, "90s", "5s").ShouldNot(HaveOccurred())
		})

	})
	AfterEach(func() {
		Eventually(func() error {
			err := sbc.Delete(sbk.Namespace, sbk.Name, &metav1.DeleteOptions{})
			return err
		}, "30s", "5s").ShouldNot(HaveOccurred())

	})
})
