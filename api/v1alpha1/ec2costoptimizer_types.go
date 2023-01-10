package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ec2CostOptimizerSpec defines the desired state of Ec2CostOptimizer
type Ec2CostOptimizerSpec struct {
	// StopInstanceID on which start/stop operations has to be performed
	InstanceID string `json:"instance_id"`
	// START/STOP operation
	Operation string `json:"operation"`
	// OnDemand/Scheduled window
	WindowType string `json:"window_type"`
	// Scheduled window time frame
	ScheduleWindow string `json:"schedule_window,omitempty"`
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
