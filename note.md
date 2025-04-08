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
/var/lib/libvirt/images/uuid   :with uuid as id of customer

