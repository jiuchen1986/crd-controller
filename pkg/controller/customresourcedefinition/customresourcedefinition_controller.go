package customresourcedefinition

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"

	// apiextensionsv1beta1 "github.com/jiuchen1986/crd-controller/pkg/apis/apiextensions/v1beta1"
	// corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	// meta "k8s.io/apimachinery/pkg/api/meta"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	// uscheme "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructuredscheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	// "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	// "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/yaml"
)

var log = logf.Log.WithName("controller_customresourcedefinition")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CustomResourceDefinition Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCustomResourceDefinition{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("customresourcedefinition-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CustomResourceDefinition
	err = c.Watch(&source.Kind{Type: &apiextensionsv1beta1.CustomResourceDefinition{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCustomResourceDefinition implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCustomResourceDefinition{}

// ReconcileCustomResourceDefinition reconciles a CustomResourceDefinition object
type ReconcileCustomResourceDefinition struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CustomResourceDefinition object and makes changes based on the state read
// and what is in the CustomResourceDefinition.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCustomResourceDefinition) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CustomResourceDefinition")

	// Fetch the CustomResourceDefinition instance
	instance := &apiextensionsv1beta1.CustomResourceDefinition{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logCRDInfoV4(reqLogger, "CRD is not found, maybe deleted!", instance)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	logCRDInfo(reqLogger, "Found a CRD!", instance)

	if crs, err := readCRInstances(instance, r); err != nil {
		logCRDError(reqLogger, err, "Failed to get CR instances of CRD!", instance)
		return reconcile.Result{}, err
	} else {
		logCRDInfoV4(reqLogger, "Found CR instances of CRD!", instance)
		if serialized, err := yaml.Marshal(crs); err != nil {
			logCRDError(reqLogger, err, "Failed to marshal list of CR instances!", instance)
			return reconcile.Result{}, err
		} else {
			fmt.Println("\n" + string(serialized))
		}
	}

	return reconcile.Result{RequeueAfter: time.Second * 3000}, nil
}

func readCRInstances(crd *apiextensionsv1beta1.CustomResourceDefinition, r *ReconcileCustomResourceDefinition) (crs *unstructured.UnstructuredList, err error) {
	gvk := schema.GroupVersionKind{
		Group:   crd.Spec.Group,
		Version: crd.Spec.Version,
		Kind:    crd.Spec.Names.ListKind,
	}
	crs = &unstructured.UnstructuredList{}
	crs.SetGroupVersionKind(gvk)
	err = r.client.List(context.TODO(), &client.ListOptions{Namespace: ""}, crs)
	if err != nil {
		return nil, err
	}
	return
}

func logCRDInfo(logger logr.Logger, msg string, crd *apiextensionsv1beta1.CustomResourceDefinition) {
	logger.Info(msg, "CRD.Group", crd.Spec.Group, "CRD.Version", crd.Spec.Version, "CRD.Kind", crd.Spec.Names.Kind, "CRD.ListKind", crd.Spec.Names.ListKind)
}

func logCRDInfoV4(logger logr.Logger, msg string, crd *apiextensionsv1beta1.CustomResourceDefinition) {
	logger.V(4).Info(msg, "CRD.Group", crd.Spec.Group, "CRD.Version", crd.Spec.Version, "CRD.Kind", crd.Spec.Names.Kind, "CRD.ListKind", crd.Spec.Names.ListKind)
}

func logCRDError(logger logr.Logger, err error, msg string, crd *apiextensionsv1beta1.CustomResourceDefinition) {
	logger.Error(err, msg, "CRD.Group", crd.Spec.Group, "CRD.Version", crd.Spec.Version, "CRD.Kind", crd.Spec.Names.Kind, "CRD.ListKind", crd.Spec.Names.ListKind)
}
