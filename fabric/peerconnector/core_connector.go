package peerconnector

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"time"
)

type CorePeerConnector struct {
	identity *identity.X509Identity
	sign     identity.Sign
	gateway  *client.Gateway
}

func New(mspID string, certPath string, privateKeyPath string) PeerConnector {
	return &CorePeerConnector{
		identity: newIdentity(certPath, mspID),
		sign:     newSign(privateKeyPath),
		gateway:  nil,
	}
}

func NewWithRawCredentials(mspID string, cert []byte, privateKey []byte) PeerConnector {
	return &CorePeerConnector{
		identity: newIdentityFromRawCert(cert, mspID),
		sign:     newSignFromRawKey(privateKey),
		gateway:  nil,
	}
}

func (hfc *CorePeerConnector) IsConnected() bool {
	return hfc.gateway != nil
}

func (hfc *CorePeerConnector) Connect(peerEndpoint string, gatewayPeer string, tlsCertPath string) error {
	if hfc.IsConnected() {
		return fmt.Errorf("gateway is already connected")
	}

	grpcConn, err := newGrpcConnection(peerEndpoint, gatewayPeer, tlsCertPath)
	if err != nil {
		return err
	}

	gateway, err := client.Connect(
		hfc.identity,
		client.WithSign(hfc.sign),
		client.WithClientConnection(grpcConn),

		// Default timeouts for different gRPC calls
		client.WithSubmitTimeout(5*time.Second),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return err
	}

	hfc.gateway = gateway

	return nil
}

func (hfc *CorePeerConnector) ConnectWithRawTlsCert(peerEndpoint string, gatewayPeer string, tlsCert []byte) error {
	if hfc.IsConnected() {
		return fmt.Errorf("gateway is already connected")
	}

	grpcConn, err := newGrpcConnectionFromRawTlsCert(peerEndpoint, gatewayPeer, tlsCert)
	if err != nil {
		return err
	}

	gateway, err := client.Connect(
		hfc.identity,
		client.WithSign(hfc.sign),
		client.WithClientConnection(grpcConn),

		// Default timeouts for different gRPC calls
		client.WithSubmitTimeout(5*time.Second),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return err
	}

	hfc.gateway = gateway

	return nil
}

func (hfc *CorePeerConnector) GetChaincode(channelName string, chaincodeName string) (Chaincode, error) {
	network := hfc.gateway.GetNetwork(channelName)
	if network == nil {
		return nil, fmt.Errorf("channel %s not exist", channelName)
	}

	contract := network.GetContract(chaincodeName)
	if contract == nil {
		return nil, fmt.Errorf("chaincode %s not exist", chaincodeName)
	}

	return NewCoreChaincodeAPI(contract), nil
}

func (hfc *CorePeerConnector) GetChaincodeNotificationsChannel(
	ctx context.Context, channelName string, chaincodeName string,
) (<-chan *client.ChaincodeEvent, error) {
	network := hfc.gateway.GetNetwork(channelName)
	if network == nil {
		return nil, fmt.Errorf("channel %s not exist", channelName)
	}

	return network.ChaincodeEvents(ctx, chaincodeName)
}

// newSign creates a function that generates a digital
// signature from a message digest using a private key.
func newSign(keyPath string) identity.Sign {
	privateKeyPEM, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	return newSignFromRawKey(privateKeyPEM)
}

func newSignFromRawKey(key []byte) identity.Sign {
	privateKey, err := privateKeyFromPEM(key)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// newIdentity creates a client identity for this
// Gateway connection using an X.509 certificate.
func newIdentity(certPath string, mspID string) *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func newIdentityFromRawCert(cert []byte, mspID string) *identity.X509Identity {
	certificate, err := identity.CertificateFromPEM(cert)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newGrpcConnection creates the grpc connection
// used to communicate with the peer.
func newGrpcConnection(peerEndpoint string, gatewayPeer string, tlsCertPath string) (*grpc.ClientConn, error) {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return connection, nil
}

func newGrpcConnectionFromRawTlsCert(peerEndpoint string, gatewayPeer string, tlsCert []byte) (*grpc.ClientConn, error) {
	certificate, err := identity.CertificateFromPEM(tlsCert)

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return connection, nil
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

func privateKeyFromPEM(privateKeyPEM []byte) (crypto.PrivateKey, error) {
	// this function aims to support two types of generation
	// for private keys
	// * PKCS8 -> is generated by Open-SSL
	// * EC (Elliptic curves) -> is generated by golang

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse private key PEM")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		privateKeyParsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: neither is EC and PKCS8 private key")
		}
		privateKey = privateKeyParsed.(*ecdsa.PrivateKey)
	}

	return privateKey, nil
}
