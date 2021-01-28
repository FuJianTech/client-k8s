package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	ctx,_:=context.WithCancel(context.Background())
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	deploymentName := flag.String("deployment", "", "deployment name")
	imageName := flag.String("image", "", "new image name")
	appName := flag.String("app", "app", "application name")
	namespaceName := flag.String("ns","default","namespace name")
	flag.Parse()


	flag.Parse()
	if *deploymentName == "" {
		fmt.Println("You must specify the deployment name.")
		os.Exit(0)
	}
	if *imageName == "" {
		fmt.Println("You must specify the new image name.")
		os.Exit(0)
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//clientset.AppsV1().DaemonSets("default").Get()
	deployment, err := clientset.AppsV1().Deployments(*namespaceName).Get(ctx,*deploymentName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	if errors.IsNotFound(err) {
		fmt.Printf("Deployment not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting deployment%v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployment\n")
		name := deployment.GetName()
		fmt.Println("name ->", name)
		getAppName := deployment.GetObjectMeta().GetName()
		fmt.Println("getAppName ->", getAppName)
		myAppname := getAppName
		if *appName != "" {
			myAppname = *appName
		}
		containers := &deployment.Spec.Template.Spec.Containers
		found := false
		for i := range *containers {
			c := *containers
			//fmt.Printf(c[0].Image)
			fmt.Printf(c[i].Name)
			if c[i].Name == myAppname {
				found = true
				fmt.Println("Old image ->", c[i].Image)
				fmt.Println("New image ->", *imageName)
				c[i].Image = *imageName
			}
		}
		if found == false {
			fmt.Println("The application container not exist in the deployment pods.")
			os.Exit(0)
		}
		clientset.AppsV1beta1()
		_, err := clientset.AppsV1().Deployments("default").Update(ctx,deployment,metav1.UpdateOptions{

		})
		if err != nil {
			panic(err.Error())
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
