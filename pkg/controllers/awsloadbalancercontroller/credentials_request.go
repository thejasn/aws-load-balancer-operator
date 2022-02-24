package awsloadbalancercontroller

import (
	"context"
	"fmt"
	"reflect"

	cco "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	credentialRequestName = "cluster"
)

// currentCredentialsRequest returns true if credentials request exists.
func (r *AWSLoadBalancerControllerReconciler) currentCredentialsRequest(ctx context.Context, name types.NamespacedName) (bool, *cco.CredentialsRequest, error) {
	cr := &cco.CredentialsRequest{}
	if err := r.Client.Get(ctx, name, cr); err != nil {
		if errors.IsNotFound(err) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, cr, nil
}

func (r *AWSLoadBalancerControllerReconciler) ensureCredentialsRequest(ctx context.Context) error {
	credReq := types.NamespacedName{Name: credentialRequestName, Namespace: r.Namespace}

	reqLogger := log.FromContext(ctx).WithValues("credentialsrequest", credReq)
	reqLogger.Info("reconciling credentials secret for externalDNS instance")

	exists, current, err := r.currentCredentialsRequest(ctx, credReq)
	if err != nil {
		reqLogger.Info("failed to find existing credential request due to %v", err)
		return err
	}

	desired, err := desiredCredentialsRequest(ctx, credReq)
	if err != nil {
		return err
	}

	if !exists {
		if err := r.createCredentialsRequest(ctx, desired); err != nil {
			return err
		}
		_, _, err = r.currentCredentialsRequest(ctx, credReq)
		return err
	}

	if updated, err := r.updateCredentialsRequest(ctx, current, desired); err != nil {
		return err
	} else if updated {
		_, _, err = r.currentCredentialsRequest(ctx, credReq)
		return err
	}

	return nil
}

func (r *AWSLoadBalancerControllerReconciler) createCredentialsRequest(ctx context.Context, desired *cco.CredentialsRequest) error {
	if err := r.Client.Create(ctx, desired); err != nil {
		return fmt.Errorf("failed to create externalDNS credentials request %s: %w", desired.Name, err)
	}
	return nil
}

func (r *AWSLoadBalancerControllerReconciler) updateCredentialsRequest(ctx context.Context, current *cco.CredentialsRequest, desired *cco.CredentialsRequest) (bool, error) {
	var updated *cco.CredentialsRequest
	changed, err := credentialsRequestChanged(current, desired)
	if err != nil {
		return false, err
	}
	if !changed {
		return false, nil
	}
	updated = current.DeepCopy()
	updated.Name = desired.Name
	updated.Namespace = desired.Namespace
	updated.Spec = desired.Spec
	if err := r.Client.Update(ctx, updated); err != nil {
		return false, err
	}
	return true, nil
}

func desiredCredentialsRequest(ctx context.Context, name types.NamespacedName) (*cco.CredentialsRequest, error) {
	credentialsRequest := &cco.CredentialsRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Spec: cco.CredentialsRequestSpec{
			ServiceAccountNames: []string{controllerServiceAccountName},
			SecretRef: corev1.ObjectReference{
				Name:      name.Name,
				Namespace: name.Namespace,
			},
		},
	}

	codec, err := cco.NewCodec()
	if err != nil {
		return nil, err
	}

	providerSpec, err := createProviderConfig(codec)
	if err != nil {
		return nil, err
	}
	credentialsRequest.Spec.ProviderSpec = providerSpec
	return credentialsRequest, nil
}

func createProviderConfig(codec *cco.ProviderCodec) (*runtime.RawExtension, error) {
	return codec.EncodeProviderSpec(&cco.AWSProviderSpec{
		StatementEntries: GetIAMPolicy().Statement,
	})
}

func credentialsRequestChanged(current, desired *cco.CredentialsRequest) (bool, error) {

	if current.Name != desired.Name {
		return true, nil
	}

	if current.Namespace != desired.Namespace {
		return true, nil
	}

	codec, err := cco.NewCodec()
	if err != nil {
		return false, err
	}

	currentAwsSpec := cco.AWSProviderSpec{}
	err = codec.DecodeProviderSpec(current.Spec.ProviderSpec, &currentAwsSpec)
	if err != nil {
		return false, err
	}

	desiredAwsSpec := cco.AWSProviderSpec{}
	err = codec.DecodeProviderSpec(desired.Spec.ProviderSpec, &desiredAwsSpec)
	if err != nil {
		return false, err
	}

	if !(reflect.DeepEqual(currentAwsSpec, desiredAwsSpec)) {
		return true, nil
	}

	return false, nil
}
