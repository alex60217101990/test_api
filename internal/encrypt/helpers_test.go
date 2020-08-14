package encrypt

import (
	"fmt"
	"os"
	"testing"

	"github.com/alex60217101990/test_api/internal/configs"
	"github.com/alex60217101990/test_api/internal/logger"
)

func init() {
	// Load configs file
	configs.Conf = &configs.Configs{
		LoggerType: configs.Zap,
	}
	// Init loggers
	logger.InitLoggerSettings()
}

func TestInitKeys(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		logger.AppLogger.Fatal(err)
	}
	InitKeys("test-api.rsa", currentDir, false)
	InitKeys("test-api-repo.rsa", currentDir, true, 1024)
}

func TestDecryptPassword(t *testing.T) {
	// chiper, err := EncryptWithPublicKey([]byte("pGE oItaZ#aV%xZ0KWd$hJ5q#aM y%9P"), "../../auth_keys/test-api.rsa.pub")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// fmt.Println("IN: ", chiper)
	// fmt.Println("EQUIL: ", chiper == "44140c0bdab3291c93426ef91fb20d4d4664b3e3b2de922cc825ab7d92584e9f186269746b00f965253314ab2fe32813c3dfe9709d395501491bb9a90182b49e0cdd1e3f0ff992eba0ec7432acdb3fc8384ec52ec02ca078166290ba7062eb12357dee9a06e34d94cda661f431172865080ccc81c0d0c482a8909847121e5e96")
	chiper := "32d7deebf665af32215960a98a8e9c5831769413a86840413fa4e21147aa0ebc7779ec3d3b5ed0121c83af58e16e39d643d3cc752579718b84001f9e9ad33cb7ae3cd2299e950f4ac32b49c8098fe00ec37c2a198d7f4c26993ddca8a6b91ee29f5891f09adbc06dcf561ee243ff8a30dbc8f47a1c579ac990dc1efc33da8ff7"
	bts, err := DecryptWithPrivateKey(
		chiper,
		"../../auth_keys/test-api-repo.rsa",
	)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("OUT: ", string(bts))
	t.Log(string(bts))
}
