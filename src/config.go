/*
 * Handle Pond-Mail-Gateway configuration settings.
 *
 * (c) 2013-2014 Bernd Fix   >Y<
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or (at
 * your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

///////////////////////////////////////////////////////////////////////
// Import external declarations.

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

///////////////////////////////////////////////////////////////////////
// Configuration types

/*
 * Config data for POP3 server
 */
type EmailConfig struct {
	POP3    string        `json:"pop3"` // POP3 URL
	SMTP    string        `json:"smtp"` // SMTP URL
	Address string        `json:"address"`
	Poll    time.Duration `json:"pollInterval"`
}

//---------------------------------------------------------------------
/*
 * Config data for Pond client
 */
type PondConfig struct {
	Home      string `json:"home"`
	StateFile string `json:"stateFile"`
	StatePW   string `json:"statePW"`
}

//---------------------------------------------------------------------
/*
 * Message templates
 */
type TemplateConfig struct {
	MailRegFailure string `json:"mailRegFailure"`
	MailRegSuccess string `json:"mailRegSuccess"`
	MailPending    string `json:"mailPending"`
	ValidateMail   string `json:"validateMail"`
	PondRegSuccess string `json:"pondRegSuccess"`
}

//---------------------------------------------------------------------
/*
 * Database-related configuration settings
 */
type DatabaseConfig struct {
	Connect          string `json:"connect"` // database connection string
	InsertMailUser   string `json:"insertMailUser"`
	SelectMailUser   string `json:"selectMailUser"`
	UpdateMailStatus string `json:"updateMailStatus"`
	InsertPondUser   string `json:"insertPondUser"`
	SelectPondUser   string `json:"selectPondUser"`
	UpdatePondStatus string `json:"updatePondStatus"`
}

//---------------------------------------------------------------------
/*
 * Web-interface settings
 */
type WebConfig struct {
	Host        string `json:"host"`
	Key         string `json:"key"`
	Cert        string `json:"cert"`
	FormPage    string `json:"formPage"`
	ErrorPage   string `json:"errorPage"`
	CaptchaFail string `json:"captchaFail"`
}

//---------------------------------------------------------------------
/*
 * Config data for control server
 */
type ControlConfig struct {
	Port    int    `json:"port"`    // control service port
	Allowed string `json:"allowed"` // list of allowed IP addresses
}

//---------------------------------------------------------------------
/*
 * Combined configuration data
 */
type Config struct {
	Control  *ControlConfig  `json:"control"`
	Database *DatabaseConfig `json:"database"`
	Web      *WebConfig      `json:"webif"`
	Proxy    string          `json:"proxy"` // SOCKS5 URL
	Email    *EmailConfig    `json:"email"`
	Pond     *PondConfig     `json:"pond"`
	Tpls     *TemplateConfig `json:"templates"`
}

///////////////////////////////////////////////////////////////////////
// Configuration methods

/*
 * A JSON-encoded configuration file is parsed and mapped to the
 * Config data structure.
 * @param fileName string - name of configuration file
 * @return config *Config - parsed configuration data (or nil if failed)
 * @return err error - error object in case of failure
 */
func parseConfig(fileName string) (config *Config, err error) {

	// pasre configuration file
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// unmarshal to Config data structure
	config = new(Config)
	err = json.Unmarshal(file, config)
	return
}
