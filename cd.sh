# Build for linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# rsync
rsync sync neo@10.1.1.7:/home/neo/docker/images/tiki_sync
rsync index.html neo@10.1.1.7:/home/neo/docker/images/tiki_sync
rsync app.js neo@10.1.1.7:/home/neo/docker/images/tiki_sync
rsync ot.js neo@10.1.1.7:/home/neo/docker/images/tiki_sync
rsync util.js neo@10.1.1.7:/home/neo/docker/images/tiki_sync

# Deploy
ssh neo@10.1.1.7 /home/neo/docker/images/tiki_sync/deploy.sh
