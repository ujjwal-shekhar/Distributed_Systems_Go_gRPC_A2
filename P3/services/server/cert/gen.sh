rm *.pem

echo "subjectAltName=DNS:bank,DNS:localhost,IP:127.0.0.1,IP:::1,IP:0.0.0.0" > bank-ext.cnf


openssl req -newkey rsa:4096 -nodes -keyout bank-key.pem -out bank-req.pem \
        -subj "/C=IN/ST=Telangana/L=Hyderabad/O=IIITH/OU=College/CN=ujjwal-shekhar"

openssl x509 -req -in bank-req.pem -days 60 -CA ../../common/cert/ca-cert.pem -CAkey ../../common/cert/ca-key.pem -CAcreateserial -out bank-cert.pem -extfile bank-ext.cnf

echo "bank's signed certificate"
openssl x509 -in bank-cert.pem -noout -text