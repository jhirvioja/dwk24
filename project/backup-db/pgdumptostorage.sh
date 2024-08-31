#!/bin/bash -eu

DB_USERNAME=${DB_USERNAME}
DB_PASSWORD=${DB_PASSWORD}
DB_DATABASE=${DB_DATABASE}
GCS_BUCKET=${GCS_BUCKET}
PGHOST=${PGHOST}

DUMP_FILE="${DB_DATABASE}-$(date +%Y%m%d%H%M%S).sql"

echo "Creating PostgreSQL dump.."
PGPASSWORD=$DB_PASSWORD pg_dump -U $DB_USERNAME -d $DB_DATABASE -F c -b -v -f $DUMP_FILE

echo "Uploading dump to Google Cloud Storage.."
gsutil cp $DUMP_FILE gs://$GCS_BUCKET/$DUMP_FILE

echo "Cleaning up.."
rm $DUMP_FILE

echo "Backup completed and uploaded to $GCS_BUCKET successfully."
