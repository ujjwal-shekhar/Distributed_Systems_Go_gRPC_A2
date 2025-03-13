rm *.pem

echo "subjectAltName=DNS:gateway,DNS:localhost,IP:127.0.0.1,IP:::1,IP:0.0.0.0" > gateway-ext.cnf

openssl req -newkey rsa:4096 -nodes -keyout gateway-key.pem -out gateway-req.pem \
        -subj "/C=IN/ST=Telangana/L=Hyderabad/O=IIITH/OU=College/CN=ujjwal-shekhar"

openssl x509 -req -in gateway-req.pem -days 60 -CA ../../common/cert/ca-cert.pem -CAkey ../../common/cert/ca-key.pem -CAcreateserial -out gateway-cert.pem -extfile gateway-ext.cnf

echo "gateway's signed certificate"
openssl x509 -in gateway-cert.pem -noout -text