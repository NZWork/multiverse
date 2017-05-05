# Build for linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# rsync
rsync multiverse neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse
rsync index.html neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse
rsync app.js neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse
rsync ot.js neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse
rsync util.js neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse

# Deploy
ssh neo@10.1.1.7 /home/neo/docker/images/tiki_multiverse/deploy.sh
