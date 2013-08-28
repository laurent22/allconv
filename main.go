package main

// TODO: Tool to loop through available currencies
// TODO: Format currency numbers
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

func printUsage() {
	fmt.Println("Usage: aconv <command> [<value>]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("   list          Lists all the possible conversions.")
	fmt.Println("   <from>2<to>   Converts from <from> to <to>. eg. hex2bin, dec2oct, eur2usd, etc.")
	fmt.Println("   help          Displays this help page.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("   aconv bin2hex 1100110010   # Convert binary to hexadecimal")
	fmt.Println("   aconv hex2dec ff5c         # Convert hexadecimal to decimal")
	fmt.Println("   aconv eur2usd 10           # Convert Euros to US Dollars")
	fmt.Println("   aconv aud2jpy 5000         # Convert Australian Dollars to Japanese Yens")
	flag.PrintDefaults()
}

func processCommand(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("No command specified.")
	}
	
	conv := conversions.NewConversions(nil)
	
	command := strings.ToLower(args[0])
	
	switch command {
		
		case "help": 
		
			printUsage()
			return "", nil
		
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
			return s, nil
			
		default: 
		
			if len(args) < 2 {
				return "", errors.New("No value specified.")
			}
		
			fromUnit, toUnit, err := parseConversionCommand(args[0])
			if err != nil {
				return "", err
			}
			value := args[1]
							
			result, err := conv.Convert(fromUnit, toUnit, value)
			if err != nil {
				return "", errors.New("Could not convert input: " + fmt.Sprint(err))
			}
			
			return result, nil
			
	}	
}

func main() {	
	flag.Parse()
	
	result, err := processCommand(flag.Args())
	if err != nil {
		fmt.Println(err)
		fmt.Println("")
		printUsage()
		os.Exit(1)
	}
	
	if result != "" {
		fmt.Println(result)
	}
}