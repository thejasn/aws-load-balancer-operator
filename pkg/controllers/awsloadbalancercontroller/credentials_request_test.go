package awsloadbalancercontroller

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/openshift/aws-load-balancer-operator/pkg/controllers/utils/test"
	cco "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	testOperatorNamespace           = test.OperatorNamespace
	testOperandNamespace            = test.OperandNamespace
	testCredentialsRequestName      = credentialRequestName
	testCredentialsRequestNamespace = credentialRequestNamespace
)

func TestEnsureCredentialsRequest(t *testing.T) {

	managedTypesList := []client.ObjectList{
		&cco.CredentialsRequestList{},
	}

	eventWaitTimeout := time.Duration(1 * time.Second)

	testcases := []struct {
		name            string
		existingObjects []runtime.Object
		expectedEvents  []test.Event
		errExpected     bool
	}{
		{
			name:            "Initial bootstrap",
			existingObjects: make([]runtime.Object, 0),
			expectedEvents: []test.Event{
				{
					EventType: watch.Added,
					ObjType:   "credentialsrequest",
					NamespacedName: types.NamespacedName{
						Namespace: testCredentialsRequestNamespace,
						Name:      testCredentialsRequestName,
					},
				},
			},
			errExpected: false,
		},
		{
			name: "Change in Credential Request",
			existingObjects: []runtime.Object{
				testCredentialsRequest(),
			},
			expectedEvents: []test.Event{
				{
					EventType: watch.Modified,
					ObjType:   "credentialsrequest",
					NamespacedName: types.NamespacedName{
						Namespace: testCredentialsRequestNamespace,
						Name:      testCredentialsRequestName,
					},
				},
			},
			errExpected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(test.Scheme).WithRuntimeObjects(tc.existingObjects...).Build()

			r := &AWSLoadBalancerControllerReconciler{
				Client:    cl,
				Namespace: test.OperandNamespace,
				Image:     test.OperandImage,
			}

			c := test.NewEventCollector(t, cl, managedTypesList, len(tc.expectedEvents))

			// get watch interfaces from all the type managed by the operator
			c.Start(context.TODO())
			defer c.Stop()

			err := r.ensureCredentialsRequest(context.TODO())
			// error check
			if err != nil {
				if !tc.errExpected {
					t.Fatalf("got unexpected error: %v", err)
				}
			} else if tc.errExpected {
				t.Fatalf("error expected but not received")
			}

			// collect the events received from Reconcile()
			collectedEvents := c.Collect(len(tc.expectedEvents), eventWaitTimeout)

			// compare collected and expected events
			idxExpectedEvents := test.IndexEvents(tc.expectedEvents)
			idxCollectedEvents := test.IndexEvents(collectedEvents)
			if diff := cmp.Diff(idxExpectedEvents, idxCollectedEvents); diff != "" {
				t.Fatalf("found diff between expected and collected events: %s", diff)
			}
		})
	}
}

func testCredentialsRequest() *cco.CredentialsRequest {
	return &cco.CredentialsRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testCredentialsRequestName,
			Namespace: testCredentialsRequestNamespace,
		},
		Spec: cco.CredentialsRequestSpec{
			ProviderSpec: testAWSProviderSpec(),
		},
	}
}

func testAWSProviderSpec() *runtime.RawExtension {
	codec, _ := cco.NewCodec()
	providerSpec, _ := codec.EncodeProviderSpec(&cco.AWSProviderSpec{})
	return providerSpec
}
