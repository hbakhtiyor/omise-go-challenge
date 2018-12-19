package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/omise/omise-go"
)

var (
	OmisePublicKey  = os.Getenv("OP_KEY")
	OmiseSecretKey  = os.Getenv("OS_KEY")
	DefaultCurrency = "thb"
	options         *Options
)

func main() {
	programName := os.Args[0]
	commaFlag := flag.String("comma", ",", "The field `delimiter` for CSV files.")
	flag.Parse()

	if *commaFlag == "" {
		// flag parser not trigering for empty flags, e.g. -comma=
		fmt.Println("flag needs an argument: -comma")
		flag.Usage()
		return
	}

	options = &Options{
		Comma:   []rune(*commaFlag)[0],
		Headers: []string{"Name", "AmountSubunits", "CCNumber", "CVV", "ExpMonth", "ExpYear"},
	}

	if OmisePublicKey == "" {
		fmt.Printf("%v: set OP_KEY enviroment\n", programName)
		return
	}
	if OmiseSecretKey == "" {
		fmt.Printf("%v: set OS_KEY enviroment\n", programName)
		return
	}

	// use pipe streaming if data available
	if info, err := os.Stdin.Stat(); err != nil {
		fmt.Printf("%v: %v\n", programName, err)
		return
	} else if info.Mode()&os.ModeNamedPipe != 0 {
		// use empty string for indicating pipe mode
		if err := run([]string{""}); err != nil {
			fmt.Printf("%v: %v\n", programName, err)
		}
		return
	}

	if flag.NArg() < 1 {
		fmt.Printf("Usage: %v [file] [url] ...\n", programName)
		flag.PrintDefaults()
		return
	}

	if err := run(flag.Args()); err != nil {
		fmt.Printf("%v: %v\n", programName, err)
		return
	}
}

func run(files []string) error {
	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if err != nil {
		return err
	}

	for _, fileName := range files {
		options.Extension = filepath.Ext(fileName)

		if options.Extension == ".tsv" {
			options.Comma = '\t'
		}

		reader, err := getReader(fileName)
		if err != nil {
			fmt.Printf("%v: %v\n", fileName, err)
			continue
		}

		chPayment := make(chan *Payment)

		go func(reader io.ReadCloser) {
			if err := processCSV(reader, chPayment, options); err != nil {
				fmt.Printf("%v: %v\n", fileName, err)
			}
		}(reader)

		result := &DonationResult{TopPayments: []*Payment{{}, {}, {}}, Currency: DefaultCurrency}

		fmt.Println("performing donations...")

		var wg sync.WaitGroup
		cpus := runtime.NumCPU()
		for payment := range chPayment {
			result.sortTopPayments(payment)
			result.TotalPayments++
			wg.Add(1)
			go func(payment *Payment) {
				defer wg.Done()

				if succeed, err := createTokenAndCharge(client, payment); err != nil {
					log.Fatal(err)
				} else if succeed {
					result.SuccessSum += payment.Amount
				} else {
					result.FailSum += payment.Amount
				}
			}(payment)
			if result.TotalPayments%cpus == 0 {
				wg.Wait()
				result.printSummary(true)
			}
		}

		wg.Wait()

		result.cleanSummary()
		fmt.Printf("done.\n")
		result.printSummary(false)
	}

	return nil
}
