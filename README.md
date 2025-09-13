# kvm-nbd-wizard

`kvm-nbd-wizard` is a command-line tool designed to simplify the process of mounting and unmounting KVM (Kernel-based Virtual Machine) disk images using NBD (Network Block Device). This tool provides a interface to guide users through the steps of connecting, mounting, unmounting, and disconnecting virtual machine disk images.

**LVM is not supported.**

## Features

- Interactive prompts for easy operation
- Automatic detection of available KVM disk images
- Support for mounting and unmounting partitions

## Prerequisites

- KVM and QEMU installed on your system
- `qemu-nbd` command-line tool
- Root or sudo access (required for mounting operations)

## Installation

Build for tool:

```
$ go build
```

## Usage

Run the tool with root privileges:

```
$ sudo ./kvm-nbd-wizard
```

Follow the on-screen prompts to:
1. Select an action (connect, disconnect, or exit)
2. Choose a VM (for connect action)
3. Select a disk image
4. Specify NBD device and mount point
5. Mount or unmount the selected disk image

### Example

Connect:

```
$ sudo ./kvm-nbd-wizard
Select action:
1. Connect
2. Disconnect
3. Exit
Enter choice (1/2/3): 1
Available shut off VMs:
1. bookworm64
2. bookworm64-test1
Select VM number: 1
Available disk images:
1. /var/lib/libvirt/images/bookworm64.qcow2
Select disk image number (default is 1):
Using default selection: 1
Enter NBD device (default is /dev/nbd0):
NBD device /dev/nbd0 connected
Available partitions:
1. /dev/nbd0p1 (/dev/nbd0p1: Linux rev 1.0 ext4 filesystem data, UUID=c6d5ebcc-6882-47d8-a88a-fb2b01de4ed6 (extents) (64bit) (large files) (huge files)
)
2. /dev/nbd0p5 (/dev/nbd0p5: Linux swap file, 4k page size, little endian, version 1, size 249599 pages, 0 bad pages, no label, UUID=bb7d71f0-e8fb-4f46-8d46-1430a17e647a
)
Select partition number: 1
Enter mount point (default is /mnt/vm_partition):
Successfully mounted /dev/nbd0p1 to /mnt/vm_partition
Partition /dev/nbd0p1 mounted at /mnt/vm_partition
NBD device /dev/nbd0 remains connected. Use 'disconnect' option to unmount and disconnect.
Select action:
1. Connect
2. Disconnect
3. Exit
Enter choice (1/2/3):
```

Disconnect:

```
Select action:
1. Connect
2. Disconnect
3. Exit
Enter choice (1/2/3): 2
Enter NBD device (default is /dev/nbd0):
Unmounted /mnt/vm_partition
NBD device /dev/nbd0 disconnected
Select action:
1. Connect
2. Disconnect
3. Exit
Enter choice (1/2/3):
```

## Important Notes

- Always ensure that the VM is stopped before attempting to mount its disk image.
- Be careful when modifying the contents of mounted disk images to avoid data corruption.
- Remember to use the disconnect option to properly unmount and disconnect disk images when finished.
- This tool requires root privileges to perform mounting operations.

## License

This project is licensed under the [MIT License](./LICENSE).
