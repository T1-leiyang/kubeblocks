/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package types

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/apis/monitor/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/configuration/core"
	intctrlutil "github.com/apecloud/kubeblocks/pkg/controllerutil"
)

type ScrapeConfig map[string]any

func fromCollectorDataSource(dataSourceSpec v1alpha1.CollectorDataSourceSpec, cli client.Client, ctx context.Context, namespace string) ([]ScrapeConfig, error) {
	configs := make([]ScrapeConfig, 0)
	clusterName := dataSourceSpec.ClusterRef

	if cli == nil || ctx == nil {
		return nil, core.MakeError("client or context is nil")
	}

	resource := intctrlutil.NewResourceFetcher(&intctrlutil.ResourceCtx{
		Context:     ctx,
		Client:      cli,
		ClusterName: clusterName,
		Namespace:   namespace,
	})

	if err := resource.Cluster().ClusterDef().Complete(); err != nil {
		return nil, err
	}
	clusterDefName := resource.ClusterObj.Spec.ClusterDefRef
	for _, spec := range dataSourceSpec.CollectorSpecs {
		resource.ComponentName = spec.ComponentName
		if err := resource.ClusterComponent().ClusterDefComponent().Complete(); err != nil {
			return nil, err
		}
		if resource.ClusterDefComObj == nil {
			return nil, core.MakeError("failed to found componentDef[%s] in the clusterDefinition[%s]", resource.ComponentName, clusterDefName)
		}
		for _, config := range spec.ScrapeConfigs {
			configs = append(configs, buildEngineValMap(clusterName,
				resource.ComponentName,
				config,
				resource.ClusterDefComObj,
				clusterDefName),
			)
		}
	}
	return configs, nil
}

func buildEngineValMap(clusterName string, componentName string, config v1alpha1.ScrapeConfig, obj *appsv1alpha1.ClusterComponentDefinition, clusterDefName string) ScrapeConfig {
	monitorType := obj.CharacterType
	if config.Metrics != nil && config.Metrics.MonitorType != "" {
		monitorType = config.Metrics.MonitorType
	}

	valMap := map[string]any{
		"cluster_name":    clusterName,
		"component_name":  componentName,
		"container_name":  config.ContainerName,
		"collector_name":  monitorType,
		"metrics.enabled": config.Metrics != nil,
		"logs.enabled":    config.Logs != nil,
		// for clusterDefinition
		"cluster_def_name":   clusterDefName,
		"component_def_name": obj.Name,
	}

	if config.Metrics != nil {
		if len(config.Metrics.MetricsSelector) > 0 {
			valMap["metrics.enabled_metrics"] = config.Metrics.MetricsSelector
		}
		valMap["metrics.collection_interval"] = config.Metrics.CollectionInterval
	}
	if config.Logs == nil {
		return valMap
	}

	logCollector := map[string]v1alpha1.InputConfig{}
	isAll := isAllLogs(config.Logs.LogTypes)
	for _, logConfig := range obj.LogConfigs {
		if isAll || isMatchLogs(logConfig, config.Logs.LogTypes) {
			logCollector[logConfig.Name] = v1alpha1.InputConfig{
				Include: []string{logConfig.FilePathPattern},
			}
		}
	}
	if len(logCollector) > 0 {
		valMap["logs.logs_collector"] = logCollector
	}
	return valMap
}

func isMatchLogs(config appsv1alpha1.LogConfig, types []string) bool {
	for _, s := range types {
		if config.Name == s {
			return true
		}
	}
	return false
}

func isAllLogs(types []string) bool {
	return len(types) == 1 && types[0] == "*"
}