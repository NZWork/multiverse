# Build for linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

rsync -re 'ssh -p 1027' multiverse neo@amoy.layer.nevoz.com:/home/neo/docker/images/tiki_multiverse
rsync -re 'ssh -p 1027' app.js neo@amoy.layer.nevoz.com:/home/neo/docker/images/tiki_multiverse
ssh neo@amoy.layer.nevoz.com -p 1027 /home/neo/docker/images/tiki_multiverse/deploy.sh

# rsync
#rsync multiverse neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse
#rsync index.html neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse
#rsync app.js neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse
#rsync ot.js neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse
#rsync util.js neo@10.1.1.7:/home/neo/docker/images/tiki_multiverse

# Deploy
#ssh neo@10.1.1.7 /home/neo/docker/images/tiki_multiverse/deploy.sh
