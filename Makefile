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


sqlc: 
	@echo "generating sqlc files"
	sqlc generate

prisma: schema.prisma
	@echo "generating schema sql file"
	npx prisma migrate diff --from-empty --to-schema-datamodel schema.prisma --script > ./migrations/schema.sql

	@echo "generating dbml file"
	npx prisma generate --schema=prisma/schema.prisma
	

