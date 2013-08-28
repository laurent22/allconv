# allconv

A command line tool to convert between number bases and currencies.

## Installation

    go build
    
## Usage

    Usage: ac <command> [<value>]
    
    Commands:
       list          Lists all the possible conversions.
       <from>2<to>   Converts from <from> to <to>. eg. hex2bin, dec2oct, eur2usd, etc.
       help          Displays this help page.
    
    Examples:
       ac bin2hex 1100110010   # Convert binary to hexadecimal
       ac hex2dec ff5c         # Convert hexadecimal to decimal
       ac eur2usd 10           # Convert Euros to US Dollars
       ac aud2jpy 5000         # Convert Australian Dollars to Japanese Yens
