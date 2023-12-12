//go:build lifecycle || all
// +build lifecycle all

package test

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"fmt"
	"k8gbterratest/utils"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

// TestK8gbRepeatedlyRecreatedFromIngress creates GSLB, then keeps operator live and than recreates GSLB again from Ingress.
// This is usual lifecycle scenario and we are testing spec strategy has expected values.
func TestK8gbRepeatedlyRecreatedFromIngress(t *testing.T) {
	// name of ingress and gslb
	const name = "test-gslb-failover-simple"

	assertStrategy := func(t *testing.T, options *k8s.KubectlOptions) {
		utils.AssertGslbSpec(t, options, name, "spec.strategy.splitBrainThresholdSeconds", "300")
		utils.AssertGslbSpec(t, options, name, "spec.strategy.dnsTtlSeconds", "5")
		utils.AssertGslbSpec(t, options, name, "spec.strategy.primaryGeoTag", settings.PrimaryGeoTag)
		utils.AssertGslbSpec(t, options, name, "spec.strategy.type", "failover")
	}

	// Path to the Kubernetes resource config we will test
	ingressResourcePath, err := filepath.Abs("../examples/ingress-annotation-failover-simple.yaml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("k8gb-test-repeatedly-recreated-from-ingress-%s", strings.ToLower(random.UniqueId()))

	// Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// - Random namespace
	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)

	defer k8s.DeleteNamespace(t, options, namespaceName)

	defer k8s.KubectlDelete(t, options, ingressResourcePath)

	utils.CreateGslb(t, options, settings, ingressResourcePath)

	k8s.WaitUntilIngressAvailable(t, options, name, utils.DefaultRetries, 1*time.Second)

	ingress := k8s.GetIngress(t, options, name)

	require.Equal(t, ingress.Name, name)

	// assert Gslb strategy has expected values
	assertStrategy(t, options)

	k8s.KubectlDelete(t, options, ingressResourcePath)

	utils.AssertGslbDeleted(t, options, ingress.Name)

	// recreate ingress
	utils.CreateGslb(t, options, settings, ingressResourcePath)

	k8s.WaitUntilIngressAvailable(t, options, name, utils.DefaultRetries, 1*time.Second)

	ingress = k8s.GetIngress(t, options, name)

	require.Equal(t, ingress.Name, name)
	// assert Gslb strategy has expected values
	assertStrategy(t, options)
}

// TestK8gbSpecKeepsStableAfterIngressUpdates, If ingress is updated and GSLB has non default values, the GSLB stays
// stable and is not updated.
func TestK8gbSpecKeepsStableAfterIngressUpdates(t *testing.T) {
	t.Parallel()
	// name of ingress and gslb
	const name = "test-gslb-lifecycle"

	assertStrategy := func(t *testing.T, options *k8s.KubectlOptions) {
		utils.AssertGslbSpec(t, options, name, "spec.strategy.splitBrainThresholdSeconds", "600")
		utils.AssertGslbSpec(t, options, name, "spec.strategy.dnsTtlSeconds", "5")
		utils.AssertGslbSpec(t, options, name, "spec.strategy.primaryGeoTag", settings.PrimaryGeoTag)
		utils.AssertGslbSpec(t, options, name, "spec.strategy.type", "failover")
	}

	kubeResourcePath, err := filepath.Abs("../examples/failover-lifecycle.yaml")
	ingressResourcePath, err := filepath.Abs("../examples/ingress-annotation-failover.yaml")
	require.NoError(t, err)
	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("k8gb-test-spec-keeps-stable-after-ingress-updates-%s", strings.ToLower(random.UniqueId()))

	// Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// - Random namespace
	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)
	defer k8s.DeleteNamespace(t, options, namespaceName)

	// create gslb
	utils.CreateGslb(t, options, settings, kubeResourcePath)
	k8s.WaitUntilIngressAvailable(t, options, name, utils.DefaultRetries, 1*time.Second)

	assertStrategy(t, options)

	// reapply ingress
	utils.CreateGslb(t, options, settings, ingressResourcePath)

	k8s.WaitUntilIngressAvailable(t, options, name, utils.DefaultRetries, 1*time.Second)

	ingress := k8s.GetIngress(t, options, name)

	require.Equal(t, ingress.Name, name)
	// assert Gslb strategy has initial values, ingress doesn't change it
	assertStrategy(t, options)
}
