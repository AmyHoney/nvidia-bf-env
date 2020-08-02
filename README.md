# vSphere for Kubernetes Integration with Bitfusion (VMware NGC Bootcamp Project 2020)

[TOC]

### Background

With the release of key new capabilities in vSphere 7, the tools of vSphere include Kubernetes, Lifecycle Manager, Intrinsic Security and Application Acceleration. By utilizing these tools, the operational efficiency and collaboration can be improved, the automated configuration failures and security risks can be reduced, and realize higher automation, accelerate development and innovation. In the NGC Bootcamp 2020, our team's goals are as follows: 

- Getting to know the features and architectures of vSphere for Kubernetes, Bitfusion, and Govmomi
- Programming integratedly with Docker and Kubernetes under vSphere
- Implementing automated processes to setup some applications and develop the environment with vSphere for Kubernetes and Bitfusion



### Tasks

- Auto-deploy bitfusion client and run a tensorflow example to test the liberary：[Bitfusion_Client_Depolyment_bash](https://gitlab.eng.vmware.com/liuqi/ngc-bootcamp-2020-bj-cpbu/tree/master/Bitfusion_Client_Deployment_bash)
- Auto-develop bitfusion client dockerfile
- Auto-clone a vm with govmomi
- Kubernetes : [ngc_k8s_operator_bifusion_depolyment](https://gitlab.eng.vmware.com/liuqi/ngc-bootcamp-2020-bj-cpbu/tree/master/ngc_k8s_operator_bitfusion_deployment)



### How to RUN

#### 1、Bitfusion client ENV in vm setup

Ubuntu 16.04

```console
$ sudo bash ./Bitfusion_Client_Deployment_bash/nvidia-env-with-ubuntu16.04.sh
```

RHEL7.4 & CentOS7 

```console
$ sudo bash Bitfusion_Client_Deployment_bash/nvidia-env-with-rhel-and-CentOS.sh
```

#### 2、Bitfusion Docker ENV setup

**Docker Image Build**：Copy the Dockerfile into individual directory

Ubuntu 16.04

```console
$ docker build -t ubuntu16.04_nvidia:v1 . --network=host
```

RHEL7.4 & CentOS7

```console
$ docker build -t centos_rhel_nvidia:v1 . --network=host
```

**Docker Run**

Ubuntu 16.04

```console
$ docker run -itd --net=host ubuntu16.04_nvidia:v1 /bin/bash
```

URHEL7.4 & CentOS7

```console
$ docker run -itd --net=host centos_rhel_nvidia:v1 /bin/bash
```

#### 3、Clone a VM with Govmomi

Firstly, install the govmomi library for interacting with VMware vSphere APIs

```console
$ go get -u github.com/vmware/govmomi
```

Secondly, install  the govc, a vSphere CLI built on top of govmomi (make sure to set the environment variable [GOPATH](https://github.com/golang/go/wiki/SettingGOPATH))

```console
$ go install github.com/vmware/govmomi/govc
```

Then, clone a vm

```console
$ govc vm.clone -u https://administrator@vsphere.local:Admin!23@10.110.165.188 -k=true -ds hdd01-180 -pool 10.110.166.180/Resources  -vm photon-001 govmomi-clone-02
```

#### 4、Auto deploy Bitfusion client

Use Kubernetes Operator and [govmomi](https://github.com/vmware/govmomi) with vSphere.

Setup CRD

```console
$ cd ./ngc_k8s_operator_bitfusion_deployment && make && make install && make run
```

Trigger Controller to auto deploy

```console
$ kubectl apply -f config/samples/bit_v1_vmbit.yaml
```

