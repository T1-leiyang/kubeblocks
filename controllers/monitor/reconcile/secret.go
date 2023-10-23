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

package reconcile

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitorv1alpha1 "github.com/apecloud/kubeblocks/apis/monitor/v1alpha1"
	"github.com/apecloud/kubeblocks/controllers/monitor/types"
	"github.com/apecloud/kubeblocks/pkg/constant"
)

const OteldSecretName = "oteld-secret"

func Secret(reqCtx types.ReconcileCtx, params types.OTeldParams) error {
	desired := []*corev1.Secret{}

	cg := reqCtx.GetConfigGenerator()
	if cg == nil {
		return fmt.Errorf("config generator is nil")
	}
	daemonsetInstance := reqCtx.GetOteldInstance(monitorv1alpha1.ModeDaemonSet)
	if daemonsetInstance != nil {
		secret, _ := buildSecretForOteld(reqCtx.Config, daemonsetInstance, reqCtx.Namespace, reqCtx.GetExporters(), cg)
		if secret != nil {
			desired = append(desired, secret)
		}
	}

	deploymentInstance := reqCtx.GetOteldInstance(monitorv1alpha1.ModeDeployment)
	if deploymentInstance != nil {
		secret, _ := buildSecretForOteld(reqCtx.Config, deploymentInstance, reqCtx.Namespace, reqCtx.GetExporters(), cg)
		if secret != nil {
			desired = append(desired, secret)
		}
	}

	if err := expectedSecret(reqCtx, params, desired); err != nil {
		return err
	}

	if err := deleteSecret(reqCtx, params, desired); err != nil {
		return err
	}

	return nil
}

func expectedSecret(reqCtx types.ReconcileCtx, params types.OTeldParams, desired []*corev1.Secret) error {
	for _, configmap := range desired {
		desired := configmap

		existing := &corev1.Secret{}
		getErr := params.Client.Get(reqCtx.Ctx, client.ObjectKey{Name: desired.Name, Namespace: desired.Namespace}, existing)
		if getErr != nil && apierrors.IsNotFound(getErr) {
			if createErr := params.Client.Create(reqCtx.Ctx, desired); createErr != nil {
				if apierrors.IsAlreadyExists(createErr) {
					return nil
				}
				return fmt.Errorf("failed to create: %w", createErr)
			}
			reqCtx.Log.V(2).Info("created", "configmap.name", desired.Name, "configmap.namespace", desired.Namespace)
			continue
		} else if getErr != nil {
			return getErr
		}

		updated := existing.DeepCopy()
		if updated.Annotations == nil {
			updated.Annotations = map[string]string{}
		}
		if updated.Labels == nil {
			updated.Labels = map[string]string{}
		}

		updated.Data = desired.Data
		updated.StringData = desired.StringData
		updated.ObjectMeta.OwnerReferences = desired.ObjectMeta.OwnerReferences

		for k, v := range desired.ObjectMeta.Annotations {
			updated.ObjectMeta.Annotations[k] = v
		}
		for k, v := range desired.ObjectMeta.Labels {
			updated.ObjectMeta.Labels[k] = v
		}

		patch := client.MergeFrom(existing)

		if err := params.Client.Patch(reqCtx.Ctx, updated, patch); err != nil {
			return fmt.Errorf("failed to apply changes: %w", err)
		}

		reqCtx.Log.V(2).Info("applied", "configmap.name", desired.Name, "configmap.namespace", desired.Namespace)
	}
	return nil
}

func deleteSecret(reqCtx types.ReconcileCtx, params types.OTeldParams, desired []*corev1.Secret) error {
	listopts := []client.ListOption{
		client.InNamespace(reqCtx.Namespace),
		client.MatchingLabels(map[string]string{
			constant.AppManagedByLabelKey: constant.AppName,
			constant.AppNameLabelKey:      OTeldName,
		}),
	}

	secretList := &corev1.SecretList{}
	if params.Client.List(reqCtx.Ctx, secretList, listopts...) != nil {
		return nil
	}

	for _, configMap := range secretList.Items {
		isdel := true
		for _, keep := range desired {
			if keep.Name == configMap.Name && keep.Namespace == configMap.Namespace {
				isdel = false
				break
			}
		}

		if isdel {
			if err := params.Client.Delete(reqCtx.Ctx, &configMap); err != nil {
				return fmt.Errorf("failed to delete: %w", err)
			}
		}
	}
	return nil
}