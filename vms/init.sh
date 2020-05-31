#!/bin/bash
# This script sets up environment variables and initializes gcloud

export PROJECT_ID=sven-asm-vms
export CLUSTER_NAME=vm-cluster
export CLUSTER_LOCATION=us-west3-a
export SERVICE_NAMESPACE=hipster

export GCE_NAME=vm-ig-zprn

gcloud config set project ${PROJECT_ID}
gcloud config set compute/zone ${CLUSTER_LOCATION}
