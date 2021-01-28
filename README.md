# Kubernetes client-go sample
This is a simple code sample for using kubernetes [client-go](https://github.com/kubernetes/client-go) package.

You can use this code to update a specify deployment's application image (More than one container in a pod).

## Build

```bash
go build main.go
```

## Usage

```bash
Usage of ./main:
  -app string
    	application name (默认是deployment的名字)
  -deployment string
    	deployment name
  -image string
    	new image name
  -kubeconfig string
    	(optional) absolute path to the kubeconfig file (default "/Users/jimmy/.kube/config")
  -ns string
        命名空间(默认命名空间是 default)
```

- `-image`: 需要更新的镜像名字
- `-deployment`: deployment name
- `-app`: 应用的app名字
- `-kubeconfig`: k8s应用文件 (default "$HOME/.kube/config")
- '-ns': 命名空间名字 默认是 default

## 例子

```bash
./update-deployment-image -image test:Build_8 -deployment filebeat-test
Found deployment
name -> filebeat-test
Old image -> test:Build_7
New image -> test:Build_8
```



