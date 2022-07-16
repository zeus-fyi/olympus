
auth:
	doctl auth init

account.get:
	doctl account get

vpc.list:
	doctl vpcs list

docker.login:
	doctl registry login

chart.upgrade:
	helm upgrade chartmuseum ../../chartmuseum -f ../../chartmuseum/values.yaml

repo.add:
	helm repo add chartmuseum zeus-fyi.sfo3.digitaloceanspaces.com -u "" -p ""