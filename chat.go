package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

const chatProtocol = "/chat/1.0.0"

func getHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	targetF := flag.String("d", "", "target peer to dial")
	flag.Parse()

	rand.Seed(666)
	port1 := rand.Intn(100) + 10000
	if *targetF != "" {
		port1 = port1 + 1
	}

	h1 := makeBasicHost(port1)

	fullAddr := getHostAddress(h1)

	fmt.Println("I am ", fullAddr)
	if *targetF != "" {
		sendMessage(ctx, h1, *targetF)
	} else {
		h1.SetStreamHandler(chatProtocol, onChatRequest)
	}
	<-ctx.Done()

	fmt.Println("end")
}

func makeBasicHost(port int) host.Host {
	priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	host, _ := libp2p.New(
		libp2p.ListenAddrs(listen),
		libp2p.Identity(priv),
	)
	return host
}

func onChatRequest(s network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')
		if str == "" {
			return
		}
		if str != "\n" {
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, _ := stdReader.ReadString('\n')

		rw.WriteString(sendData + "\n")
		rw.Flush()
	}
}

func sendMessage(ctx context.Context, h1 host.Host, target string) {

	maddr, err := ma.NewMultiaddr(target)
	if err != nil {
		fmt.Println(err)
		return
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	h1.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	s, err := h1.NewStream(context.Background(), info.ID, chatProtocol)
	if err != nil {
		fmt.Println(err)
		return
	}
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go readData(rw)
	go writeData(rw)

}
