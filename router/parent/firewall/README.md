

# README.rulesets for doxi / dogtown-naxi-rules

- Readme-Version: 2014-04-04
- [latest ruleset-commits](https://bitbucket.org/lazy_dogtown/doxi-rules/src)
- [Doxi-News Blog](http://blog.dorvakt.org/)

these rulesets are now available as independent git-repo @ 
[bitbucket.org/lazy_dogtown/doxi-rules](https://bitbucket.org/lazy_dogtown/doxi-rules)

for tools to manage your doxi-rules you might want to install doxi-tools
[bitbucket.org/lazy_dogtown/doxi](https://bitbucket.org/lazy_dogtown/doxi)

to keep track of changes and ruleset-updates you could either 
subscribe to the [doxi-news - blog](http://blog.dorvakt.org/) ([rss-feed](http://blog.dorvakt.org/feeds/posts/default)), 
subscribe to the naxsi-mailinglist 
https://groups.google.com/forum/?fromgroups#!forum/naxsi-discuss or
subscribe to the [ruleset-commit-feed](https://bitbucket.org/lazy_dogtown/doxi-rules/rss)
or follow that project on Bitbucket

License: see License.txt



all not-mentioned files here are part of naxsi/nginx - default-configuration


# configuration rules 

please note: due to changes in naxsi after 0.49 this file-layout might get 
obsolete. 

### rules.conf

- your global includes-file; you might setup different rules.con - files,
- maybe tuned for each virtualhost.


### learning-mode.rules 

- rules to configure/enable learning-mode 
    
### active-mode.rules 

- rules to configure active-mode (block)


# detection rules

### app_server.rules

- rules you might want to enable when running nginx as lb/proxy 
for app-servers like tomcat / rails etc and you're shure to
have no php/asp/cgi - files lying around

### malware.rules

**NOTE: for a better coverage you might want to try a real ids
like snort or suricata  with et-rulesets rules to detect malicious
content in- and outbound. **
    
- this ruleset is designed to detect malicious request that give a 
hint for hacked / misused / C&C-servers and tries to detect
web-backdoors, webshells and other malicious access to unwanted
files/services.
    
- **CAUTION:** these rules are quite noise, so if included you might want to
tune and create whitelists for your applications
    
### scanner.rules
    
- detect scanners (WebAppScanners/Testing-Tools
- detetc vuln-scanning-bots or attack-tools) by UA or by certain requests.
- some of these rules could be included into web_[app|server].rules,
like scanners for certain webapp/server-vulns, but when there's a 
clear sign for an automated scanning-process the sigs are include here
- **CAUTION:** these rules are quite noise, so if included you might want to
tune and create whitelists for your applications
    

### web_app.rules

- detect exploit/misuse-attempts againts web-applications; please see 
scanner.rules for some details on webapp-based scanners

### web_server.rules
    
- generic rules to protect a webserver from misconfiguration 
and known mistakes / exploit-vectors 


# misc. rules (obsolete, not maintained after jan 2014)

# misc_whitelisting.rules 

- whitelistings for different webapps/actions that are known to fail
on certain parameters 
    
