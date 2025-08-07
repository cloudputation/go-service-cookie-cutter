FROM alpine:3

# ARGs
ARG NAME=service-seed
ARG SERVICE_USERNAME=service-seed
ARG PRODUCT_VERSION=0.0.2

# ENVs
ENV NAME=$NAME
ENV VERSION=$PRODUCT_VERSION
ENV ROOTDIR="/service-seed"
ENV SERVICE_SEED_CONFIG_FILE_PATH="/etc/service-seed/config.hcl"
ENV SERVICE_SEED_LOG_DIRECTORY="/var/log/service-seed"
ENV SERVICE_SEED_DATA_DIRECTORY="/var/lib/service-seed"
ENV TERRAFORM_PATH="/usr/local/bin/terraform"
ENV TERRAGRUNT_PATH="/usr/local/bin/terragrunt"

WORKDIR ${ROOTDIR}

# Install runtime dependencies
RUN apk add --no-cache dumb-init git

# Create service directories
RUN mkdir -p /etc/service-seed \
    && mkdir -p ${SERVICE_SEED_LOG_DIRECTORY} \
    && mkdir -p ${SERVICE_SEED_DATA_DIRECTORY}

# Set service user
RUN addgroup -g 991 ${SERVICE_USERNAME} \
    && adduser -D -u 991 -G ${SERVICE_USERNAME} ${SERVICE_USERNAME}

# Copy artifacts from builder
COPY ./API_VERSION ./API_VERSION
COPY ./artifacts/terraform ${TERRAFORM_PATH}
COPY ./artifacts/terragrunt ${TERRAGRUNT_PATH}
COPY ./build/service-seed /bin/service-seed
COPY ./.release/defaults/example.config.hcl /etc/service-seed/config.hcl
COPY ./.release/docker/docker-entrypoint.sh /bin/docker-entrypoint.sh

# Set permissions
RUN chown -R ${SERVICE_USERNAME}:${SERVICE_USERNAME} ${ROOTDIR} \
    && chown -R ${SERVICE_USERNAME}:${SERVICE_USERNAME} ${SERVICE_SEED_LOG_DIRECTORY} \
    && chown -R ${SERVICE_USERNAME}:${SERVICE_USERNAME} ${SERVICE_SEED_DATA_DIRECTORY} \
    && chmod +x /bin/docker-entrypoint.sh \
    && chmod +x ${TERRAFORM_PATH} \
    && chmod +x ${TERRAGRUNT_PATH}

# Expose port 9595
EXPOSE 9595

# Set user
USER ${SERVICE_USERNAME}

# Entrypoint to run the executable
ENTRYPOINT ["/bin/docker-entrypoint.sh"]
CMD ["/bin/service-seed"]
