/*
 * Main control loop:
 * - Poll POP3 mailbox of gateway periodically for incoming mail messages.
 * - Poll Pond server periodically for incoming Pond messages.
 * - Handle incoming messages
 * - Respond to notifications from the administrative control server.
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
	"./pond"
	"code.google.com/p/go.crypto/openpgp"
	"crypto/rand"
	"database/sql"
	"flag"
	"fmt"
	"github.com/bfix/gospel/logger"
	"github.com/bfix/gospel/network"
	"io"
	"strconv"
	"time"
)

///////////////////////////////////////////////////////////////////////
// package-global variables

type Global struct {
	configFile string
	logFile    string
	config     *Config
	prng       io.Reader

	prvkey   string
	identity *openpgp.Entity
	pubkey   []byte

	db *sql.DB

	client *pond.Client
}

var (
	g Global
)

///////////////////////////////////////////////////////////////////////
// Main methods

/*
 * Initialize parameters
 */
func init() {
	flag.StringVar(&g.configFile, "c", "pondgw.conf", "configuration file")
	flag.StringVar(&g.logFile, "l", "", "log file")
	flag.StringVar(&g.prvkey, "k", "private.asc", "private key file")
	flag.Parse()
}

//---------------------------------------------------------------------
/*
 * Application entry point
 */
func main() {
	fmt.Println("================================")
	fmt.Println("POND/EMAIL GATEWAY,    v0.1")
	fmt.Println("(c) 2014 by Bernd Fix   >Y<")
	fmt.Println("================================")

	if len(g.logFile) > 0 {
		logger.LogToFile(g.logFile)
	}

	// handle configuration file
	logger.Println(logger.INFO, "Reading configuration file '"+g.configFile+"'...")
	var err error
	g.config, err = parseConfig(g.configFile)
	if err != nil {
		fmt.Println("Failed to read configuration file!")
		fmt.Println("*** " + err.Error())
		logger.Println(logger.ERROR, "Failed to read configuration file: '"+err.Error())
		return
	}

	// initialize modules (and global parameters)
	g.prng = rand.Reader
	if err = InitPondModule(); err == nil {
		if err = InitMailModule(); err == nil {
			err = InitUserModule()
		}
	}
	if err != nil {
		logger.Println(logger.ERROR, err.Error())
		return
	}

	// create control service and start control daemon
	ch := make(chan bool)
	ctrl := &ControlSrv{ch}
	ctrlList := []network.Service{ctrl}
	logger.Printf(logger.INFO, "Starting control daemon on port %d...\n", g.config.Control.Port)
	network.RunService("tcp", ":"+strconv.Itoa(g.config.Control.Port), ctrlList)

	// prepare and run handlers
	mailMsgIn := make(chan MailMessage)
	mailCtrl := make(chan int)
	go PollMailServer(mailMsgIn, mailCtrl)
	
	//g.client.StartKeyExchange("#1", "23MasterOfDesasterRulez!")

	heartbeat := time.NewTicker(6 * time.Hour)
	for {
		select {
		// handle mail message
		case msg := <-mailMsgIn:
			if err = HandleIncomingMailMessage(msg); err != nil {
				logger.Println(logger.ERROR, err.Error())
			}

		// request for termination
		case <-ch:
			logger.Println(logger.INFO, "Terminating application.")
			mailCtrl <- MAIL_CMD_QUIT
			return

		// handle heartbeat and drop timed-out sessions
		case now := <-heartbeat.C:
			logger.Println(logger.INFO, "Heartbeat: "+now.String())
		}
	}
}