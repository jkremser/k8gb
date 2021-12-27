package main

/*
Copyright 2021 The k8gb Contributors.

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
	"context"
	"fmt"
	"os"

	str "github.com/AbsaOSS/gopkg/strings"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	"github.com/k8gb-io/k8gb/controllers"
	"github.com/k8gb-io/k8gb/controllers/depresolver"
	"github.com/k8gb-io/k8gb/controllers/logging"
	"github.com/k8gb-io/k8gb/controllers/providers/dns"
	"github.com/k8gb-io/k8gb/controllers/providers/metrics"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	externaldns "sigs.k8s.io/external-dns/endpoint"
	// +kubebuilder:scaffold:imports
)

var (
	runtimescheme = runtime.NewScheme()
	version       = "development"
	commit        = "none"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(runtimescheme))

	utilruntime.Must(k8gbv1beta1.AddToScheme(runtimescheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var f *dns.ProviderFactory
	resolver := depresolver.NewDependencyResolver()
	config, err := resolver.ResolveOperatorConfig()
	deprecations := resolver.GetDeprecations()
	// Initialize desired log or default log in case of configuration failed.
	logging.Init(config)
	log := logging.Logger()
	log.Info().
		Str("version", version).
		Str("commit", commit).
		Msg("K8gb status")
	if err != nil {
		log.Err(err).Msg("can't resolve environment variables")
		return err
	}
	log.Debug().
		Str("config", str.ToString(config)).
		Msg("Resolved config")

	ctrl.SetLogger(logging.NewLogrAdapter(log))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             runtimescheme,
		MetricsBindAddress: config.MetricsAddress,
		Port:               9443,
		LeaderElection:     false,
		LeaderElectionID:   "8020e9ff.absa.oss",
	})
	if err != nil {
		log.Err(err).Msg("Unable to start k8gb")
		return err
	}

	err = inferClusterGeoTag(config, mgr.GetClient())
	if err != nil {
		log.Err(err).Msg("Can't infer the CLUSTER_GEO_TAG")
		return err
	}

	for _, d := range deprecations {
		log.Warn().Msg(d)
	}

	log.Info().Msg("Registering components")

	// Add external-dns DNSEndpoints resource
	// https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#adding-3rd-party-resources-to-your-operator
	schemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
	schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
	if err := schemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
		log.Err(err).Msg("Extending scheme")
		return err
	}

	reconciler := &controllers.GslbReconciler{
		Config:      config,
		Client:      mgr.GetClient(),
		DepResolver: resolver,
		Scheme:      mgr.GetScheme(),
	}

	log.Info().Msg("Starting metrics")
	metrics.Init(config)
	defer metrics.Metrics().Unregister()
	err = metrics.Metrics().Register()
	if err != nil {
		log.Err(err).Msg("Register metrics error")
		return err
	}

	log.Info().Msg("Resolving DNS provider")
	f, err = dns.NewDNSProviderFactory(reconciler.Client, *reconciler.Config)
	if err != nil {
		log.Err(err).Msg("Unable to create factory")
		return err
	}
	reconciler.DNSProvider = f.Provider()
	log.Info().Str("provider", reconciler.DNSProvider.String()).Msg("Started")

	if err = reconciler.SetupWithManager(mgr); err != nil {
		log.Err(err).Msg("Unable to create controller Gslb")
		return err
	}
	metrics.Metrics().SetRuntimeInfo(version, commit)
	// +kubebuilder:scaffold:builder
	log.Info().Msg("Starting k8gb")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Err(err).Msg("Problem running k8gb controller")
		return err
	}
	log.Info().Msg("Gracefully finished, bye!\n")
	return nil
}

func inferClusterGeoTag(operatorConfig *depresolver.Config, client client.Client) error {
	if len(operatorConfig.ClusterGeoTag) != 0 {
		return nil
	}
	nodeList := &corev1.NodeList{}
	err := client.List(context.TODO(), nodeList)
	if err == nil {
		return fmt.Errorf("unable to get the nodes in the cluster (RBAC?)")
	}
	if len(nodeList.Items) == 0 {
		return fmt.Errorf("no available nodes in the cluster")
	}
	// values of these annotations/labels are cloud provider specific
	const region = "topology.kubernetes.io/region" // top-lvl (example: africa-east-1)
	const zone = "topology.kubernetes.io/zone"     // lower-lvl   (example: africa-east-1a)
	// assuming all the nodes to have the same topology.kubernetes.io/region and topology.kubernetes.io/zone so
	for _, node := range nodeList.Items {
		switch {
		case len(node.Annotations[zone]) != 0:
			operatorConfig.ClusterGeoTag = node.Annotations[zone]
		case len(node.Labels[zone]) != 0:
			operatorConfig.ClusterGeoTag = node.Labels[zone]
		case len(node.Annotations[region]) != 0:
			operatorConfig.ClusterGeoTag = node.Annotations[region]
		case len(node.Labels[region]) != 0:
			operatorConfig.ClusterGeoTag = node.Labels[region]
		}
		if len(operatorConfig.ClusterGeoTag) != 0 {
			return nil
		}
	}
	if len(operatorConfig.ClusterGeoTag) == 0 {
		return fmt.Errorf("%v nodes were tried, but none of them were annotated."+
			" Either set the CLUSTER_GEO_TAG explicitly as operator conviguraion using env variable,"+
			" or mark the nodes with label (or annotation) '%s' or '%s'", len(nodeList.Items), zone, region)
	}
	return nil
}
