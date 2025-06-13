# Notes

## [SETUP] Before cloning the image, need to reset all machine identify files

### Clean up machine-specific files

sudo virt-customize -a /path/to/your/image.img \
  --delete /etc/machine-id \
  --delete /var/lib/dbus/machine-id \
  --touch /etc/machine-id \
  --delete /etc/ssh/ssh_host_rsa_key \
  --delete /etc/ssh/ssh_host_rsa_key.pub \
  --delete /etc/ssh/ssh_host_ecdsa_key \
  --delete /etc/ssh/ssh_host_ecdsa_key.pub \
  --delete /etc/ssh/ssh_host_ed25519_key \
  --delete /etc/ssh/ssh_host_ed25519_key.pub

## [SETUP] Adding permissions to libvirt images directory

### Create a dedicated subdirectory with appropriate permissions

sudo mkdir /var/libvirt/images/yourusername
sudo chown yourusername:yourusername /var/libvirt/images/yourusername

## [RUN] Qemu create image

sudo qemu-img create -b /var/lib/libvirt/images/alexng/base/focal-server-cloudimg-amd64.img -f qcow2 -F qcow2 /var/lib/libvirt/images/alexng/7a4a5c55-000c-44d5-b41e-903b71bf32fe/focal-server-cloudimg-amd64.img

## [CODE] Code rules

* Calling to other services must use saga pattern to ensure that the system is resilient to failures and can recover gracefully.
Eg: service A calls service B, if service B fails, service A should run the compensating transaction to revert the changes made by service A.
