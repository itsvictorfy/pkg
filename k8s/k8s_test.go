package k8s

import (
	"context"
	"fmt"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing" // Add this import
)

func TestCheckClusterConnectivity(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *KubeClient
		wantErr bool
	}{
		{
			name: "successful connectivity check",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset(&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-deployment",
						Namespace: "kube-system",
					},
				})
				return &KubeClient{Clientset: clientset}
			},
			wantErr: false,
		},
		{
			name: "failed connectivity check",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset()
				clientset.PrependReactor("list", "deployments", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("simulated connectivity failure")
				})
				return &KubeClient{Clientset: clientset}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := tt.setup()
			err := kube.CheckClusterConnectivity("test-env")
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckClusterConnectivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTriggerJobFromCronJob(t *testing.T) {
	tests := []struct {
		name    string
		cronJob *batchv1.CronJob
		wantErr bool
	}{
		{
			name: "successful job trigger",
			cronJob: &batchv1.CronJob{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cronjob",
					Namespace: "default",
				},
				Spec: batchv1.CronJobSpec{
					JobTemplate: batchv1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "cronjob not found",
			cronJob: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset()
			if tt.cronJob != nil {
				_, _ = clientset.BatchV1().CronJobs("default").Create(context.TODO(), tt.cronJob, metav1.CreateOptions{})
			}

			kube := &KubeClient{Clientset: clientset}
			_, err := kube.TriggerJobFromCronJob("test-cronjob", "default")
			if (err != nil) != tt.wantErr {
				t.Errorf("TriggerJobFromCronJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateJob(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *KubeClient
		jobName   string
		namespace string
		wantErr   bool
	}{
		{
			name: "successful job creation",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset()
				return &KubeClient{Clientset: clientset}
			},
			jobName:   "test-job",
			namespace: "default",
			wantErr:   false,
		},
		{
			name: "failed job creation",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset()
				clientset.PrependReactor("create", "jobs", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("simulated job creation error")
				})
				return &KubeClient{Clientset: clientset}
			},
			jobName:   "test-job",
			namespace: "default",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := tt.setup()
			err := kube.CreateJob(tt.jobName, tt.namespace, "test-image", nil, nil, nil, nil, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScaleDownDeployment(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *KubeClient
		wantErr bool
	}{
		{
			name: "successful scale down",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset(&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-deployment",
						Namespace: "default",
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: new(int32),
					},
				})
				return &KubeClient{Clientset: clientset}
			},
			wantErr: false,
		},
		{
			name: "deployment not found",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset()
				return &KubeClient{Clientset: clientset}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := tt.setup()
			err := kube.ScaleDownDeployment("test-deployment", "default")
			if (err != nil) != tt.wantErr {
				t.Errorf("ScaleDownDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestartDeployment(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *KubeClient
		wantErr bool
	}{
		{
			name: "successful restart",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset(&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-deployment",
						Namespace: "default",
					},
				})
				return &KubeClient{Clientset: clientset}
			},
			wantErr: false,
		},
		{
			name: "deployment not found",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset()
				return &KubeClient{Clientset: clientset}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := tt.setup()
			err := kube.RestartDeployment("test-deployment", "default")
			if (err != nil) != tt.wantErr {
				t.Errorf("RestartDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestIsJobCompleted(t *testing.T) {
	tests := []struct {
		name      string
		job       *batchv1.Job
		jobName   string
		namespace string
		want      bool
		wantErr   bool
	}{
		{
			name: "job completed successfully",
			job: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-job",
					Namespace: "default",
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
			jobName:   "test-job",
			namespace: "default",
			want:      true,
			wantErr:   false,
		},
		{
			name: "job failed",
			job: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-job",
					Namespace: "default",
				},
				Status: batchv1.JobStatus{
					Failed: 1,
				},
			},
			jobName:   "test-job",
			namespace: "default",
			want:      false,
			wantErr:   true,
		},
		{
			name: "job still in progress",
			job: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-job",
					Namespace: "default",
				},
				Status: batchv1.JobStatus{
					Succeeded: 0,
					Failed:    0,
				},
			},
			jobName:   "test-job",
			namespace: "default",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "job not found",
			job:       nil, // No job created in the fake client
			jobName:   "test-job",
			namespace: "default",
			want:      false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the fake clientset
			clientset := fake.NewSimpleClientset()
			if tt.job != nil {
				_, _ = clientset.BatchV1().Jobs(tt.namespace).Create(context.TODO(), tt.job, metav1.CreateOptions{})
			}

			// Create the KubeClient
			kube := &KubeClient{Clientset: clientset}

			// Call the method under test
			got, err := kube.IsJobCompleted(tt.jobName, tt.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsJobCompleted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsJobCompleted() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestCreateConfigMap(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *KubeClient
		configMap *corev1.ConfigMap
		wantErr   bool
	}{
		{
			name: "successful ConfigMap creation",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset()
				return &KubeClient{Clientset: clientset}
			},
			configMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{"key": "value"},
			},
			wantErr: false,
		},
		{
			name: "ConfigMap with existing name",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-configmap",
						Namespace: "default",
					},
				})
				return &KubeClient{Clientset: clientset}
			},
			configMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{"key": "new-value"},
			},
			wantErr: true,
		},
		{
			name: "clientset failure during ConfigMap creation",
			setup: func() *KubeClient {
				clientset := fake.NewSimpleClientset()
				clientset.PrependReactor("create", "configmaps", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("simulated clientset error")
				})
				return &KubeClient{Clientset: clientset}
			},
			configMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{"key": "value"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kube := tt.setup()
			err := kube.CreateConfigMap(tt.configMap.Name, tt.configMap.Namespace, tt.configMap.Data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateConfigMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
