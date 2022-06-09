FROM eclipse-mosquitto:1.6.10

WORKDIR /app

# Copy the certificate authority's certificate into the container's trusted authorities
COPY ./config/certs/ca.crt /usr/local/share/ca-certificates/ca.crt

# Update the container's trusted authorities to include the self-signed certificate
RUN update-ca-certificates
