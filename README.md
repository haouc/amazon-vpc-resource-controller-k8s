# amazon-vpc-resource-controller-k8s

Controller for managing [Amazon AWS VPC](https://aws.amazon.com/vpc/) resources for [Kubernetes Pods](https://kubernetes.io/docs/concepts/workloads/pods/pod/) 

## Installation

Prerequisite steps

- Create a repository in [ECR](https://console.aws.amazon.com/ecr/home) with the name ```amazon/amazon-k8s-cni```
- Install [cert manager](https://cert-manager.io/docs/installation/kubernetes/) on your cluster

To build the docker image.
```
make docker-build AWS_ACCOUNT=<account-id> AWS_REGION=<region-code>
```
To push the docker image.
```
make docker-push AWS_ACCOUNT=<account-id> AWS_REGION=<region-code>
```

After pushing the docker image, to apply the yaml for the vpc-resource-controller and the vpc-admission-webhook on your cluster.
```
make deploy AWS_ACCOUNT=<account-id> AWS_REGION=<region-code>
```
To build, push and deploy in one command
```
make docker-build docker-push deploy AWS_ACCOUNT=<account-id> AWS_REGION=<region-code>
```
 
## Requirements for managing pod-eni

In order to manage Branch ENI for pods, the account has to be Whitelisted for ENI Trunking/Branching.

## Replacing aws-sdk-go (Internal only)

Workaround to copy the aws-sdk-go to internal package. This only needs to be done if we add a functionality in aws-sdk-go that is not already present in the internal package.

1. Unzip the staging file with private API calls to  ```private_sdk_dir```
2. In go.mod, add ```replace github.com/aws/aws-sdk-go => private_sdk_dir```
3. Run ```go mod vendor``` in project root.
4. Copy vendor directory for aws-sdk-go to internal ```cp -r vendor/github.com/aws/aws-sdk-go internal```
5. Copy go.mod ```cp private_sdk_dir/go.mod internal/aws-sdk-go```
6. Clean up vendor directory ```rm -r vendor```
7. Remove from go.mod ```replace github.com/aws/aws-sdk-go => private_sdk_dir```

This reduces the size of the internal package from ~120 Mbs to ~5 Mbs, this method is error prone but considering it's not going to be repeated again and the APIs will be public soon we will use it for now.
