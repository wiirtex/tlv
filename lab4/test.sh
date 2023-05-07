# /bin/sh

sleep 1
ifconfig veth1 10.1.1.2/24 up
mount -t proc proc /proc

touch container.disk
mkdir containerfs

dd if=/dev/zero of=container.disk count=1024
losetup -f container.disk


mkfs -t ext4 /dev/loop0
mount /dev/loop0 containerfs

echo "TLV lab4" > containerfs/tlv.txt
cat containerfs/tlv.txt

sh