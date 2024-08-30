FROM google/cloud-sdk:slim

RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

COPY pgdumptostorage.sh /pgdumptostorage.sh

RUN chmod +x /pgdumptostorage.sh

ENTRYPOINT ["/pgdumptostorage.sh"]