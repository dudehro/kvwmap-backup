#!/bin/bash
pg_dump -Fc -U ${PGUSER} -h ${PGHOST} -f /dumps/${PGDATABASE}.dump ${PGDATABASE}
