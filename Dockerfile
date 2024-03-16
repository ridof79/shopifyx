FROM golang:latest

COPY main /app/main

WORKDIR /app

ENV JWT_EXPIRED_MINUTES=50000
ENV DB_NAME=shopifyx_data
ENV DB_PORT=5433
ENV DB_HOST=localhost
ENV DB_USERNAME=postgres
ENV DB_PASSWORD=admin
ENV JWT_SECRET='secret1'
ENV BCRYPT_SALT=8 
ENV S3_ID=AKIA6ODU4QHXLYBDJCWL
ENV S3_SECRET_KEY=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
ENV S3_BUCKET_NAME=shopifyx

EXPOSE 8000

# Eksekusi binary ketika container dijalankan
CMD ["./main"]
