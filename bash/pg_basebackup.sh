#!/bin/bash

pg_basebackup --pgdata=/var/www/postgres_basebackup --xlog-method=stream --progress -d postgres://postgres@localhost
