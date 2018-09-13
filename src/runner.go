package main

import (
	"github.com/miekg/dns"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func check(conf *config) uint8{
	req := new(dns.Msg)
	req.SetQuestion(dns.Fqdn(conf.Record), dns.StringToType[conf.RecordType])

	client := new(dns.Client)
	client.SingleInflight = true
	resolvers := conf.CustomResolvers
	if conf.TrySystemResolver {
		_, err := os.Stat("/etc/resolv.conf")
		if os.IsExist(err) {
			cfg, err := dns.ClientConfigFromFile("/etc/resolv.conf")
			if err != nil {
				for _, elem := range cfg.Servers {
					resolvers = append(resolvers, elem)
				}
			}
		}
	}

	var in *dns.Msg
	var rtt time.Duration
	var err error
	gotResult := false
	for _, server := range resolvers {
		in, rtt, err = client.Exchange(req, server)
		log.Printf("Server: %s RTT: %s\n", server, rtt)

		// unexpected result
		if err != nil {
			log.Printf("Unable to resolve record using %s\n", server)
		} else {
			gotResult = true
			break
		}
	}
	if !gotResult {
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
			filenames := make([]string, len(files))
			for i, v := range files {
				filenames[i] = v.Name()
			}
			sort.Strings(filenames)
			for _, fi := range filenames {
				runScriptIterative(filepath.Join(path, fi))
			}
		} else {
			// is a single file
			log.Printf("Executing %s:\n", path)
			out, err := exec.Command(path).Output()
			if err != nil {
				log.Print(err)
			} else {
				log.Print(out)
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