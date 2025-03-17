package main

import (
	"GoVirService/internal/cloud-init"
	"GoVirService/internal/core"
	"GoVirService/internal/model"
	"flag"
	"fmt"
	"log"
	"os"
)

var ubuntu = model.OS{
	Name: "UbuntuFocal",
}

func main() {
	// virt-install --name=hal9000 --ram=2048 --vcpus=1 --import --disk
	// path=/var/lib/libvirt/images/focal-server-cloudimg-amd64.img,format=qcow2
	//  --disk path=/var/lib/libvirt/images/ubuntu-with-init.iso,device=cdrom
	// --os-variant=ubuntu20.04 --network bridge=virbr0,model=virtio
	// --graphics vnc,listen=0.0.0.0 --noautoconsole
	// sourceIso := flag.String("source", "", "path to the source ISO file")
	// cloudinit := flag.String("cloudinit", "", "path to the cloud-init ISO file")
	// ram := flag.Int("ram", 2048, "amount of RAM in MiB")
	// vcpus := flag.Uint("vcpus", 1, "number of VCPUs")
	// flag.Parse()
	//
	// core.InstallDomain(core.InstallDomainParams{
	// 	OS:            ubuntu,
	// 	Memory:        *ram,
	// 	VCPUs:         *vcpus,
	// 	SourcePath:    *sourceIso,
	// 	CloudInitPath: *cloudinit,
	// })

	passwd, err := cloudinit.HashPassword("061204", cloudinit.BCrypt)
	if err != nil {
		log.Fatalf("failed to hash password: %s", err)
	}

	customer := cloudinit.Customer{
		Name:   "trietdeptrai",
		Passwd: passwd,
		SSHKeys: []string{
			//khoakomlem
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCisXMl0mvnYlCfwjsDa6WEqfXQ30MYkN4f5OUPKojc7OsAIP50soK9nn0r0fzAnJ/ze5szdD3XtqEvZ0VF/DFChMi8EDgSFwL6IfsSh1+56Nl1XHCXeJWZLyHkzpsDIwS7e38l25j95sWKL55yup8geB2d03YVTEdWMhhlq89Q8GN8xDVpXy9SKgFrQTP8BCmi7UZZX03LTE3kXepQ3TZgPPZZUtUFFKNA8eq4qlofOF//9u6xxrh7eOi6a8WVuxLNIOMjCoyfD7xbdLY+Vw7d2hCzjKeU1dKAD3ZGIxPllYn+eDS7hJm1pZP1WMwdEJUywQJ25Fu557s2nlc6DnXl",
			//alexng
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCkMgZmAp5Wz+vd3jLgiMdZrT+W2D1UjUnJ/mKL/QpimQ7HfrKM6ejtIwg9LwHlvC+GVnoEDdeKq/LwpBIs4Vtqnu+ZkJTe36ew2szjaY3k7fzxpfePVsoPPu5zovsGCOA88c8vnZPjtSJXmpXWtZFIXGRUQbj6EbmPU3wuXseg/23d/ZQmwvhlFhxB2pNVBjza1hXhDFzA62wQ/zMMEqgUm4uztSb8heyLj6Gc9cVg+kVcLqEyMdN1/IFLlN+8LgRNwssuA0r7E31WjfmtZZvpE0peptFgWm5/yhlF/jSyNyrWLfUSB+HmXu5THDOScvJyiWa7UieUMXaja/h5fEc4ZIxLWL8vSsm8W8YOpf3Pz/oHf0dGJTsf+qE9lwW3b+tAApbA+u49ZZ8aebHcH14uTCZcdVr28qSuZds13RILL/1n0qaoTmGwa4W6OV8LVuH3Ijp1DpTY+yKgwBIeIgaY++EuhWk9f4k2jsJdudu6pHUaQhod2MM0Y63zD/roBhkaV6/qGHTYa/RLinj8k4TqgfH9YAF8+BFUgDr08/221aPsA3OU6ARemxfuZPpr8zS6HLWV1xu86HKpzcZ3r6plyiM8iuR1Ldlc4iXghicNKFqJ1tVClfcCqluD5lw1EALhpgWYgfC2Zf8ESWiOmO3ulxC7RBaYVWApzljd/KnMMw== alexng@pop-os",
		},
		LockPasswd: false,
	}

	err = cloudinit.WriteCloudInitFiles(&customer, "./cloud-init-files/ubuntu")
	if err != nil {
		log.Fatalf("failed to write cloud-init files: %s", err)
	}
}

func test2() {
	var userdata = flag.String("userdata", "", "path to the user data file")
	var metadata = flag.String("metadata", "", "path to the metadata file")

	flag.Parse()

	output := flag.Arg(0)

	// resovle to our pwd
	fmt.Println("Current working directory: ", os.Getenv("PWD"))
	fmt.Println("Output: ", output)
	iso, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer iso.Close()

	userdataFile, err := os.Open(*userdata)
	if err != nil {
		panic(err)
	}
	defer userdataFile.Close()

	metadataFile, err := os.Open(*metadata)
	if err != nil {
		panic(err)
	}
	defer metadataFile.Close()

	core.GenISO(core.GenISOParams{
		Userdata:  userdataFile,
		Metadata:  metadataFile,
		ResultIso: iso,
	})

	iso.Close()
	log.Printf("ISO image created successfully")
}
