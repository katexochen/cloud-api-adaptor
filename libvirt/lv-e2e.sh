#!/usr/bin/env bash

set -exuo pipefail

rm -rf /root/.kcli/clusters/peer-pods
virsh destroy peer-pods-ctlplane-0 || true
virsh destroy peer-pods-worker-0 || true
virsh undefine peer-pods-ctlplane-0 --remove-all-storage || true
virsh undefine peer-pods-worker-0 --remove-all-storage || true
podvms=$(virsh list --all | tail -n +3 | cut -d' ' -f5 | grep 'podvm') || true
for podvm in $podvms; do
    virsh destroy "$podvm" || true
    virsh undefine "$podvm" --nvram --remove-all-storage || true
done
podvms=$(virsh list --all | tail -n +3 | cut -d' ' -f6 | grep 'podvm') || true
for podvm in $podvms; do
    virsh destroy "$podvm" || true
    virsh undefine "$podvm" --nvram --remove-all-storage || true
done
virsh net-destroy default || true

kvm-ok

if ! sudo virsh pool-list --all | grep default >/dev/null; then
    sudo virsh pool-define-as default dir - - - - "/var/lib/libvirt/images"
    sudo virsh pool-build default
fi
sudo virsh pool-start default || true
sudo virsh net-start default || true

sudo setfacl -m "u:${USER}:rwx" /var/lib/libvirt/images
sudo adduser "$USER" libvirt
sudo setfacl -m "u:${USER}:rwx" /var/run/libvirt/libvirt-sock

[[ -f ~/.ssh/id_rsa ]] || ssh-keygen -t rsa -f ~/.ssh/id_rsa -N ""

[[ -f ./install/overlays/libvirt/id_rsa ]] ||
    ssh-keygen -f ./install/overlays/libvirt/id_rsa -N "" &&
    cat ./install/overlays/libvirt/id_rsa.pub >> ~/.ssh/authorized_keys &&
    chmod 600 ~/.ssh/authorized_keys

IP="$(hostname -I | cut -d' ' -f1)"
virsh -c "qemu+ssh://$USER@${IP}/system?keyfile=$(pwd)/id_rsa&no_verify=1" nodeinfo

echo "libvirt_uri=\"qemu+ssh://$USER@${IP}/system?no_verify=1\"" > libvirt.properties
echo "libvirt_ssh_key_file=\"id_rsa\"" >> libvirt.properties
cat libvirt.properties

export PATH=$PATH:/usr/local/go/bin
export CLOUD_PROVIDER=libvirt
export TEST_PROVISION="yes"
export TEST_TEARDOWN="no"
export TEST_PROVISION_FILE="$PWD/libvirt.properties"
export TEST_PODVM_IMAGE="${PWD}/system.qcow2"
export TEST_E2E_TIMEOUT="50m"

make test-e2e
