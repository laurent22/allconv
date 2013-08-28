package main

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"	
	"log"
	"fmt"
)

type GoogleCalculatorResponse struct {
	Lhs string `json:"lhs"`
	Rhs string `json:"rhs"`
	Error string `json:"error"`
	Icc bool `json:"icc"`
}

func main() {
	allCurrencies := [][]string{
		[]string{"AED", "UAE Dirham"},
		[]string{"AFN", "Afghani"},
		[]string{"ALL", "Lek"},
		[]string{"AMD", "Armenian Dram"},
		[]string{"ANG", "Netherlands Antillean Guilder"},
		[]string{"AOA", "Kwanza"},
		[]string{"ARS", "Argentine Peso"},
		[]string{"AUD", "Australian Dollar"},
		[]string{"AWG", "Aruban Florin"},
		[]string{"AZN", "Azerbaijanian Manat"},
		[]string{"BAM", "Convertible Mark"},
		[]string{"BBD", "Barbados Dollar"},
		[]string{"BDT", "Taka"},
		[]string{"BGN", "Bulgarian Lev"},
		[]string{"BHD", "Bahraini Dinar"},
		[]string{"BIF", "Burundi Franc"},
		[]string{"BMD", "Bermudian Dollar"},
		[]string{"BND", "Brunei Dollar"},
		[]string{"BOB", "Boliviano"},
		[]string{"BOV", "Mvdol"},
		[]string{"BRL", "Brazilian Real"},
		[]string{"BSD", "Bahamian Dollar"},
		[]string{"BTN", "Ngultrum"},
		[]string{"BWP", "Pula"},
		[]string{"BYR", "Belarussian Ruble"},
		[]string{"BZD", "Belize Dollar"},
		[]string{"CAD", "Canadian Dollar"},
		[]string{"CDF", "Congolese Franc"},
		[]string{"CHE", "WIR Euro"},
		[]string{"CHF", "Swiss Franc"},
		[]string{"CHW", "WIR Franc"},
		[]string{"CLF", "Unidades de fomento"},
		[]string{"CLP", "Chilean Peso"},
		[]string{"CNY", "Yuan Renminbi"},
		[]string{"COP", "Colombian Peso"},
		[]string{"COU", "Unidad de Valor Real"},
		[]string{"CRC", "Costa Rican Colon"},
		[]string{"CUC", "Peso Convertible"},
		[]string{"CUP", "Cuban Peso"},
		[]string{"CVE", "Cape Verde Escudo"},
		[]string{"CZK", "Czech Koruna"},
		[]string{"DJF", "Djibouti Franc"},
		[]string{"DKK", "Danish Krone"},
		[]string{"DOP", "Dominican Peso"},
		[]string{"DZD", "Algerian Dinar"},
		[]string{"EGP", "Egyptian Pound"},
		[]string{"ERN", "Nakfa"},
		[]string{"ETB", "Ethiopian Birr"},
		[]string{"EUR", "Euro"},
		[]string{"FJD", "Fiji Dollar"},
		[]string{"FKP", "Falkland Islands Pound"},
		[]string{"GBP", "Pound Sterling"},
		[]string{"GEL", "Lari"},
		[]string{"GHS", "Ghana Cedi"},
		[]string{"GIP", "Gibraltar Pound"},
		[]string{"GMD", "Dalasi"},
		[]string{"GNF", "Guinea Franc"},
		[]string{"GTQ", "Quetzal"},
		[]string{"GYD", "Guyana Dollar"},
		[]string{"HKD", "Hong Kong Dollar"},
		[]string{"HNL", "Lempira"},
		[]string{"HRK", "Croatian Kuna"},
		[]string{"HTG", "Gourde"},
		[]string{"HUF", "Forint"},
		[]string{"IDR", "Rupiah"},
		[]string{"ILS", "New Israeli Sheqel"},
		[]string{"INR", "Indian Rupee"},
		[]string{"IQD", "Iraqi Dinar"},
		[]string{"IRR", "Iranian Rial"},
		[]string{"ISK", "Iceland Krona"},
		[]string{"JMD", "Jamaican Dollar"},
		[]string{"JOD", "Jordanian Dinar"},
		[]string{"JPY", "Yen"},
		[]string{"KES", "Kenyan Shilling"},
		[]string{"KGS", "Som"},
		[]string{"KHR", "Riel"},
		[]string{"KMF", "Comoro Franc"},
		[]string{"KPW", "North Korean Won"},
		[]string{"KRW", "Won"},
		[]string{"KWD", "Kuwaiti Dinar"},
		[]string{"KYD", "Cayman Islands Dollar"},
		[]string{"KZT", "Tenge"},
		[]string{"LAK", "Kip"},
		[]string{"LBP", "Lebanese Pound"},
		[]string{"LKR", "Sri Lanka Rupee"},
		[]string{"LRD", "Liberian Dollar"},
		[]string{"LSL", "Loti"},
		[]string{"LTL", "Lithuanian Litas"},
		[]string{"LVL", "Latvian Lats"},
		[]string{"LYD", "Libyan Dinar"},
		[]string{"MAD", "Moroccan Dirham"},
		[]string{"MDL", "Moldovan Leu"},
		[]string{"MGA", "Malagasy Ariary"},
		[]string{"MKD", "Denar"},
		[]string{"MMK", "Kyat"},
		[]string{"MNT", "Tugrik"},
		[]string{"MOP", "Pataca"},
		[]string{"MRO", "Ouguiya"},
		[]string{"MUR", "Mauritius Rupee"},
		[]string{"MVR", "Rufiyaa"},
		[]string{"MWK", "Kwacha"},
		[]string{"MXN", "Mexican Peso"},
		[]string{"MXV", "Mexican Unidad de Inversion (UDI)"},
		[]string{"MYR", "Malaysian Ringgit"},
		[]string{"MZN", "Mozambique Metical"},
		[]string{"NAD", "Namibia Dollar"},
		[]string{"NGN", "Naira"},
		[]string{"NIO", "Cordoba Oro"},
		[]string{"NOK", "Norwegian Krone"},
		[]string{"NPR", "Nepalese Rupee"},
		[]string{"NZD", "New Zealand Dollar"},
		[]string{"OMR", "Rial Omani"},
		[]string{"PAB", "Balboa"},
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
		[]string{"RWF", "Rwanda Franc"},
		[]string{"SAR", "Saudi Riyal"},
		[]string{"SBD", "Solomon Islands Dollar"},
		[]string{"SCR", "Seychelles Rupee"},
		[]string{"SDG", "Sudanese Pound"},
		[]string{"SEK", "Swedish Krona"},
		[]string{"SGD", "Singapore Dollar"},
		[]string{"SHP", "Saint Helena Pound"},
		[]string{"SLL", "Leone"},
		[]string{"SOS", "Somali Shilling"},
		[]string{"SRD", "Surinam Dollar"},
		[]string{"SSP", "South Sudanese Pound"},
		[]string{"STD", "Dobra"},
		[]string{"SVC", "El Salvador Colon"},
		[]string{"SYP", "Syrian Pound"},
		[]string{"SZL", "Lilangeni"},
		[]string{"THB", "Baht"},
		[]string{"TJS", "Somoni"},
		[]string{"TMT", "Turkmenistan New Manat"},
		[]string{"TND", "Tunisian Dinar"},
		[]string{"TOP", "Pa’anga"},
		[]string{"TRY", "Turkish Lira"},
		[]string{"TTD", "Trinidad and Tobago Dollar"},
		[]string{"TWD", "New Taiwan Dollar"},
		[]string{"TZS", "Tanzanian Shilling"},
		[]string{"UAH", "Hryvnia"},
		[]string{"UGX", "Uganda Shilling"},
		[]string{"USD", "US Dollar"},
		[]string{"UYI", "Uruguay Peso en Unidades Indexadas (URUIURUI)"},
		[]string{"UYU", "Peso Uruguayo"},
		[]string{"UZS", "Uzbekistan Sum"},
		[]string{"VEF", "Bolivar Fuerte"},
		[]string{"VND", "Dong"},
		[]string{"VUV", "Vatu"},
		[]string{"WST", "Tala"},
		[]string{"XAF", "CFA Franc BEAC"},
		[]string{"XAG", "Silver"},
		[]string{"XAU", "Gold"},
		[]string{"XCD", "East Caribbean Dollar"},
		[]string{"XFU", "UIC-Franc"},
		[]string{"XOF", "CFA Franc BCEAO"},
		[]string{"XPD", "Palladium"},
		[]string{"XPF", "CFP Franc"},
		[]string{"XPT", "Platinum"},
		[]string{"XSU", "Sucre"},
		[]string{"YER", "Yemeni Rial"},
		[]string{"ZAR", "Rand"},
		[]string{"ZMK", "Zambian Kwacha"},
		[]string{"ZWL", "Zimbabwe Dollar"},
	}	
	
	var goodOnes [][]string
	for _, c := range allCurrencies {
		gcUrl := "http://www.google.com/ig/calculator?hl=en&q=1" + strings.ToUpper(c[0]) + "%3D%3F" + strings.ToUpper(c[0])
		resp, err := http.Get(gcUrl)
		if err != nil { panic(err) }
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil { panic(err) }
		jsonString := string(body)
		jsonString = strings.Replace(jsonString, "lhs:", "\"lhs\":", -1)
		jsonString = strings.Replace(jsonString, "rhs:", "\"rhs\":", -1)
		jsonString = strings.Replace(jsonString, "error:", "\"error\":", -1)
		jsonString = strings.Replace(jsonString, "icc:", "\"icc\":", -1)
		var googleResp GoogleCalculatorResponse
		err = json.Unmarshal([]byte(jsonString), &googleResp)
		if googleResp.Error != "" && googleResp.Error != "0" {
			log.Println("Not good: " + c[0])
		} else {
			log.Println("OK: " + c[0])
			goodOnes = append(goodOnes, c)
		}
	}
	
	fmt.Println("=======================================")
	fmt.Println("[][]string{")
	for _, c := range goodOnes {
		fmt.Printf("	[]string{\"%s\", \"%s\"},\n", c[0], c[1])
	}
	fmt.Println("}")
}