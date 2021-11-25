#!/usr/bin/env python3

map = { 'web':'pkorduan/kvwmap-server', 'mysql':'mysql', 'pgsql':'pkorduan/postgis', 'gdal':'pkorduan/gdal-http', 'mapserver':'' }
for k, v in map.items():
    print(k + ':' + v)
