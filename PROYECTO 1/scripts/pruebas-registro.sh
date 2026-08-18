#!/bin/bash

# Validación del registro privado Zot
# Carnet: 202112092

ZOT="http://192.168.122.239:5000"

echo "======================================"
echo " REGISTRO PRIVADO ZOT"
echo "======================================"
echo

echo "=== Estado Registry ==="
curl -s "$ZOT/v2/"
echo
echo

echo "=== Catálogo ==="
curl -s "$ZOT/v2/_catalog"
echo
echo

echo "=== API1 ==="
curl -s "$ZOT/v2/api1-202112092/tags/list"
echo
echo

echo "=== API2 ==="
curl -s "$ZOT/v2/api2-202112092/tags/list"
echo
echo

echo "=== API3 ==="
curl -s "$ZOT/v2/api3-202112092/tags/list"
echo
