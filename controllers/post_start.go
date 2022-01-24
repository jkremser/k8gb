package controllers

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
	"context"
	"fmt"

	"github.com/k8gb-io/k8gb/controllers/depresolver"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *GslbReconciler) PostStartHook(operatorConfig *depresolver.Config, mgr ctrl.Manager) error {
	if err := inferClusterGeoTag(operatorConfig, mgr.GetAPIReader()); err != nil {
		log.Err(err).Msgf("Can't infer the %s", depresolver.ClusterGeoTagKey)
		return err
	}
	if err := createOrUpdateExternalDNSConfigMap(operatorConfig, mgr.GetClient()); err != nil {
		log.Err(err).Msg("Can't create/update config map for external-dns")
		return err
	}
	return nil
}

func inferClusterGeoTag(cfg *depresolver.Config, client client.Reader) error {
	if len(cfg.ClusterGeoTag) != 0 {
		// env var has the highest precedence
		return nil
	}
	nodeList := &corev1.NodeList{}
	err := client.List(context.TODO(), nodeList)
	if err != nil {
		return fmt.Errorf("unable to get the nodes in the cluster, error: %v", err)
	}
	if len(nodeList.Items) == 0 {
		return fmt.Errorf("no available nodes in the cluster")
	}
	// values of this annotation/label is cloud provider specific
	const region = "topology.kubernetes.io/region" // (example: africa-east-1)
	// assuming all the nodes to have the same topology.kubernetes.io/region so
	var foundTag string
	for _, node := range nodeList.Items {
		switch {
		case len(node.Annotations[region]) != 0:
			foundTag = node.Annotations[region]
		case len(node.Labels[region]) != 0:
			foundTag = node.Labels[region]
		}
		if len(cfg.ClusterGeoTag) != 0 && cfg.ClusterGeoTag != foundTag {
			return fmt.Errorf("%v nodes were tried, but they don't have the same annotation/label on "+
				"them ('%s' != '%s'). Remedy: make sure the value of '%s' is same on each node in the cluster."+
				"Details: '%s' != '%s', where the latter was found on node '%s'",
				len(nodeList.Items), cfg.ClusterGeoTag, foundTag, region, cfg.ClusterGeoTag, foundTag, node.GetName())
		}
		cfg.ClusterGeoTag = foundTag
	}
	if len(cfg.ClusterGeoTag) == 0 {
		return fmt.Errorf("%v nodes were tried, but none of them were annotated. Either set the %s explicitly"+
			" as operator conviguraion using env variable, or mark the nodes with label (or annotation) '%s'",
			len(nodeList.Items), depresolver.ClusterGeoTagKey, region)
	}
	return nil
}

func createOrUpdateExternalDNSConfigMap(cfg *depresolver.Config, c client.Client) error {
	const (
		ns             = "k8gb"
		cmName         = "external-dns-env"
		txtOwnerKey    = "EXTERNAL_DNS_TXT_OWNER_ID"
		txtOwnerPrefix = "k8gb"
	)
	txtOwnerValue := fmt.Sprintf("%s-%s-%s", txtOwnerPrefix, cfg.DNSZone, cfg.ClusterGeoTag)

	cm := corev1.ConfigMap{
		Data: map[string]string{
			txtOwnerKey: txtOwnerValue,
		},
	}
	cm.SetName(cmName)
	cm.SetNamespace(ns)
	if err := c.Create(context.TODO(), &cm); err != nil {
		if errors.IsAlreadyExists(err) {
			if updateErr := c.Update(context.TODO(), &cm); updateErr != nil {
				return updateErr
			}
			log.Info().Msg("Configmap for external-dns has been updated")
		}
		return err
	}
	log.Info().Msg("Configmap for external-dns has been created")
	return nil
}
