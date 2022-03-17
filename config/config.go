package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

// Config represents the configuration information.
type Config struct {
	DBName       string `json:"db_name"`       // DBName is the type of database to use. Right now, only "sqlite3" is supported
	DBPath       string `json:"db_path"`       // DBPath is the name of the database itself.
	HostAddress  string `json:"host_address"`  // HostAddress is the address to listen for connections on.
	HostPort     string `json:"host_port"`     // HostPort is the port to listen on.
	LogFile      string `json:"log_file"`      // LogFile is an optional file to log messages to.
	RelyingParty string `json:"relying_party"` // RelyingParty is the name of the WebAuthn relying party.
	RPID         string
	RPOrigin     string
}

// LoadConfig loads a configuration at the provided filepath, returning the
// parsed configuration.
func LoadConfig(filepath string) (*Config, error) {
	// Get the config file
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(configFile, config)
	return config, err
}

// SonrConfig is the configuration information loaded into the highway node instance
type SonrConfig struct {
	// HighwayAddress is the address of the Sonr Highway node.
	HighwayAddress string `json:"highway_address"`

	// HighwayPort is the port of the Sonr Highway node.
	HighwayPort int `json:"highway_port"`

	// HighwayPort is the port of the Sonr Highway node for grpc
	GrpcPort string `json:"grpc_port"`

	// HighwayPort is the port of the Sonr Highway node for http
	HttpPort string `json:"http_port"`

	// HighwayNetwork is the network of the Sonr Highway node.
	HighwayNetwork string `json:"highway_network"`

	// HighwayDID is the DID of the Sonr Highway node.
	HighwayDID string `json:"highway_did"`

	// IPFSPort is the port of the IPFS node.
	IPFSPort int `json:"ipfs_port"`

	// IPFSPath is the path of the IPFS node.
	IPFSPath string `json:"ipfs_path"`

	// LibP2PLowWater is the low water mark for the libp2p connection pool.
	LibP2PLowWater int `json:"libp2p_low_water"`

	// LibP2PHighWater is the high water mark for the libp2p connection pool.
	LibP2PHighWater int `json:"libp2p_high_water"`

	// LibP2PRendevouz is the rendevouz point for the libp2p connection pool.
	LibP2PRendevouz string `json:"libp2p_rendevouz"`

	// LibP2PBootstrapPeers is the list of bootstrap peers for the libp2p connection pool.
	LibP2PBootstrapPeers []string `json:"libp2p_bootstrap_peers"`

	// HomeDir is the home directory of the Sonr node.
	HomeDir string `json:"home_dir"`

	// CacheDir is the cache directory of the Sonr node.
	CacheDir string `json:"cache_dir"`

	// ConfigDir is the config directory of the Sonr node.
	ConfigDir string `json:"config_dir"`

	// WalletDir is the wallet directory of the Sonr node.
	WalletDir string `json:"wallet_dir"`

	// DeviceId is the device id of the Sonr node.
	DeviceId string `json:"device_id"`

	// PublicIP is the public IP of the Sonr node.
	PublicIP string `json:"public_ip"`

	// PrivateIP is the private IP of the Sonr node.
	PrivateIP string `json:"private_ip"`

	// AccountName is the account name of the Sonr node.
	AccountName string `json:"account_name"`

	// MongoUri is URI to connect to the mongodb
	MongoUri string `json:"mongo_uri"`

	// MongoCollectionName is the name of the collection we use
	MongoCollectionName string `json:"mongo_collection_name"`

	// MongoDbName is the name of the database we sue
	MongoDbName string `json:"mongo_db_name"`

	// secret key for jwts
	SecretKey string `json:"jwt"`

	//dev accoutn name for initial genesis block
	DevAccount string `json:"dev_account"`

	SqlName      string `json:"sql_name"`
	SqlPath      string `json:"sql_path"`
	RelyingParty string `json:"relying_party"`
	RPOrigin     string `json:"rp_origin"`
	RPPort       string `json:"rp_port"`
}

func (sc *SonrConfig) Save() (*SonrConfig, error) {
	viper.Set("highway.address", sc.HighwayAddress)
	viper.Set("highway.port", sc.HighwayPort)
	viper.Set("highway.network", sc.HighwayNetwork)
	viper.Set("highway.did", sc.HighwayDID)
	viper.Set("ipfs.port", sc.IPFSPort)
	viper.Set("ipfs.path", sc.IPFSPath)
	viper.Set("libp2p.lowWater", sc.LibP2PLowWater)
	viper.Set("libp2p.highWater", sc.LibP2PHighWater)
	viper.Set("libp2p.rendevouz", sc.LibP2PRendevouz)
	viper.Set("libp2p.bootstrap_peers", sc.LibP2PBootstrapPeers)
	viper.Set("home_dir", sc.HomeDir)
	viper.Set("cache_dir", sc.CacheDir)
	viper.Set("config_dir", sc.ConfigDir)
	viper.Set("wallet_dir", sc.WalletDir)
	viper.Set("device_id", sc.DeviceId)
	viper.Set("public_ip", sc.PublicIP)
	viper.Set("private_ip", sc.PrivateIP)
	viper.Set("account_name", sc.AccountName)
	err := viper.WriteConfig()
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// Return the config dir path as a Folder.
func (sc *SonrConfig) ConfigFolder() Folder {
	return Folder(sc.ConfigDir)
}

// Return the home dir path as a Folder.
func (sc *SonrConfig) HomeFolder() Folder {
	return Folder(sc.HomeDir)
}

// Return the cache dir path as a Folder.
func (sc *SonrConfig) CacheFolder() Folder {
	return Folder(sc.CacheDir)
}

// Create or return the wallet directory as a Folder.
func (sc *SonrConfig) WalletFolder() Folder {
	return Folder(sc.WalletDir)
}
