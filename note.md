// Create a dedicated subdirectory with appropriate permissions
sudo mkdir /var/libvirt/images/yourusername
sudo chown yourusername:yourusername /var/libvirt/images/yourusername


// Customize libvirt file system for running with root perm
* sudo vi /etc/libvirt/libvirtd.conf
// Uncomment the following lines and modify


auth_unix_ro = "none" // need authen for read only?
auth_unix_rw = "none" // need authen for read/write?
unix_sock_group = "root" // uncomment this line
unix_sock_ro_perms = "0777" // uncomment this line
unix_sock_rw_perms = "0770" // uncomment this line


// destionation for customer's images
/var/lib/libvirt/images/{$USER}/uuid   :with uuid as id of customer and USER as username, example: /var/lib/libvirt/images/alexng/7a4a5c55-000c-44d5-b41e-903b71bf32fe


// qemu create image
sudo qemu-img create -b /var/lib/libvirt/images/alexng/base/focal-server-cloudimg-amd64.img -f qcow2 -F qcow2 /var/lib/libvirt/images/alexng/7a4a5c55-000c-44d5-b41e-903b71bf32fe/focal-server-cloudimg-amd64.img


// change permission (recursively) for /etc/nginx/users.d/*
sudo chown -R yourusername:yourusername /etc/nginx/users.d/*

// add nginx to sudoers file(for user access to)
sudo EDITOR=vi visudo -> username ALL=(ALL) NOPASSWD: /usr/sbin/nginx




