# Doublehop
>This tool aims to assist in using PowerShell Remoting / WinRM to execute PowerShell commands on remote hosts through WinRM double hop. 
> 
> This tool assumes you have internal network access and with access to a CLI.
> 
>Inspired by: ([Slayerlabs - Kerberos Double-Hop Workarounds]([linkurl](https://posts.slayerlabs.com/double-hop/))). 


## Usage
```
doublehop Usage():
    -c string
        Command we are executing.
  -j string
        Host we are executing command against.
  -l string
        Inital host we are jumping from.
  -m string
        Add hosts comma seperated. FORMAT: 'host1,host2'
  -p string
        Password for user
  -u string
```
***

## Example
```
# Execute against single host:
doublehop.exe -u "acme.local\pwneduser" -p "MySecurePass" -l wks01.acme.local -j server1.acme.local -c ipconfig

# Execute against multiple jump systems:
doublehop.exe -u "acme.local\pwneduser" -p "MySecurePass" -l wks01.acme.local -m "server1.acme.local,server2.amce.local" -c ipconfig
```