#!/usr/bin/env bash

if [ "$(uname)" == "Darwin" ]; then
  echo "Detected Mac OS"
  if ! [ -x "$(command -v k3d )" ]; then
  echo "Installing k3d"
  brew install k3d
  else
    echo "k3d is already installed"
  fi
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  if ! [ -x "$(command -v k3s )" ]; then
  echo "Installing k3s"
  sudo /usr/local/bin/k3s-uninstall.sh
  curl -sfL https://get.k3s.io | sh -s - --docker
  sudo chmod 777 /etc/rancher/k3s/k3s.yaml
  cp -r /etc/rancher/k3s/k3s.yaml ~/.kube/config
  kubectl get nodes
  kubectl get ns
  else
    echo "k3s is already installed"
    k3s --version
  fi
fi