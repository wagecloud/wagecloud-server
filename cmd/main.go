package main

import (
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
	sourceIso := flag.String("source", "", "path to the source ISO file")
	cloudinit := flag.String("cloudinit", "", "path to the cloud-init ISO file")
	ram := flag.Int("ram", 2048, "amount of RAM in MiB")
	vcpus := flag.Uint("vcpus", 1, "number of VCPUs")

	flag.Parse()

	core.InstallDomain(core.InstallDomainParams{
		OS:            ubuntu,
		Memory:        *ram,
		VCPUs:         *vcpus,
		SourcePath:    *sourceIso,
		CloudInitPath: *cloudinit,
	})

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
