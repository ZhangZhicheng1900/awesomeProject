package main

import (
	"fmt"
	"os"

	"k8s.io/client-go/tools/clientcmd"
	clientgoapi "k8s.io/client-go/tools/clientcmd/api"
)

const haproxyIngressUserName = "haproxy-ingress"
const haproxyIngressContext = "haproxy-ingress"

func help() {
	fmt.Printf("%s clusterName clusterAddress clusterCA userToken\n", os.Args[0])
	fmt.Printf("\nexample:\n")
	fmt.Printf("%s V4_WGQ_1 https://1.1.1.1:6443 cacacaca token\n", os.Args[0])
}

func main() {
	if len(os.Args) != 5 {
		help()
		os.Exit(1)
	}
	clusterName := os.Args[1]
	clusterAddress := os.Args[2]
	clusterCA := os.Args[3]
	userToken := os.Args[4]

	c := clientgoapi.NewConfig()
	c.Clusters[clusterName] = &clientgoapi.Cluster{
		Server:                   clusterAddress,
		CertificateAuthorityData: []byte(clusterCA),
	}
	c.AuthInfos[haproxyIngressUserName] = &clientgoapi.AuthInfo{
		Token: userToken,
	}
	c.Contexts[haproxyIngressContext] = &clientgoapi.Context{
		AuthInfo: haproxyIngressUserName,
		Cluster:  clusterName,
	}
	c.CurrentContext = haproxyIngressContext

	bytes, err := clientcmd.Write(*c)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "generate kubeconfig failed, %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s", string(bytes))
}
