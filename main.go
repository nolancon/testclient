package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	podresourcesapi "k8s.io/kubelet/pkg/apis/podresources/v1"
	"k8s.io/kubernetes/pkg/kubelet/util"
	"time"
)

var timeout = 2 * time.Minute

func main() {
	socket := "unix:///var/lib/kubelet/pod-resources/kubelet.sock"
	client, _, err := getV1Client(socket, timeout, 1000)
	if err != nil {
		fmt.Println("Can't create client:", err)
		return
	}
	req := podresourcesapi.ListPodResourcesRequest{}
	resp, err := client.List(context.TODO(), &req)
	if err != nil {
		fmt.Println("Can't receive response:", err)
		return
	}
	for {
		for idx := range resp.PodResources {
			podresource := resp.PodResources[idx]
			fmt.Printf("podresource: %v", podresource)
		}
		time.Sleep(5 * time.Second)
	}
}

// getV1Client returns a client for the PodResourcesLister grpc service
func getV1Client(socket string, connectionTimeout time.Duration, maxMsgSize int) (podresourcesapi.PodResourcesListerClient, *grpc.ClientConn, error) {
	addr, dialer, err := util.GetAddressAndDialer(socket)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithContextDialer(dialer), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing socket %s: %v", socket, err)
	}
	return podresourcesapi.NewPodResourcesListerClient(conn), conn, nil
}
