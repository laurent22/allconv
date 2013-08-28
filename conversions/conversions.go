package conversions

import (
	"errors"
	"strconv"
	"strings"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type ConsolePrinter interface {
	Println(a ...interface{}) (int, error)
	Print(a ...interface{}) (int, error)
	Printf(format string, a ...interface{}) (int, error)
}

type Conversion struct {
	category string
	from string
	to string
	convert func(input string) (string, error)
}

type Conversions struct {
	inner []Conversion
	printer ConsolePrinter
	currencies [][]string
}

// {lhs: "1 British pound",rhs: "9.2661276 Chinese yuan",error: "",icc: true}

type GoogleCalculatorResponse struct {
	Lhs string `json:"lhs"`
	Rhs string `json:"rhs"`
	Error string `json:"error"`
	Icc bool `json:"icc"`
}

var googleCalculatorRateCache map[string]float64

func NewConversions(printer ConsolePrinter) *Conversions {
	output := new(Conversions)
	output.printer = printer
	output.currencies = [][]string{
		[]string{"AED", "UAE Dirham"},
		[]string{"ANG", "Netherlands Antillean Guilder"},
		[]string{"ARS", "Argentine Peso"},
		[]string{"AUD", "Australian Dollar"},
		[]string{"BGN", "Bulgarian Lev"},
		[]string{"BHD", "Bahraini Dinar"},
		[]string{"BND", "Brunei Dollar"},
		[]string{"BOB", "Boliviano"},
		[]string{"BRL", "Brazilian Real"},
		[]string{"BWP", "Pula"},
		[]string{"CAD", "Canadian Dollar"},
		[]string{"CHF", "Swiss Franc"},
		[]string{"CLP", "Chilean Peso"},
		[]string{"CNY", "Yuan Renminbi"},
		[]string{"COP", "Colombian Peso"},
		[]string{"CRC", "Costa Rican Colon"},
		[]string{"CZK", "Czech Koruna"},
		[]string{"DKK", "Danish Krone"},
		[]string{"DOP", "Dominican Peso"},
		[]string{"DZD", "Algerian Dinar"},
		[]string{"EGP", "Egyptian Pound"},
		[]string{"EUR", "Euro"},
		[]string{"FJD", "Fiji Dollar"},
		[]string{"GBP", "Pound Sterling"},
		[]string{"HKD", "Hong Kong Dollar"},
		[]string{"HNL", "Lempira"},
		[]string{"HRK", "Croatian Kuna"},
		[]string{"HUF", "Forint"},
		[]string{"IDR", "Rupiah"},
		[]string{"ILS", "New Israeli Sheqel"},
		[]string{"INR", "Indian Rupee"},
		[]string{"JMD", "Jamaican Dollar"},
		[]string{"JOD", "Jordanian Dinar"},
		[]string{"JPY", "Yen"},
		[]string{"KES", "Kenyan Shilling"},
		[]string{"KRW", "Won"},
		[]string{"KWD", "Kuwaiti Dinar"},
		[]string{"KYD", "Cayman Islands Dollar"},
		[]string{"KZT", "Tenge"},
		[]string{"LBP", "Lebanese Pound"},
		[]string{"LKR", "Sri Lanka Rupee"},
		[]string{"LTL", "Lithuanian Litas"},
		[]string{"LVL", "Latvian Lats"},
		[]string{"MAD", "Moroccan Dirham"},
		[]string{"MDL", "Moldovan Leu"},
		[]string{"MKD", "Denar"},
		[]string{"MUR", "Mauritius Rupee"},
		[]string{"MXN", "Mexican Peso"},
		[]string{"MXV", "Mexican Unidad de Inversion (UDI)"},
		[]string{"MYR", "Malaysian Ringgit"},
		[]string{"NAD", "Namibia Dollar"},
		[]string{"NGN", "Naira"},
		[]string{"NIO", "Cordoba Oro"},
		[]string{"NOK", "Norwegian Krone"},
		[]string{"NPR", "Nepalese Rupee"},
		[]string{"NZD", "New Zealand Dollar"},
		[]string{"OMR", "Rial Omani"},
		[]string{"PEN", "Nuevo Sol"},
		[]string{"PGK", "Kina"},
		[]string{"PHP", "Philippine Peso"},
		[]string{"PKR", "Pakistan Rupee"},
		[]string{"PLN", "Zloty"},
		[]string{"PYG", "Guarani"},
		[]string{"QAR", "Qatari Rial"},
		[]string{"RON", "New Romanian Leu"},
		[]string{"RSD", "Serbian Dinar"},
		[]string{"RUB", "Russian Ruble"},
		[]string{"SAR", "Saudi Riyal"},
		[]string{"SCR", "Seychelles Rupee"},
		[]string{"SEK", "Swedish Krona"},
		[]string{"SGD", "Singapore Dollar"},
		[]string{"SLL", "Leone"},
		[]string{"SVC", "El Salvador Colon"},
		[]string{"THB", "Baht"},
		[]string{"TND", "Tunisian Dinar"},
		[]string{"TRY", "Turkish Lira"},
		[]string{"TTD", "Trinidad and Tobago Dollar"},
		[]string{"TWD", "New Taiwan Dollar"},
		[]string{"TZS", "Tanzanian Shilling"},
		[]string{"UAH", "Hryvnia"},
		[]string{"UGX", "Uganda Shilling"},
		[]string{"USD", "US Dollar"},
		[]string{"UYU", "Peso Uruguayo"},
		[]string{"UZS", "Uzbekistan Sum"},
		[]string{"VEF", "Bolivar Fuerte"},
		[]string{"VND", "Dong"},
		[]string{"YER", "Yemeni Rial"},
		[]string{"ZAR", "Rand"},
		[]string{"ZMK", "Zambian Kwacha"},
	}
	
	baseConv := func(input string, inputBase int, format string) (string, error) {
		n, err := strconv.ParseUint(input, inputBase, 64)
		if err != nil { return "", err }
		return fmt.Sprintf(format, n), nil
	}
	
	output.Add(Conversion{
		"number", "dec", "hex", func(input string) (string, error) {
			return baseConv(input, 10, "%x")
		},
	})
	
	output.Add(Conversion{
		"number", "dec", "bin", func(input string) (string, error) {
			return baseConv(input, 10, "%b")
		},
	})
	
	output.Add(Conversion{
		"number", "dec", "oct", func(input string) (string, error) {
			return baseConv(input, 10, "%o")
		},
	})
	
	output.Add(Conversion{
		"number", "hex", "dec", func(input string) (string, error) {
			return baseConv(input, 16, "%d")
		},
	})
	
	output.Add(Conversion{
		"number", "hex", "bin", func(input string) (string, error) {
			return baseConv(input, 16, "%b")
		},
	})
	
	output.Add(Conversion{
		"number", "hex", "oct", func(input string) (string, error) {
			return baseConv(input, 16, "%o")
		},
	})
	
	output.Add(Conversion{
		"number", "bin", "dec", func(input string) (string, error) {
			return baseConv(input, 2, "%d")
		},
	})
	
	output.Add(Conversion{
		"number", "bin", "hex", func(input string) (string, error) {
			return baseConv(input, 2, "%x")
		},
	})
	
	output.Add(Conversion{
		"number", "bin", "oct", func(input string) (string, error) {
			return baseConv(input, 2, "%o")
		},
	})
	
	output.Add(Conversion{
		"number", "oct", "dec", func(input string) (string, error) {
			return baseConv(input, 8, "%d")
		},
	})
	
	output.Add(Conversion{
		"number", "oct", "hex", func(input string) (string, error) {
			return baseConv(input, 8, "%x")
		},
	})
	
	output.Add(Conversion{
		"number", "oct", "bin", func(input string) (string, error) {
			return baseConv(input, 8, "%b")
		},
	})
	
	currencyConv := func(input string, from string, to string) (string, error) {
		if googleCalculatorRateCache == nil { googleCalculatorRateCache = make(map[string]float64) }
		cacheKey := from + "_" + to
		floatInput, err := strconv.ParseFloat(input, 64)
		if err != nil { return "", err }
		output, ok := googleCalculatorRateCache[cacheKey]
		if ok { 
			s := strconv.FormatFloat(output * floatInput, 'f', 8, 64)
			return s, nil
		}
		gcUrl := "http://www.google.com/ig/calculator?hl=en&q=1" + strings.ToUpper(from) + "%3D%3F" + strings.ToUpper(to)
		if printer != nil { printer.Println("Getting latest rates from " + gcUrl + " ...") }
		resp, err := http.Get(gcUrl)
		if err != nil { return "", err }
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil { return "", err }
		jsonString := string(body)
		jsonString = strings.Replace(jsonString, "lhs:", "\"lhs\":", -1)
		jsonString = strings.Replace(jsonString, "rhs:", "\"rhs\":", -1)
		jsonString = strings.Replace(jsonString, "error:", "\"error\":", -1)
		jsonString = strings.Replace(jsonString, "icc:", "\"icc\":", -1)
		var googleResp GoogleCalculatorResponse
		err = json.Unmarshal([]byte(jsonString), &googleResp)
		if googleResp.Error != "" { return "", errors.New("Google Calculator error: " + googleResp.Error) }
		rhs := strings.Split(googleResp.Rhs, " ")
		if len(rhs) <= 0 { return "", errors.New("Invalid response format: " + googleResp.Rhs) }
		conv, err := strconv.ParseFloat(rhs[0], 64)
		if err != nil { return "", err }
		googleCalculatorRateCache[cacheKey] = conv
		return strconv.FormatFloat(conv * floatInput, 'f', 8, 64), nil
	}
	
	for i := 0; i < len(output.currencies); i++ {
		for j := 0; j < len(output.currencies); j++ {
			c1 := strings.ToLower(output.currencies[i][0])
			c2 := strings.ToLower(output.currencies[j][0])
			output.Add(Conversion{
				"currency", c1, c2, func(input string) (string, error) {
					return currencyConv(input, c1, c2)
				},
			})
		}
	}
	
	return output
}

func (this *Conversions) Convert(from string, to string, input string) (string, error) {
	for _, c := range this.inner {
		if c.from == from && c.to == to {
			return c.convert(input)
		}
	}
	return "", errors.New("Unsupported conversion: \"" + from + "\" to \"" + to + "\"") 
}

func (this *Conversions) Add(c Conversion) {
	this.inner = append(this.inner, c)
}

func (this *Conversions) NiceCategoryName(s string) string {
	return s
}

func (this *Conversions) NiceUnitName(category string, s string) string {
	s = strings.ToLower(s)

	if category == "number" {
		if s == "hex" { return "Hexadecimal" }
		if s == "bin" { return "Binary" }
		if s == "dec" { return "Decimal" }
		if s == "oct" { return "Octal" }
	}
	
	if category == "currency" {
		for _, row := range this.currencies {
			if strings.ToLower(row[0]) == s {
				return row[1]
			}
		}
	}
	
	return s
}

func (this *Conversions) CategoryNames() []string {
	var output []string 
	for _, c := range this.inner {
		found := false
		for _, n := range output {
			if c.category == n {
				found = true
				break
			}
		}
		if !found {
			output = append(output, c.category)
		}		
	}
	return output
}

func (this *Conversions) UnitNames(category string) []string {
	var output []string 
	for _, c := range this.inner {
		if c.category != category {
			continue
		}
		found := false
		for _, n := range output {
			if c.from == n {
				found = true
				break
			}
		}
		if !found {
			output = append(output, c.from)
		}		
	}
	return output
}