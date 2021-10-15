docker rm backup-interface
docker build -t gkaemmert/backup-interface:1.0 .
docker run -p 8082:8082 \
-v /home/work/go/src/github.com/kvwmap-backup/backup-config/:/backup-interface/backup-config \
--name backup-interface gkaemmert/backup-interface:1.0
