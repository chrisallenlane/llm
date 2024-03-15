// Package config encapsulates app config
package config

import "gorm.io/gorm"

// Config models a configuration value
type Config struct {
	gorm.Model
	Name  string `gorm:"unique"`
	Value string
}
