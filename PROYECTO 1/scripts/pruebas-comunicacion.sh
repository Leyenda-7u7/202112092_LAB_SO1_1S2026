#!/bin/bash

# Prueba completa de comunicación entre API1, API2 y API3
# Carnet: 202112092

API1="http://192.168.122.101:8081"
API2="http://192.168.122.101:8082"
API3="http://192.168.122.22:8083"

echo "========================================"
echo " PRUEBA DE LAS 6 COMUNICACIONES REST"
echo "========================================"
echo

echo "===== API1 -> API2 ====="
curl -s "$API1/api1/202112092/call-api2"
echo
echo

echo "===== API1 -> API3 ====="
curl -s "$API1/api1/202112092/call-api3"
echo
echo

echo "===== API2 -> API1 ====="
curl -s "$API2/api2/202112092/call-api1"
echo
echo

echo "===== API2 -> API3 ====="
curl -s "$API2/api2/202112092/call-api3"
echo
echo

echo "===== API3 -> API1 ====="
curl -s "$API3/api3/202112092/call-api1"
echo
echo

echo "===== API3 -> API2 ====="
curl -s "$API3/api3/202112092/call-api2"
echo
echo

echo "========================================"
echo " FIN DE PRUEBAS"
echo "========================================"
