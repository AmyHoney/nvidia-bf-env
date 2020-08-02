package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	// "github.com/vmware/govmomi/govc/vm"
	"net/url"
)

var ctx = context.Background()

// concat vmware
func NewClient(ip,user,pwd string) *vim25.Client{

	u := &url.URL{
		Scheme: "https",
		Host:   ip,
		Path:   "/sdk",
	}
	u.User = url.UserPassword(user, pwd)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		panic(err)
	}

	fmt.Println("Log in successful")
	return client.Client

}


//find VM by name
func FindVMByName(c *vim25.Client, vmName string) (*object.VirtualMachine,error){
	m := view.NewManager(c)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		panic(err)
	}

	defer v.Destroy(ctx)

	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		panic(err)
	}

	var vmware *object.VirtualMachine
	for _, vm := range vms {
		// 判断是否含有待查找的vm，相同则为查找到的主机
		if vm.Summary.Config.Name == vmName {
			fmt.Printf("%s: %s\n", vm.Summary.Config.Name, vm.Summary.Config.GuestFullName)
			vmw := object.NewVirtualMachine(c, vm.Reference())
			fmt.Println("Find successful")
			return vmw,nil
		}
	}

	return vmware,errors.New("虚拟机不存在！")
}


func main()  {
	var ip_add = "10.110.165.188"
	var usr = "administrator@vsphere.local"
	var pwd = "Admin!23"
	var client = NewClient(ip_add,usr,pwd)
	FindVMByName(client,"photon-001")
}

