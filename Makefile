
dev:
	air .air.toml

# Install linux with cloudinit
install:
	sudo virt-install --name=ubuntu --ram=2048 --vcpus=1 --import --disk path=/var/lib/libvirt/images/khoakomlem/base/focal-server-cloudimg-amd64.img,format=qcow2 --disk path=/var/lib/libvirt/images/khoakomlem/cloudinit/cloudinit_fea77729-e562-4b35-aaef-e0398e9ceee3.iso,device=cdrom --os-variant=ubuntu20.04 --network bridge=virbr0,model=virtio --graphics vnc,listen=0.0.0.0 --noautoconsole

# Create a cloudinit iso
cloudinit:
	cd cloud-init-files/ubuntu &&\
	genisoimage -output ubuntu-with-init.iso -V cidata -r -J user-data meta-data network-config && \
	sudo mv -f ubuntu-with-init.iso /var/lib/libvirt/images && \
	cd ../..

ip:
	sudo virsh net-dhcp-leases default

remove:
	sudo virsh destroy ubuntu && sudo virsh undefine ubuntu

sqlc:
	sqlc generate

init-migrate:
	npx prisma migrate diff --from-empty --to-schema-datamodel prisma/schema.prisma --script > prisma/migrations/0_init/migration.sql


