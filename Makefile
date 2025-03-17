dev: 
	air .air.toml

# Install linux with cloudinit
install:
	virt-install --name=hal9000 --ram=2048 --vcpus=1 --import --disk path=/var/lib/libvirt/images/focal-server-cloudimg-amd64.img,format=qcow2 --disk path=/var/lib/libvirt/images/ubuntu-with-init.iso,device=cdrom --os-variant=ubuntu20.04 --network bridge=virbr0,model=virtio --graphics vnc,listen=0.0.0.0 --noautoconsole

# Create a cloudinit iso
cloudinit:
	genisoimage -output ubuntu-with-init.iso \ -V cidata -r -J user-data meta-data && \
	sudo mv -f ubuntu-with-init.iso /var/lib/libvirt/images

ip:
	sudo virsh net-dhcp-leases default

remove:
	sudo virsh destroy hal9000 && sudo virsh undefine hal9000

	sudo virt-install --name=hal9000 --ram=2048 --vcpus=1 --import --disk path=/var/lib/libvirt/images/clone-os.img,format=qcow2 --disk path=/var/lib/libvirt/images/ubuntu-with-init.iso,device=cdrom --os-variant=ubuntu20.04 --network bridge=virbr0,model=virtio --graphics vnc,listen=0.0.0.0 --import --noautoconsole
