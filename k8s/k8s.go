package k8s

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	*kubernetes.Clientset `json:"-"` // Exclude from JSON
}

// CheckClusterConnectivity checks the connectivity to the cluster
func (k *KubeClient) CheckClusterConnectivity(env string) error { //V
	_, err := k.AppsV1().Deployments("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("connectivity check failed: %v", err)
	}
	return nil
}

// Initialze Kubernetes Client per Cluster
func (k *KubeClient) InitClient(env string) error { //V
	var kubeconfigPath string
	if runtime.GOOS == "windows" {
		kubeconfigPath = filepath.Join(os.Getenv("USERPROFILE"), ".kube", "config")
	} else {
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig file: %v", err)
	}
	k.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("unable to create %s client from config: %v", env, err)
	}
	err = k.CheckClusterConnectivity(env)
	if err != nil {
		return fmt.Errorf("connection to %s cluster failed: %v", env, err)
	}
	return nil
}

// Creates a Job from a CronJob
func (k *KubeClient) TriggerJobFromCronJob(cronjobName, namespace string) (string, error) {
	jobName := fmt.Sprintf("%s-api-trigger-%s", cronjobName, time.Now().Format("20060102150405"))
	cronjob, err := k.BatchV1().CronJobs(namespace).Get(context.TODO(), cronjobName, metav1.GetOptions{})
	if err != nil {
		return cronjobName, fmt.Errorf("kube: unable to retrieve cronjob %v", err)
	}
	podTemplateSpec := cronjob.Spec.JobTemplate.Spec.Template
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: podTemplateSpec,
		},
	}
	_, err = k.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return cronjobName, fmt.Errorf("kube: unable to create job %v", err)
	}
	return jobName, nil
}
func (k *KubeClient) CreateJob(jobName, namespace, image string, commands []string, envVars, labels map[string]string, volumeNameP, mountPathP *string) error {
	var env []corev1.EnvVar
	for key, value := range envVars {
		env = append(env, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	if volumeNameP != nil && mountPathP != nil {
		volume := corev1.Volume{
			Name: *volumeNameP,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: *volumeNameP,
				},
			},
		}

		volumeMount := corev1.VolumeMount{
			Name:      *volumeNameP,
			MountPath: *mountPathP,
		}

		volumes = append(volumes, volume)
		volumeMounts = append(volumeMounts, volumeMount)
	}

	// Define the Job
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   jobName,
			Labels: labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:         jobName,
							Image:        image,
							Command:      commands,
							Env:          env,
							VolumeMounts: volumeMounts, // Add only if volumeMounts are defined
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes:       volumes, // Add only if volumes are defined
				},
			},
		},
	}

	// Create the Job in the specified namespace
	_, err := k.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("kube: unable to create job %v", err)
	}

	return nil
}

// ScaleDownDeployment scales down the deployment to 0 replicas
func (k *KubeClient) ScaleDownDeployment(name, namespace string) error {
	deployment, err := k.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("kube: unable to retrieve deployment %v", err)
	}
	deployment.Spec.Replicas = new(int32)
	*deployment.Spec.Replicas = 0
	_, err = k.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("kube: unable to scale down deployment %v", err)
	}
	return nil
}

// ScaleUpDeployment scales up the deployment to the specified number of replicas
func (k *KubeClient) ScaleUpDeployment(name, namespace string, replicas int32) error {
	deployment, err := k.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("kube: unable to retrieve deployment %v", err)
	}
	deployment.Spec.Replicas = new(int32)
	*deployment.Spec.Replicas = replicas
	_, err = k.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("kube: unable to scale up deployment %v", err)
	}
	return nil
}

// ScaleDownDeploymentsInNamespace scales down all deployments in the namespace to 0 replicas
func (k *KubeClient) ScaleDownAllDeploymentsInNamespace(namespace string) error {
	deployments, err := k.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error getting %v deployments: %v", namespace, err)
	}
	replicas := int32(0)
	for _, d := range deployments.Items {
		d.Spec.Replicas = &replicas
		_, err := k.AppsV1().Deployments(namespace).Update(context.TODO(), &d, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("error scaling down %s: %v", d.Name, err)
		}
	}
	return nil
}

// RestartDeployment restarts the deployment by updating the annotations
func (k *KubeClient) RestartDeployment(deploymentName, namespace string) error {
	timestamp := time.Now().Format(time.RFC3339)
	patchData := fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"kubectl.kubernetes.io/restartedAt": "%s"}}}}}`, timestamp)
	_, err := k.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch deployment: %v", err)
	}
	return nil
}

func (k *KubeClient) CreateConfigMap(name, namespace string, data map[string]string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	_, err := k.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("kube: unable to create configmap %v", err)
	}
	return nil
}

func (k *KubeClient) IsJobCompleted(jobName, namespace string) (bool, error) {
	job, err := k.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("kube: unable to retrieve job %v", err)
	}
	if job.Status.Succeeded > 0 {
		return true, nil
	}

	if job.Status.Failed > 0 {
		return false, fmt.Errorf("job failed")
	}
	return false, nil
}
