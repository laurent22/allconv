# allconv

A command line tool to convert between number bases and currencies.

## Installation

    go build -o aconv main.go
    
## Usage

    Usage: aconv <command> [flags] [<value>]

    Commands:
       list          Lists all the possible conversions.
       <from>2<to>   Converts from <from> to <to>. eg. hex2bin, dec2oct, eur2usd, etc.
       help          Displays this help page.

    Flags:
       --format   Output format - either "simple", "withUnit" or "full". (Default: full)

    Examples:
       aconv bin2hex 1100110010   # Convert binary to hexadecimal
       aconv hex2dec ff5c         # Convert hexadecimal to decimal
       aconv eur2usd 10           # Convert Euros to US Dollars
       aconv aud2jpy 5000         # Convert Australian Dollars to Japanese Yens

## License

http://opensource.org/licenses/MIT
