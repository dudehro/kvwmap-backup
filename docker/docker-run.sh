#docker rm -f kvwmap-backup
echo "Updating images..."
docker pull golang:alpine
docker pull alpine:latest
echo "Starting image build..."
docker build -t gkaemmert/kvwmap-backup:1.0 .
echo "Starting container in foreground..."
docker run --rm  -p 8082:8082 \
-v /home/work/github.com/kvwmap-backup/backup-config/:/app/backup-config \
--name kvwmap-backup gkaemmert/kvwmap-backup:1.0
