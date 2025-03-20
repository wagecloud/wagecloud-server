package main

import (
	"fmt"
	"log"

	"github.com/wagecloud/wagecloud-server/internal/service/cloudinit"
	"github.com/wagecloud/wagecloud-server/internal/service/libvirt"
)

func main() {
	passwd, err := cloudinit.HashPassword("061204", cloudinit.BCrypt)
	if err != nil {
		log.Fatalf("failed to hash password: %s", err)
	}

	customer := &cloudinit.Customer{
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

	err = cloudinit.WriteCloudInitFiles(customer, "./cloud-init-files/ubuntu/vm1")
	if err != nil {
		log.Fatalf("failed to write cloud-init files: %s", err)
	}

	// err = cloudinit.GenISO("./cloud-init-files/ubuntu/vm1", "/var/lib/libvirt/vm1.iso")
	// if err != nil {
	// 	log.Fatalf("failed to generate cloud-init ISO: %s", err)
	// }

	domain, err := libvirt.CreateDomain(&libvirt.Spec{
		CPU: 1,
		Memory: &libvirt.Memory{
			Value: 2048,
			Unit:  libvirt.MiB,
		},
		OS: libvirt.OSTypeUbuntu,
	})

	if err != nil {
		log.Fatalf("Failed to create domain: %v", err)
	}

	defer domain.Free()

	err = libvirt.StartDomain(domain)
	if err != nil {
		log.Fatalf("Failed to start domain: %v", err)
	}

	fmt.Println("Started domain")
}
