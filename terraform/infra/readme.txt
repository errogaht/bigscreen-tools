sudo ufw route allow proto tcp from any to 1.1.1.1 port 8080
sudo ufw status numbered
sudo ufw route allow proto tcp from any to any port 443 //allow port from docker container
sudo ufw allow in to any port 22 // allow port directly from host service, not docker