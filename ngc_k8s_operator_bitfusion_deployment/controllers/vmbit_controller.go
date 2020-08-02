/*
Copyright 2020 ncg.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	// "fmt"
	// "os/exec"
	// "time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bitv1 "ncg/api/v1"

	_ "github.com/vmware/govmomi/govc/about"
	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/cluster"
	_ "github.com/vmware/govmomi/govc/cluster/group"
	_ "github.com/vmware/govmomi/govc/cluster/override"
	_ "github.com/vmware/govmomi/govc/cluster/rule"
	_ "github.com/vmware/govmomi/govc/datacenter"
	_ "github.com/vmware/govmomi/govc/datastore"
	_ "github.com/vmware/govmomi/govc/datastore/cluster"
	_ "github.com/vmware/govmomi/govc/datastore/disk"
	_ "github.com/vmware/govmomi/govc/datastore/maintenance"
	_ "github.com/vmware/govmomi/govc/datastore/vsan"
	_ "github.com/vmware/govmomi/govc/device"
	_ "github.com/vmware/govmomi/govc/device/cdrom"
	_ "github.com/vmware/govmomi/govc/device/floppy"
	_ "github.com/vmware/govmomi/govc/device/scsi"
	_ "github.com/vmware/govmomi/govc/device/serial"
	_ "github.com/vmware/govmomi/govc/device/usb"
	_ "github.com/vmware/govmomi/govc/disk"
	_ "github.com/vmware/govmomi/govc/disk/snapshot"
	_ "github.com/vmware/govmomi/govc/dvs"
	_ "github.com/vmware/govmomi/govc/dvs/portgroup"
	_ "github.com/vmware/govmomi/govc/env"
	_ "github.com/vmware/govmomi/govc/events"
	_ "github.com/vmware/govmomi/govc/export"
	_ "github.com/vmware/govmomi/govc/extension"
	_ "github.com/vmware/govmomi/govc/fields"
	_ "github.com/vmware/govmomi/govc/folder"
	_ "github.com/vmware/govmomi/govc/host"
	_ "github.com/vmware/govmomi/govc/host/account"
	_ "github.com/vmware/govmomi/govc/host/autostart"
	_ "github.com/vmware/govmomi/govc/host/cert"
	_ "github.com/vmware/govmomi/govc/host/date"
	_ "github.com/vmware/govmomi/govc/host/esxcli"
	_ "github.com/vmware/govmomi/govc/host/firewall"
	_ "github.com/vmware/govmomi/govc/host/maintenance"
	_ "github.com/vmware/govmomi/govc/host/option"
	_ "github.com/vmware/govmomi/govc/host/portgroup"
	_ "github.com/vmware/govmomi/govc/host/service"
	_ "github.com/vmware/govmomi/govc/host/storage"
	_ "github.com/vmware/govmomi/govc/host/vnic"
	_ "github.com/vmware/govmomi/govc/host/vswitch"
	_ "github.com/vmware/govmomi/govc/importx"
	_ "github.com/vmware/govmomi/govc/library"
	_ "github.com/vmware/govmomi/govc/library/session"
	_ "github.com/vmware/govmomi/govc/library/subscriber"
	_ "github.com/vmware/govmomi/govc/license"
	_ "github.com/vmware/govmomi/govc/logs"
	_ "github.com/vmware/govmomi/govc/ls"
	_ "github.com/vmware/govmomi/govc/metric"
	_ "github.com/vmware/govmomi/govc/metric/interval"
	_ "github.com/vmware/govmomi/govc/namespace/cluster"
	_ "github.com/vmware/govmomi/govc/object"
	_ "github.com/vmware/govmomi/govc/option"
	_ "github.com/vmware/govmomi/govc/permissions"
	_ "github.com/vmware/govmomi/govc/pool"
	_ "github.com/vmware/govmomi/govc/role"
	_ "github.com/vmware/govmomi/govc/session"
	_ "github.com/vmware/govmomi/govc/sso/group"
	_ "github.com/vmware/govmomi/govc/sso/service"
	_ "github.com/vmware/govmomi/govc/sso/user"
	_ "github.com/vmware/govmomi/govc/tags"
	_ "github.com/vmware/govmomi/govc/tags/association"
	_ "github.com/vmware/govmomi/govc/tags/category"
	_ "github.com/vmware/govmomi/govc/task"
	_ "github.com/vmware/govmomi/govc/vapp"
	_ "github.com/vmware/govmomi/govc/vcsa/log"
	_ "github.com/vmware/govmomi/govc/version"
	_ "github.com/vmware/govmomi/govc/vm"
	_ "github.com/vmware/govmomi/govc/vm/disk"
	_ "github.com/vmware/govmomi/govc/vm/guest"
	_ "github.com/vmware/govmomi/govc/vm/network"
	_ "github.com/vmware/govmomi/govc/vm/option"
	_ "github.com/vmware/govmomi/govc/vm/rdm"
	_ "github.com/vmware/govmomi/govc/vm/snapshot"
)

// VmbitReconciler reconciles a Vmbit object
type VmbitReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=bit.ncg.com,resources=vmbits,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bit.ncg.com,resources=vmbits/status,verbs=get;update;patch

func (r *VmbitReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("vmbit", req.NamespacedName)

	// your logic here

	// cmd := exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
	// fmt.Println("NetWorkStatus Start:", time.Now().Unix())
	// err := cmd.Run()
	// fmt.Println("NetWorkStatus End  :", time.Now().Unix())
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println("Net Status , OK")
	// }

	// cmd := exec.Command("govc", "vm.clone", "-u", "https://administrator@vsphere.local:Admin!23@10.110.165.188", "-k=true", "-ds", "hdd01-180", "-pool", "10.110.166.180/Resources", "-vm", "photon-001", "govmomi-clone-01")

	// if err := cmd.Run(); err != nil {
	// 	r.Log.Info("vm.clone Success!")
	// } else {
	// 	r.Log.Info("vm.clone Failure!")
	// }
	var source_vm = "bitfusion-client-u1604"
	// "govmomi-clone-bitfusion-client_rhel01"
	var dest_vm = "govmomi-clone-bitfusion-client01"
	var gvmomiArgs = []string{"vm.clone", "-u", "https://administrator@vsphere.local:Admin!23@10.110.165.188", "-k=true", "-ds", "hdd01-180", "-pool", "10.110.166.180/Resources", "-vm", source_vm, dest_vm}

	// gvmomiArgs = append(gvmomiArgs, "vm.clone")
	// gvmomiArgs = append(gvmomiArgs, "-u")
	// gvmomiArgs = append(gvmomiArgs, "https://administrator@vsphere.local:Admin!23@10.110.165.188")
	// gvmomiArgs = append(gvmomiArgs, "-k=true")
	// gvmomiArgs = append(gvmomiArgs, "-ds")
	// gvmomiArgs = append(gvmomiArgs, "hdd01-180")
	// gvmomiArgs = append(gvmomiArgs, "-pool")
	// gvmomiArgs = append(gvmomiArgs, "10.110.166.180/Resources")
	// gvmomiArgs = append(gvmomiArgs, "-vm")
	// gvmomiArgs = append(gvmomiArgs, "photon-001")
	// gvmomiArgs = append(gvmomiArgs, "govmomi-clone-01")

	res := cli.Run(gvmomiArgs)

	// fmt.Println(gvmomiArgs)
	// fmt.Println(res)

	if res == 0 {
		r.Log.Info("vm.clone Success, name : govmomi-clone-bitfusion-client01!")
	} else {
		r.Log.Info("vm.clone Failure!")
		// r.Log.WithValues("vm.clone Failure!", res)
	}

	return ctrl.Result{}, nil
}

func (r *VmbitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bitv1.Vmbit{}).
		Complete(r)
}
