package server

import (
	"net/http"
	"reflect"

	"github.com/rancher/wrangler/v3/pkg/webhook"
	"github.com/sirupsen/logrus"

	"github.com/cloudweav/cloudweav/pkg/webhook/clients"
	"github.com/cloudweav/cloudweav/pkg/webhook/config"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/pod"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/templateversion"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/virtualmachine"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/virtualmachinebackup"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/virtualmachineimage"
	"github.com/cloudweav/cloudweav/pkg/webhook/types"
)

func Mutation(clients *clients.Clients, options *config.Options) (http.Handler, []types.Resource, error) {
	resources := []types.Resource{}
	settingCache := clients.CloudweavFactory.Cloudweavhci().V1beta1().Setting().Cache()
	storageClassCache := clients.StorageFactory.Storage().V1().StorageClass().Cache()
	nadCache := clients.CNIFactory.K8s().V1().NetworkAttachmentDefinition().Cache()
	vmBackupCache := clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineBackup().Cache()
	mutators := []types.Mutator{
		pod.NewMutator(settingCache),
		templateversion.NewMutator(),
		virtualmachine.NewMutator(settingCache, nadCache),
		virtualmachineimage.NewMutator(storageClassCache),
		virtualmachinebackup.NewMutator(vmBackupCache),
	}

	router := webhook.NewRouter()
	for _, m := range mutators {
		addHandler(router, types.AdmissionTypeMutation, m, options)
		resources = append(resources, m.Resource())
	}

	return router, resources, nil
}

func addHandler(router *webhook.Router, admissionType string, admitter types.Admitter, options *config.Options) {
	rsc := admitter.Resource()
	kind := reflect.Indirect(reflect.ValueOf(rsc.ObjectType)).Type().Name()
	router.Kind(kind).Group(rsc.APIGroup).Type(rsc.ObjectType).Handle(types.NewAdmissionHandler(admitter, admissionType, options))
	logrus.Infof("add %s handler for %+v.%s (%s)", admissionType, rsc.Names, rsc.APIGroup, kind)
}
