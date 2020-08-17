package encrypt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/logger"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/ssh"
)

const (
	rsaBiteSize int = 4096
)

func hash(data []byte) []byte {
	s := sha1.Sum(data)
	return s[:]

}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	if configs.Conf.IsDebug {
		logger.AppLogger.Info("Private Key generated")
	}
	return privateKey, nil
}

func DecodeAPIKey() (string, error) {
	var (
		err    error
		apiKey []byte
	)
	apiKey, err = base64.StdEncoding.DecodeString(configs.Conf.APIKey)
	if err != nil {
		return "", err
	}
	return string(apiKey), err
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
}

func EncodePublicKeyToPEM(publicKey *rsa.PublicKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
func generatePublicSSHKey(publicKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	if configs.Conf.IsDebug {
		logger.AppLogger.Info("Public key generated")
	}

	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}
	if configs.Conf.IsDebug {
		logger.AppLogger.Info(fmt.Sprintf("Key saved to: %s", saveFileTo))
	}

	return nil
}

func GetRSAKeys(ctx context.Context, pubPath, prvPath string) (prvKey *rsa.PrivateKey, pubKey *rsa.PublicKey, err error) {
	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() (err error) {
		pubKey, err = getPublicKeyFile(pubPath)
		return err
	})

	eg.Go(func() (err error) {
		prvKey, err = getPrivateKey(prvPath)
		return err
	})

	return prvKey, pubKey, eg.Wait()
}

func getPrivateKey(file string) (*rsa.PrivateKey, error) {
	signBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}

	return signKey, nil
}

func getPublicKeyFile(file string) (*rsa.PublicKey, error) {
	verifyBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}

	return verifyKey, nil
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	if err != nil {
		logger.AppLogger.Fatal(err)
	}
}

// initKeys generate new RSA keys pair
func InitKeys(keyName, keysPath string, force bool, byteSize ...uint16) error {
	// Check if file already exists
	if _, err := os.Stat(keysPath); err != nil {
		logger.AppLogger.Error(err)
		return err
	}

	if parent := filepath.Base(keysPath); parent != "test_api" {
		keysPath = filepath.Join(filepath.Dir(keysPath), "../")
	}

	keysPath = filepath.Join(keysPath, "auth_keys")
	if _, err := os.Stat(keysPath); os.IsNotExist(err) {
		if err = os.MkdirAll(keysPath, 0755); err != nil {
			logger.AppLogger.Error(err)
			return err
		}
	}

	privateKeyPath := filepath.Join(keysPath, keyName)
	if info, err := os.Stat(privateKeyPath); err == nil && !info.IsDir() && !force {
		logger.AppLogger.Warn("keys files already exists")
		return nil
	}

	// generate new actual RSA keys
	var (
		privateKey *rsa.PrivateKey
		err        error
	)
	if len(byteSize) > 0 && ((byteSize[0]<<10)%(1<<10) == 0) {
		privateKey, err = generatePrivateKey(int(byteSize[0]))
	} else {
		privateKey, err = generatePrivateKey(rsaBiteSize)
	}
	if err != nil {
		logger.AppLogger.Error(err)
		return err
	}

	savePEMKey(filepath.Join(keysPath, keyName), privateKey)
	savePublicPEMKey(filepath.Join(keysPath, keyName+".pub"), privateKey.PublicKey)

	return nil
}

// GenerateCert generates certificate and private key based on the given host.
func GenerateCert(host string) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"I have your data"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		DNSNames:              []string{host},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader, cert, cert, &priv.PublicKey, priv,
	)

	p := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	b := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		},
	)

	return b, p, err
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, keyPath string) (ciphertext string, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "encrypt with public key")
		}
	}()
	pub, err := getPublicKeyFile(keyPath)
	if err != nil {
		return "", err
	}
	// hash := sha512.New()
	// return rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)

	bts, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bts), nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext string, keyPath string) (plaintext []byte, err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "decrypt with public key")
		}
	}()
	priv, err := getPrivateKey(keyPath)
	if err != nil {
		return nil, err
	}
	// hash := sha512.New()
	// return rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	cipherbts, err := hex.DecodeString(ciphertext)
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherbts)
}
