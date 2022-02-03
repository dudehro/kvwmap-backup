#docker rm -f kvwmap-backup
echo "Starting container in foreground..."
script_path=$(dirname $(readlink -f $0))
echo $script_path
docker run \
-v $script_path/bin:/app/bin \
--name kvwmap-backup gkaemmert/kvwmap-backup:0.1
