package config

import (
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type Config struct {
    DBPath string `yaml:"db_path"`
}

func LoadConfig() (*Config, error) {
    dir, err := os.Getwd()
    if err != nil {
        return nil, err
    }
    for {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            break // found project root
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            return nil, os.ErrNotExist
        }
        dir = parent
    }

    configPath := filepath.Join(dir, "config.yaml")
    f, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var cfg Config
    decoder := yaml.NewDecoder(f)
    if err := decoder.Decode(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
