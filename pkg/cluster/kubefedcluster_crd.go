package cluster

import (
	"context"
	errs "github.com/pkg/errors"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EnsureKubeFedClusterCrd creates a KubeFedCluster CRD in the cluster.
// If the creation returns an error that is of the type "AlreadyExists" then the error is ignored,
// if the error is of another type then it is returned
func EnsureKubeFedClusterCrd(scheme *runtime.Scheme, client client.Client) error {
	decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()
	kubeFedCrd := &v1beta1.CustomResourceDefinition{}
	_, _, err := decoder.Decode([]byte(kubeFedClusterCrd), nil, kubeFedCrd)
	if err != nil {
		return errs.Wrap(err, "unable to decode the KubeFedCluster CRD")
	}
	err = client.Create(context.TODO(), kubeFedCrd)
	if err != nil && !errors.IsAlreadyExists(err) {
		return errs.Wrap(err, "unable to create the KubeFedCluster CRD")
	}
	return nil
}

const kubeFedClusterCrd = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  labels:
    controller-tools.k8s.io: "1.0"
  name: kubefedclusters.core.kubefed.k8s.io
spec:
  additionalPrinterColumns:
  - JSONPath: .status.conditions[?(@.type=='Ready')].status
    name: ready
    type: string
  - JSONPath: .metadata.creationTimestamp
    name: age
    type: date
  group: core.kubefed.k8s.io
  names:
    kind: KubeFedCluster
    plural: kubefedclusters
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            apiEndpoint:
              description: The API endpoint of the member cluster. This can be a hostname,
                hostname:port, IP or IP:port.
              type: string
            caBundle:
              description: CABundle contains the certificate authority information.
              format: byte
              type: string
            secretRef:
              description: Name of the secret containing the token required to access
                the member cluster. The secret needs to exist in the same namespace
                as the control plane and should have a "token" key.
              properties:
                name:
                  description: Name of a secret within the enclosing namespace
                  type: string
              required:
              - name
              type: object
          required:
          - apiEndpoint
          - secretRef
          type: object
        status:
          properties:
            conditions:
              description: Conditions is an array of current cluster conditions.
              items:
                properties:
                  lastProbeTime:
                    description: Last time the condition was checked.
                    format: date-time
                    type: string
                  lastTransitionTime:
                    description: Last time the condition transit from one status to
                      another.
                    format: date-time
                    type: string
                  message:
                    description: Human readable message indicating details about last
                      transition.
                    type: string
                  reason:
                    description: (brief) reason for the condition's last transition.
                    type: string
                  status:
                    description: Status of the condition, one of True, False, Unknown.
                    type: string
                  type:
                    description: Type of cluster condition, Ready or Offline.
                    type: string
                required:
                - type
                - status
                - lastProbeTime
                type: object
              type: array
            region:
              description: Region is the name of the region in which all of the nodes
                in the cluster exist.  e.g. 'us-east1'.
              type: string
            zones:
              description: Zones are the names of availability zones in which the
                nodes of the cluster exist, e.g. 'us-east1-a'.
              items:
                type: string
              type: array
          required:
          - conditions
          type: object
  version: v1beta1
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`
