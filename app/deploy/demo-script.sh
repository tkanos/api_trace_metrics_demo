# script for minikube demo

# install xhyve driver - less resource overhead for the vm
# https://github.com/kubernetes/minikube/blob/master/DRIVERS.md#xhyve-driver
#brew install docker-machine-driver-xhyve

# install minikube
# https://github.com/kubernetes/minikube/releases

# create a cluster. driver is important!
minikube stop
minikube delete
minikube start --memory 8192  --extra-config=apiserver.Authorization.Mode=RBAC --kubernetes-version v1.9.4
# --bootstrapper kubeadm
#--extra-config=controller-manager.HPAController.HorizontalPodAutoscalerUseRESTClients=false
#  --extra-config=kubelet.HorizontalPodAutoscalerUseRESTClients=false
minikube addons enable efk
minikube addons enable heapster
minikube addons enable metrics-server

#kubectl cluster-info

# Install helm
#curl -Lo /tmp/helm-linux-amd64.tar.gz https://kubernetes-helm.storage.googleapis.com/helm-v2.9.1-linux-amd64.tar.gz
#tar -xvf /tmp/helm-linux-amd64.tar.gz -C /tmp/
#chmod +x  /tmp/linux-amd64/helm && sudo mv /tmp/linux-amd64/helm /usr/local/bin/

# define helm role
kubectl -n kube-system create sa tiller
kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller

# Initialize helm, install Tiller(the helm server side component)
helm init --service-account tiller

# Make sure we get the latest list of chart
helm repo update
sleep 60s
# install Prometheus
helm install  stable/prometheus --name prometheus --namespace monitoring --set server.service.type=NodePort  

# install grafana
#helm install -f grafana-values.yaml stable/grafana --name grafana --namespace monitoring

# install zipkin
kubectl apply -f ./zipkin.yaml -n monitoring


# install Pixel App
cd /home/f_dutratineesilva/Projects/go/src/github.com/tkanos/api_trace_metrics_demo/app
eval $(minikube docker-env)
docker build -t pixel-api ../.
helm install --name pixel-api ./pixel-api --set service.type=NodePort,deployment.env.zipkin='http://zipkin.monitoring.svc.cluster.local' --namespace pixel
