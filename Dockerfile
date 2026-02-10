FROM alpine:3

# ARGs - Build-time configuration
ARG NAME="service-seed"
ARG AUTHOR="Francis Robinson"
ARG OWNER="Behavox"
ARG SERVICE_USERNAME="service-seed"
ARG PRODUCT_VERSION=0.0.2
ARG SERVICE_PORT=9595
ARG CONFIG_ENV=local

# ENVs - Runtime configuration
ENV NAME=$NAME
ENV VERSION=$PRODUCT_VERSION
ENV SERVICE_SEED_ROOTDIR="/service-seed"
ENV SERVICE_USERNAME="service-seed"
ENV SERVICE_PORT=${SERVICE_PORT}
ENV SERVICE_SEED_CONFIG_PATH="/etc/service-seed/config.hcl"
ENV SERVICE_SEED_LOG_DIRECTORY="/var/log/service-seed"
ENV SERVICE_SEED_DATA_DIRECTORY="/var/lib/service-seed"
ENV TERRAFORM_PATH="/usr/local/bin/terraform"
ENV TERRAGRUNT_PATH="/usr/local/bin/terragrunt"

WORKDIR ${SERVICE_SEED_ROOTDIR}

# Install runtime dependencies (dumb-init for proper signal handling)
RUN apk add --no-cache dumb-init git

# Create service directories
RUN mkdir -p /etc/service-seed/config.d \
    && mkdir -p ${SERVICE_SEED_LOG_DIRECTORY} \
    && mkdir -p ${SERVICE_SEED_DATA_DIRECTORY}

# Set service user (non-root for security)
RUN addgroup -g 991 ${SERVICE_USERNAME} \
    && adduser -D -u 991 -G ${SERVICE_USERNAME} ${SERVICE_USERNAME}

# Copy artifacts
COPY ./API_VERSION ./API_VERSION
COPY ./artifacts/terraform ${TERRAFORM_PATH}
COPY ./artifacts/terragrunt ${TERRAGRUNT_PATH}
COPY ./build/service-seed /bin/service-seed
COPY ./.release/defaults/config/config.${CONFIG_ENV}.hcl /etc/service-seed/config.hcl
COPY ./config.d/ /etc/service-seed/config.d/

# Copy application entrypoint
COPY ./.release/docker/docker-entrypoint.sh /bin/docker-entrypoint.sh

# Set permissions
RUN chown -R ${SERVICE_USERNAME}:${SERVICE_USERNAME} ${SERVICE_SEED_ROOTDIR} \
    && chown -R ${SERVICE_USERNAME}:${SERVICE_USERNAME} /etc/service-seed \
    && chown -R ${SERVICE_USERNAME}:${SERVICE_USERNAME} ${SERVICE_SEED_LOG_DIRECTORY} \
    && chown -R ${SERVICE_USERNAME}:${SERVICE_USERNAME} ${SERVICE_SEED_DATA_DIRECTORY} \
    && chmod +x /bin/docker-entrypoint.sh \
    && chmod +x ${TERRAFORM_PATH} \
    && chmod +x ${TERRAGRUNT_PATH}

# Expose port
EXPOSE ${SERVICE_PORT}

# Set user (run as non-root)
USER ${SERVICE_USERNAME}

# Entrypoint with dumb-init for proper signal handling
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/bin/docker-entrypoint.sh", "/bin/service-seed"]
