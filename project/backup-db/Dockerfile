FROM google/cloud-sdk:slim

RUN apt-get update && \
    apt install -y postgresql-common && \
    /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh -y && \
    apt-get install -y postgresql-client-16 && \
    rm -rf /var/lib/apt/lists/*

COPY pgdumptostorage.sh /pgdumptostorage.sh

RUN chmod +x /pgdumptostorage.sh

ENTRYPOINT ["/pgdumptostorage.sh"]