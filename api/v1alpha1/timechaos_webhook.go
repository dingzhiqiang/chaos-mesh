// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var timechaoslog = logf.Log.WithName("timechaos-resource")

// SetupWebhookWithManager setup TimeChaos's webhook with manager
func (in *TimeChaos) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-pingcap-com-v1alpha1-timechaos,mutating=true,failurePolicy=fail,groups=pingcap.com,resources=timechaos,verbs=create;update,versions=v1alpha1,name=mtimechaos.kb.io

var _ webhook.Defaulter = &TimeChaos{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *TimeChaos) Default() {
	timechaoslog.Info("default", "name", in.Name)

	in.Spec.Selector.DefaultNamespace(in.GetNamespace())
	in.Spec.DefaultClockIds()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-pingcap-com-v1alpha1-timechaos,mutating=false,failurePolicy=fail,groups=pingcap.com,resources=timechaos,versions=v1alpha1,name=vtimechaos.kb.io

var _ ChaosValidator = &TimeChaos{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *TimeChaos) ValidateCreate() error {
	timechaoslog.Info("validate create", "name", in.Name)
	return in.Validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *TimeChaos) ValidateUpdate(old runtime.Object) error {
	timechaoslog.Info("validate update", "name", in.Name)
	return in.Validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *TimeChaos) ValidateDelete() error {
	timechaoslog.Info("validate delete", "name", in.Name)

	// Nothing to do?
	return nil
}

// Validate validates chaos object
func (in *TimeChaos) Validate() error {
	specField := field.NewPath("spec")
	allErrs := in.ValidateScheduler(specField)
	allErrs = append(allErrs, in.ValidatePodMode(specField)...)
	allErrs = append(allErrs, in.Spec.validateTimeOffset(specField.Child("timeOffset"))...)

	if len(allErrs) > 0 {
		return fmt.Errorf(allErrs.ToAggregate().Error())
	}
	return nil
}

// ValidateScheduler validates the scheduler and duration
func (in *TimeChaos) ValidateScheduler(spec *field.Path) field.ErrorList {
	return ValidateScheduler(in, spec)
}

// ValidatePodMode validates the value with podmode
func (in *TimeChaos) ValidatePodMode(spec *field.Path) field.ErrorList {
	return ValidatePodMode(in.Spec.Value, in.Spec.Mode, spec.Child("value"))
}

// validateTimeOffset validates the timeOffset
func (in *TimeChaosSpec) validateTimeOffset(timeOffset *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	_, err := time.ParseDuration(in.TimeOffset)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(timeOffset,
			in.TimeOffset,
			fmt.Sprintf("parse timeOffset field error:%s", err)))
	}

	return allErrs
}
