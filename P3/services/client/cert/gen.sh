rm *.pem

echo "subjectAltName=DNS:client,DNS:localhost,IP:127.0.0.1,IP:::1,IP:0.0.0.0" > client-ext.cnf

openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem \
        -subj "/C=IN/ST=Telangana/L=Hyderabad/O=IIITH/OU=College/CN=ujjwal-shekhar"

openssl x509 -req -in client-req.pem -days 60 -CA ../../common/cert/ca-cert.pem -CAkey ../../common/cert/ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client-cert.pem -noout -text