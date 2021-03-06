package main

// TODO: Add unit?

import (
	"flag"
	"./conversions"
	"strings"
	"errors"
	"os"
	"fmt"
)

func parseConversionCommand(cmd string) (string, string, error) {
	tokens := strings.Split(cmd, "2")
	if len(tokens) != 2 {
		return "", "", errors.New("Not a conversion command: \"" + cmd + "\"")
	}
	return tokens[0], tokens[1], nil
}

func printFlags() {
	longestName := 0
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "help" { return }
		if len(f.Name) > longestName {
			longestName = len(f.Name)		
		} 
	})
	
	indentOffset := longestName + 3
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "help" { return }
		s :=  "   --%s"
		s += strings.Repeat(" ", indentOffset - len(f.Name))
		s +=  "%s (Default: %s)\n"
		fmt.Printf(s, f.Name, f.Usage, f.DefValue)
	})
}

func printUsage() {
	fmt.Println("Usage: aconv [flags] <command> [<value>]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("   list          Lists all the possible conversions.")
	fmt.Println("   <from>2<to>   Converts from <from> to <to>. eg. hex2bin, dec2oct, eur2usd, etc.")
	fmt.Println("   help          Displays this help page.")
	fmt.Println("")
	fmt.Println("Flags:")
	printFlags()
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("   aconv bin2hex 1100110010   # Convert binary to hexadecimal")
	fmt.Println("   aconv hex2dec ff5c         # Convert hexadecimal to decimal")
	fmt.Println("   aconv eur2usd 10           # Convert Euros to US Dollars")
	fmt.Println("   aconv aud2jpy 5000         # Convert Australian Dollars to Japanese Yens")
}

func createFormat(formatType string) (string, error) {
	f := strings.ToLower(formatType)
	if f == "simple" {
		return "%i", nil
	}
	if f == "withunit" {
		return "%i %u", nil
	}
	if f == "full" {
		return "%i %u = %o %v", nil
	}
	return "", errors.New("Unknown format type: \"" + formatType + "\"")
}

func exitWithError(message string) {
	fmt.Println(message)
	fmt.Println("")
	printUsage()
	os.Exit(1)
}

func main() {
	var fFormat string
	var fReverse bool
	var fHelp bool
	flag.StringVar(&fFormat, "format", "full", "Output format - either \"simple\", \"withUnit\" or \"full\".")
	flag.BoolVar(&fReverse, "reverse", false, "Reverse the conversion. eg. hex2bin becomes bin2hex, etc.")
	flag.BoolVar(&fHelp, "help", false, "Displays this help page.") // Defined only to prevent Go from outputting the default help page
	flag.Parse()
	
	args := flag.Args()
	
	if len(args) < 1 {
		exitWithError("No command specified.")
	}
	
	conv := conversions.NewConversions()
	
	command := strings.ToLower(args[0])
	
	switch command {
		
		case "help": 
		
			printUsage()
			os.Exit(0)
		
		case "list": 
		
			s := ""
			categoryNames := conv.CategoryNames()
			for _, categoryName := range categoryNames {
				if s != "" {
					s += "\n"
				}
				s += strings.Title(conv.NiceCategoryName(categoryName)) + "\n"
				unitNames := conv.UnitNames(categoryName)
				for _, unitName := range unitNames {
					s += "   " + unitName + "   " + conv.NiceUnitName(categoryName, unitName) + "\n"
				}
			}
			fmt.Println(s)
			os.Exit(0)
			
		default: 
		
			if len(args) < 2 {
				exitWithError("No value specified.")
			}
		
			fromUnit, toUnit, err := parseConversionCommand(args[0])
			if err != nil {
				exitWithError(fmt.Sprint(err))
			}
			if fReverse {
				temp := fromUnit
				fromUnit = toUnit
				toUnit = temp
			}
			value := args[1]
			
			format, err := createFormat(fFormat)
			if err != nil {
				exitWithError(fmt.Sprint(err))
			}
			result, err := conv.ConvertFormat(format, fromUnit, toUnit, value)
			if err != nil {
				exitWithError("Could not convert input: " + fmt.Sprint(err))
			}
			
			fmt.Println(result)
			os.Exit(0)
			
	}
	
	panic("Unreachable")
}