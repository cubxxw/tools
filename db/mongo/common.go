// Copyright © 2024 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DefaultMaxPoolSize = 100
	DefaultMaxRetry    = 3
	DefaultConnTimeout = 5 * time.Second
)

// buildMongoURI constructs the MongoDB URI from the provided configuration.
func buildMongoURI(config *MongoConfig) string {
	credentials := ""
	if config.Username != "" && config.Password != "" {
		credentials = fmt.Sprintf("%s:%s@", config.Username, config.Password)
	}
	return fmt.Sprintf("mongodb://%s%s/%s?maxPoolSize=%d", credentials, strings.Join(config.Address, ","), config.Database, config.MaxPoolSize)
}

// shouldRetry determines whether an error should trigger a retry.
func shouldRetry(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok {
		return cmdErr.Code != 13 && cmdErr.Code != 18
	}
	return true
}

// ValidateAndSetDefaults validates the configuration and sets default values.
func (c *MongoConfig) ValidateAndSetDefaults() error {
	if c.Uri == "" && len(c.Address) == 0 {
		return errs.Wrap(errors.New("either Uri or Address must be provided"))
	}
	if c.Database == "" {
		return errs.Wrap(errors.New("database is required"))
	}
	if c.MaxPoolSize <= 0 {
		c.MaxPoolSize = DefaultMaxPoolSize
	}
	if c.MaxRetry < 0 {
		c.MaxRetry = DefaultMaxRetry
	}
	if c.ConnTimeout <= 0 {
		c.ConnTimeout = DefaultConnTimeout
	}
	if c.Uri == "" {
		c.Uri = buildMongoURI(c)
	}
	return nil
}
