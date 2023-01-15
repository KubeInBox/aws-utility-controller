package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ec2OperationType operation that has to be performed on the instance.
// +kubebuilder:validation:Enum=Start;Stopped
type Ec2OperationType string

const (
	Start Ec2OperationType = "Start"
	Stop  Ec2OperationType = "Stop"
)

// Ec2OperationWindowType allows controller to perform operations in the given window.
// +kubebuilder:validation:Enum=OnDemand;Scheduled
type Ec2OperationWindowType string

const (
	// OnDemand indicates that the operation has to be performed right now.
	OnDemand Ec2OperationWindowType = "OnDemand"
	// Scheduled indicates that the operation has to be performed in given scheduled time.
	Scheduled Ec2OperationWindowType = "Scheduled"
)

// Ec2CostOptimizerSpec defines the desired state of Ec2CostOptimizer
type Ec2CostOptimizerSpec struct {
	// StopInstanceID on which start/stop operations has to be performed
	// +kubebuilder:validation:Items:MinLength=1
	InstanceIDs []string `json:"instance_ids"`
	// START/STOP operation
	Operation Ec2OperationType `json:"operation"`
	// OnDemand/Scheduled window
	WindowType Ec2OperationWindowType `json:"window_type"`
	// Scheduled start time window, should be valid  start time, supported timezone is IST
	StartTimeWindow string `json:"start_time_window,omitempty"`
	// Scheduled end time window, should be valid  end time, supported timezone is IST
	EndTimeWindow string `json:"end_time_window,omitempty"`
}

// Ec2CostOptimizerStatus defines the observed state of Ec2CostOptimizer
type Ec2CostOptimizerStatus struct {
	// InstanceID is unique identifier for aws-ec2 instance.
	InstanceID string `json:"instance_id,omitempty"`
	// Status represents current state of machine, RUNNING, STOPPED.
	State string `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Ec2CostOptimizer is the Schema for the ec2costoptimizers API
type Ec2CostOptimizer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Ec2CostOptimizerSpec   `json:"spec,omitempty"`
	Status Ec2CostOptimizerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// Ec2CostOptimizerList contains a list of Ec2CostOptimizer
type Ec2CostOptimizerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ec2CostOptimizer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ec2CostOptimizer{}, &Ec2CostOptimizerList{})
}
