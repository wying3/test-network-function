package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// String that contains the configuration path for tnf
var configPath = flag.String("benu-config", "benu-config.yml", "path to config file")

const (
	filePerm = 0644
)

// BenuCNFConfiguration stores the BenuCNF specific test configuration.
type BenuCNFConfiguration struct {
	TrexPod                  string `json:"trexpod" yaml:"trexpod"`
	TrexContainer            string `json:"trexcontainer" yaml:"trexcontainer"`
	BNGControlPlanePod       string `json:"bngcontrolplanepod" yaml:"bngcontrolplanepod"`
	BNGControlPlaneContainer string `json:"bngcontrolplanecontainer" yaml:"bngcontrolplanecontainer"`
	BNGUserPlanePod          string `json:"bnguserplanepod" yaml:"bnguserplanepod"`
	BNGUserPlaneContainer    string `json:"bnguserplanecontainer" yaml:"bnguserplanecontainer"`
	Namespace                string `json:"namespace" yaml:"namespace"`
}

// SaveConfig writes configuration to a file at the given config path
func (c *BenuCNFConfiguration) SaveConfig(configPath string) (err error) {
	bytes, _ := yaml.Marshal(c)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(configPath, bytes, filePerm)
	return
}

// SaveConfigAsJSON writes configuration to a file in json format
func (c *BenuCNFConfiguration) SaveConfigAsJSON(configPath string) (err error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(configPath, bytes, filePerm)
	return
}

// NewConfig  returns a new decoded BenuCNFConfiguration struct
func NewConfig(configPath string) (*BenuCNFConfiguration, error) {
	var file *os.File
	var err error
	// Create config structure
	config := &BenuCNFConfiguration{}
	// Open config file
	if file, err = os.Open(configPath); err != nil {
		return nil, err
	}
	defer file.Close()
	// Init new YAML decode
	d := yaml.NewDecoder(file)
	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// parseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func parseFlags() (string, error) {
	flag.Parse()
	// Validate the path first
	if err := ValidateConfigPath(*configPath); err != nil {
		return "", err
	}
	// Return the configuration path
	return *configPath, nil
}

// GetConfig returns the Benu TestConfig configuration.
func GetConfig() (*BenuCNFConfiguration, error) {
	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := parseFlags()
	if err != nil {
		return nil, err
	}
	cfg, err := NewConfig(cfgPath)
	return cfg, err
}

// GetBenuCNFConfiguration returns the Benu CNF specific test configuration.
func GetBenuCNFConfiguration() (*BenuCNFConfiguration, error) {
	//config := &BenuCNFConfiguration{}
	config := BenuCNFConfiguration{
		TrexPod:                  "",
		TrexContainer:            "",
		BNGControlPlanePod:       "",
		BNGControlPlaneContainer: "",
		BNGUserPlanePod:          "",
		BNGUserPlaneContainer:    "",
		Namespace:                "",
	}
	return &config, nil
}
