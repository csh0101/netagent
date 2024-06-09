#!/bin/sh
# 定义路径变量
CA_DIR="./conf/ca"
SERVER_DIR="./conf/server"
CLIENT_DIR="./conf/client"

# 创建目录
mkdir -p $CA_DIR $SERVER_DIR $CLIENT_DIR


cat > $CA_DIR/ca.conf <<EOF
[ req ]
default_bits       = 4096
distinguished_name = req_distinguished_name

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = GuangDong
localityName                = Locality Name (eg, city)
localityName_default        = GuangZhou
organizationName            = Organization Name (eg, company)
organizationName_default    = Sheld
commonName                  = Common Name (e.g. server FQDN or YOUR name)
commonName_max              = 64
commonName_default          = Signed CA Test
EOF

openssl genrsa -out $CA_DIR/ca.key 4096
openssl req -new -sha256 -out $CA_DIR/ca.csr -key $CA_DIR/ca.key -config $CA_DIR/ca.conf -batch 
openssl x509 -req -days 3650 -in $CA_DIR/ca.csr -signkey $CA_DIR/ca.key -out $CA_DIR/ca.crt


cat > $SERVER_DIR/server.conf <<EOF
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = GuangDong
localityName                = Locality Name (eg, city)
localityName_default        = GuangZhou
organizationName            = Organization Name (eg, company)
organizationName_default    = Sheld
commonName                  = netagent 
commonName_max              = 64
commonName_default          = netagent 
[ req_ext ]
subjectAltName = @alt_names
[ alt_names ]
DNS.1   = localhost
IP      = 127.0.0.1
EOF


## todo server
openssl genrsa -out $SERVER_DIR/server.key 2048
openssl req -new -sha256 -out $SERVER_DIR/server.csr -key $SERVER_DIR/server.key -config $SERVER_DIR/server.conf -batch
openssl x509 -req -days 3650 -CA $CA_DIR/ca.crt -CAkey $CA_DIR/ca.key -CAcreateserial -in $SERVER_DIR/server.csr -out $SERVER_DIR/server.pem -extensions req_ext -extfile $SERVER_DIR/server.conf 

cat > $CLIENT_DIR/client.conf <<EOF
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = CN
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = GuangDong
localityName                = Locality Name (eg, city)
localityName_default        = GuangZhou
organizationName            = Organization Name (eg, company)
organizationName_default    = Sheld
commonName                  = netagent 
commonName_max              = 64
commonName_default          = netagent

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = client.local
EOF


## todo client
openssl genrsa -out $CLIENT_DIR/client.key 2048
openssl req -new -sha256 -out $CLIENT_DIR/client.csr -key $CLIENT_DIR/client.key -config $CLIENT_DIR/client.conf -batch
openssl x509 -req -days 3650 -CA $CA_DIR/ca.crt -CAkey $CA_DIR/ca.key -CAcreateserial -in $CLIENT_DIR/client.csr -out $CLIENT_DIR/client.pem -extensions req_ext -extfile $CLIENT_DIR/client.conf 

## delete temporary file 
rm $CA_DIR/ca.conf
rm $SERVER_DIR/server.conf
rm $CLIENT_DIR/client.conf