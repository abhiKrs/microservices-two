if [  -n "$(uname -a | grep Ubuntu)" ]; then
echo "Updating Ubuntu system"
sudo apt install virtualbox virtualbox-ext-pack -y
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install curl
sudo apt-get install apt-transport-https
sudo apt-get -y install certbot python3-certbot-apache
sudo certbot certonly
echo "Installing VirtualBox and its components"
sudo apt install virtualbox virtualbox-ext-pack
if ! [ -x "$(command -v docker)" ]; then
echo "Installing docker"
sudo apt install docker.io
else
  echo "docker is already installed"
  docker version
fi
if ! [ -x "$(command -v minikube)" ]; then
echo "Installing minikube"
wget https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo cp minikube-linux-amd64 /usr/local/bin/minikube
sudo chmod 755 /usr/local/bin/minikube
minikube version
else
  echo "minikube is already installed"
  minikube version
fi
if ! [ -x "$(command -v kubectl)" ]; then
echo "Installing kubectl"
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version -o json
else
  echo "kubectl is already installed"
  kubectl version -o json
fi
minikube start --force
if ! [ -x "$(command -v helm)" ]; then
echo "Installing helm"
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
sh get_helm.sh
helm version
else
  echo "helm is already installed"
  helm version
fi
else
#######################################################################################################
if ! [ -x "$(command -v kubectl)" ]; then
echo "Installing kubectl in centos "
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version -o json
else
  echo "kubectl is already installed"
fi
if ! [ -x "$(command -v kubectl)" ]; then
  echo "Adding Kubernetes yum repository..."
  sudo bash -c "cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF"

  echo "Installing kubectl..."
  sudo yum install -y kubectl
else
  echo "kubectl is already installed"
fi
##########################################################################################################
if ! [ -x "$(command -v minikube)" ]; then
  echo "Installing minikube..."
  curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.28.2/minikube-linux-amd64
  chmod +x minikube
  sudo cp minikube /usr/local/bin/

  echo "Configuring minikube..."
  minikube config set ShowBootstrapperDeprecationNotification false
  minikube config set WantUpdateNotification false
  minikube config set WantReportErrorPrompt false
  minikube config set WantKubectlDownloadMsg false
else
  echo "minikube is already installed"
fi

echo "Starting minikube..."
sudo /usr/local/bin/minikube --vm-driver=none --bootstrapper=localkube start

set tempFolder=$PWD

cd ~

rm .kube -rf
rm .minikube -rf

sudo mv /root/.kube $HOME/
sudo chown -R $USER $HOME/.kube
sudo chgrp -R $USER $HOME/.kube

sudo mv /root/.minikube $HOME/
sudo chown -R $USER $HOME/.minikube
sudo chgrp -R $USER $HOME/.minikube

cd $tempFolder

sed -i "s:/root/:/home/$USER/:g" $HOME/.kube/config
fi