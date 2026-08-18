#!/bin/bash

# Reservas DHCP KVM/libvirt
# Ejecutar en el host Ubuntu físico.

set -e

CONEXION="qemu:///system"
RED="default"

virsh -c "$CONEXION" net-update "$RED" add ip-dhcp-host \
"<host mac='52:54:00:b0:c8:76' name='vm1' ip='192.168.122.101'/>" \
--live --config

virsh -c "$CONEXION" net-update "$RED" add ip-dhcp-host \
"<host mac='52:54:00:bf:45:e0' name='vm2' ip='192.168.122.22'/>" \
--live --config

virsh -c "$CONEXION" net-update "$RED" add ip-dhcp-host \
"<host mac='52:54:00:f7:1c:e2' name='vm3' ip='192.168.122.239'/>" \
--live --config

echo
echo "Reservas configuradas:"
virsh -c "$CONEXION" net-dumpxml "$RED" | grep "<host"
