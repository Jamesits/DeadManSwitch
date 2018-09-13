package main

import (
	"github.com/miekg/dns"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func check(conf *config) uint8{
	req := new(dns.Msg)
	req.SetQuestion(dns.Fqdn(conf.Record), stringToDnsType(conf.RecordType))

	client := new(dns.Client)
	client.SingleInflight = true
	in, rtt, err := client.Exchange(req, "1.2.4.8:53")
	log.Printf("RTT: %d\n", rtt)

	// unexpected result
	if err != nil {
		log.Println("Unable to resolve record")
		return Uncertain
	}

	expectedValueFound := false
	log.Println("Result entries:")
	for _, entry := range in.Answer {
		entryString := entry.String()
		log.Println(entry)
		if strings.Contains(entryString, conf.ExpectedValue) {
			// Yay we found the correct value
			log.Println("Normal value matched")
			expectedValueFound = true
		}
		if strings.Contains(entryString, conf.TriggerValue) {
			// Something happened
			log.Println("Trigger value matched")
			return True
		}
	}

	if expectedValueFound {
		return False
	} else {
		return Uncertain
	}
}

func runScriptIterative(path string) {
	file, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("File %s not found or insufficient privilege, skipping...\n", path)
	} else {
		if file.IsDir() {
			log.Printf("Entering directory %s\n", path)
			f, _ := os.Open(path)
			files, _ := f.Readdir(-1)
			f.Close()
			for _, fi := range files {
				runScriptIterative(filepath.Join(path, fi.Name()))
			}
		} else {
			// is a single file
			if file.Mode() & 0111 != 0 {
				// executable
				log.Printf("Executing %s:\n", path)
				out, err := exec.Command(path).Output()
				if err != nil {
					log.Print(err)
				} else {
					log.Print(out)
				}
			} else {
				log.Printf("File %s is not executable, skipping...", path)
			}
		}
	}
}

func delFileIterative(path string) {
	// first we try a remove all method
	err := os.RemoveAll(path)
	// if unable to clean up, we remove as much as we can
	if err == nil  {
		log.Printf("%s removed on first try\n", path)
	} else {
		file, err := os.Stat(path)
		if os.IsNotExist(err) {
			log.Printf("%s not found or insufficient privilege, skipping...\n", path)
		} else {
			if file.IsDir() {
				log.Printf("Entering directory %s\n", path)
				f, _ := os.Open(path)
				files, _ := f.Readdir(-1)
				f.Close()
				for _, fi := range files {
					delFileIterative(filepath.Join(path, fi.Name()))
				}
			}
			os.Remove(path)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func execute(conf *config) {
	log.Println("Executing hooks...")

	// execute hooks
	for _, entry := range conf.ExecuteScripts {
		runScriptIterative(entry)
	}

	// delete files
	for _, entry := range conf.DeleteFiles {
		delFileIterative(entry)
	}

	if conf.ExitAfterTrigger {
		log.Println("Job done, RIP")
		os.Exit(0)
	}
}