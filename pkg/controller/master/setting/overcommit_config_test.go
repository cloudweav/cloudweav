package setting

import (
	"encoding/json"
	"strconv"
	"testing"

	lhv1beta2 "github.com/longhorn/longhorn-manager/k8s/pkg/apis/longhorn/v1beta2"
	longhorn "github.com/longhorn/longhorn-manager/types"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/fake"
	"github.com/cloudweav/cloudweav/pkg/settings"
	"github.com/cloudweav/cloudweav/pkg/util"
	"github.com/cloudweav/cloudweav/pkg/util/fakeclients"
)

func TestHandler_syncOvercommitConfig(t *testing.T) {
	type input struct {
		setting *cloudweavv1.Setting
	}
	type output struct {
		value string
		err   error
	}

	const (
		namespace = "default-test"
	)

	createHandler := func(clientset *fake.Clientset) *Handler {
		return &Handler{
			namespace:            namespace,
			longhornSettings:     fakeclients.LonghornSettingClient(clientset.LonghornV1beta2().Settings),
			longhornSettingCache: fakeclients.LonghornSettingCache(clientset.LonghornV1beta2().Settings),
		}
	}

	t.Run("test overcommit-config value", func(t *testing.T) {
		// arrange
		clientset := fake.NewSimpleClientset()
		longhornSettingName := string(longhorn.SettingNameStorageOverProvisioningPercentage)
		handler := createHandler(clientset)
		originalSetting := &lhv1beta2.Setting{ObjectMeta: metav1.ObjectMeta{Namespace: util.LonghornSystemNamespaceName, Name: longhornSettingName}}
		err := clientset.Tracker().Add(originalSetting)
		assert.Nil(t, err, "mock resource should add into fake controller tracker")
		inputSetting := &cloudweavv1.Setting{
			ObjectMeta: metav1.ObjectMeta{Name: settings.OvercommitConfigSettingName},
			Value:      `{"cpu":1300,"memory":1200,"storage":1100}`,
		}
		expected := settings.Overcommit{
			CPU:     1300,
			Memory:  1200,
			Storage: 1100,
		}

		// act
		err = handler.syncOvercommitConfig(inputSetting)

		// assert
		assert.Nil(t, err, "mock resource should get from fake controller")
		lhsetting, err := handler.longhornSettings.Get(util.LonghornSystemNamespaceName, longhornSettingName, metav1.GetOptions{})
		assert.Nil(t, err, "mock resource should get from fake controller")
		assert.Equal(t, lhsetting.Value, strconv.Itoa(expected.Storage), "storage not equals")
	})

	t.Run("test overcommit-config default", func(t *testing.T) {
		// arrange
		clientset := fake.NewSimpleClientset()
		longhornSettingName := string(longhorn.SettingNameStorageOverProvisioningPercentage)
		handler := createHandler(clientset)
		originalSetting := &lhv1beta2.Setting{ObjectMeta: metav1.ObjectMeta{Namespace: util.LonghornSystemNamespaceName, Name: longhornSettingName}}
		err := clientset.Tracker().Add(originalSetting)
		assert.Nil(t, err, "mock resource should add into fake controller tracker")

		inputSetting := &cloudweavv1.Setting{
			ObjectMeta: metav1.ObjectMeta{Name: settings.OvercommitConfigSettingName},
			Default:    settings.OvercommitConfig.Default,
		}

		var expected *settings.Overcommit
		assert.Nil(t, json.Unmarshal([]byte(settings.OvercommitConfig.Default), &expected), "json unmarshal overcommit-config failed")

		// act
		err = handler.syncOvercommitConfig(inputSetting)

		// assert
		assert.Nil(t, err, "mock resource should get from fake controller")
		lhsetting, err := handler.longhornSettings.Get(util.LonghornSystemNamespaceName, longhornSettingName, metav1.GetOptions{})
		assert.Nil(t, err, "mock resource should get from fake controller")
		assert.Equal(t, lhsetting.Value, strconv.Itoa(expected.Storage), "storage not equals")
	})
}
