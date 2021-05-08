package main

import (
	"flag"
	"os"
	"io"
	"fmt"
	//"encoding/json"
	"regexp"
	//"reflect"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap"
	log "github.com/sirupsen/logrus"
)

var (
	listMailBoxes bool
	debug bool
)


func init() {
	// additional flags
	flag.BoolVar(&listMailBoxes, "mbs", false, "lists the mailboxes")
	flag.BoolVar(&debug, "d", true, "run in  debug mode")
	flag.Parse()

	//log level
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {

	//logging setup
	logFile, err := os.Create("imap-client.log")
	if err != nil {
		log.Warnf("Error creating logfile, %s", err)
	} else {
		logFileName := logFile.Name()
		log.Println("Created log file: ", logFileName)
		logMultiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(logMultiWriter)
	}


	//Client setup
	log.Debugln("Connecting to server...")

	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatalf("Error connecting to imap.gmail.com, %s", err)
	}
	log.Debugln("Connected")

	defer c.Logout()

	username := ""
	password := ""

	if err := c.Login(username, password); err != nil {
		log.Fatalf("Error authenticating user %s, Error:%s", username, err)
	}
	log.Println("Logged in with user ", username)

	//List mailboxes
	if listMailBoxes != false {
		ListMailBoxes(c)
	}

	// Get payment mails:
	getPaymentMails(c)

}

func getPaymentMails(c *client.Client) {
	mboxName := "INBOX"
	mbox, err := c.Select(mboxName, true)
	if err != nil {
		log.Fatalf("Error while selecting mailbox %s", mboxName, err)
	}
	log.Debugf("Selected mailbox %s", mboxName)

	//how many  mails to get:
	seqset := new(imap.SeqSet)
	from := uint32(mbox.Messages - 60)
	to := uint32(mbox.Messages)
	seqset.AddRange(from, to)

	//Fetch func  requires a channel
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	//section := &imap.BodySectionName{}
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822Text}, messages)
		//done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	//Parsing the message body
	redate := regexp.MustCompile("on .* Info:")
	reinfo := regexp.MustCompile("Info: .*\n")
	rebalance := regexp.MustCompile("The Available Balance in your account is INR [0-9]+,[0-9]+")
	repurchase := regexp.MustCompile("A purchase of INR [0-9]+,[0-9]+.[0-9]")
	recashwindrawal := regexp.MustCompile("Cash Withdrawal of INR [0-9]+,[0-9]+.[0-9]")
	for msg := range messages {

		stringMessage := fmt.Sprintln(msg)
		//log.Println(stringMessage)

		//j, err := json.Marshal(msg)
		//if err != nil {
		//	log.Fatalf("%s", err)
		//}
		//log.Println(j)


		//log.Println(msg.Envelope.Subject) this one hangs tried from https://godoc.org/github.com/emersion/go-imap#Message

		log.Println(redate.FindAllString(stringMessage, -1))
		log.Println(reinfo.FindAllString(stringMessage, -1))
		log.Println(rebalance.FindAllString(stringMessage, -1))
		log.Println(repurchase.FindAllString(stringMessage, -1))
		log.Println(recashwindrawal.FindAllString(stringMessage, -1))
	}
	if err := <-done; err != nil {
		log.Fatalf("Error while getting messages, %s", err)
	}

	return

}


func ListMailBoxes(c *client.Client) {

	MailBoxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func () {
		done <- c.List("", "*", MailBoxes)
	}()
	for m := range MailBoxes {
		log.Println("* " + m.Name)
	}
	if err := <-done; err != nil {
		log.Fatalf("Error while listing mailboxes, %s", err)
	}
	return
}
