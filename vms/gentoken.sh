#!/bin/bash
#
# Generate and push a token to the VM
# Currently needs to be run repeatedly as token is only for 1 hour.
# Assumes PROJECT_ID, CLUSTER_LOCATION and GCE_NAME are set (and assumes VM is in same location as Cluster).

# First make sure we can ssh, otherwise scp will fail.
gcloud compute firewall-rules create default-allow-ssh --allow tcp:22

# Send a kubectl API server call to request a JWT token, get 100 hour token.
echo '{"kind":"TokenRequest","apiVersion":"authentication.k8s.io/v1","spec":{"audiences":["sven-asm-vms.svc.id.goog"],"expirationSeconds":360000}}' | kubectl create --raw /api/v1/namespaces/hipster/serviceaccounts/default/token -f - | jq -r '.status.token' > /tmp/istio-token

# SCP the token over to the VM. This places it in the home directory.
# TODO: Place directly in the appropriate location?
gcloud compute scp --project=${PROJECT_ID} --zone=${CLUSTER_LOCATION} /tmp/istio-token ${GCE_NAME}:~
