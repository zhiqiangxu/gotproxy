# gotproxy, transparent proxy in go!


**gotproxy** aims to provide global transparent proxy service, now it only works for mach osx, but support for other systems is planed!

On mach osx, it works this way:

1. all tcp traffic is redirected to a user space programe listening on a specific port, by kernel extension
2. the user space programe queries the original destination and do proxy as needed


So that we can control global outbound traffic from a single point!
