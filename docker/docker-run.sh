docker rm -f kvwmap-backup
docker build -t gkaemmert/kvwmap-backup:1.0 .
docker run -d  -p 8082:8082 \
-v /home/work/github.com/kvwmap-backup/backup-config/:/app/backup-config \
--name kvwmap-backup gkaemmert/kvwmap-backup:1.0
