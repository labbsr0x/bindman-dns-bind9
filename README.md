# Bindman DNS
![Build Status](https://travis-ci.com/labbsr0x/bindman-dns-bind9.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/labbsr0x/bindman-dns-bind9)](https://goreportcard.com/report/github.com/labbsr0x/bindman-dns-bind9)

This repository defines the component that manages Bind9 DNS Server instances.

NSUpdate commands get dispatched from REST API calls defined in the bindman webhook project [Bindman DNS Webhook](https://github.com/labbsr0x/sandman-dns-webhook).

# Configuration

The bindman is setup with the help of environment variables and volume mapping in the following way: 

## Volume Mapping

A store of records being managed is needed. Hence, a `/data` volume must be mapped to the host. There, we also expect to find the `.private` and `.key` files for secure communication with the actual `nameserver`

## Environment variables

1. `mandatory` **BINDMAN_NAMESERVER_ADDRESS**: address of the nameserver that an instance of a Bindman will manage

2. `mandatory` **BINDMAN_NAMESERVER_KEYFILE**: the zone keyfile name that will be used to authenticate with the nameserver. **MUST** match the regexp `K.*\.\+157\+.*\.key` and **MUST** be inside the `/data` volume

3. `mandatory` **BINDMAN_NAMESERVER_ZONE**: the name of the zone a bindman-dns-bind9 instance is able to manage;

4. `optional` **BINDMAN_NAMESERVER_PORT**: custom port for communication with the nameserver; defaults to `53`

5. `optional` **BINDMAN_DNS_TTL**: the dns recording rule expiration time (or time-to-live). By default, the TTL is **3600 seconds**.

6. `optional` **BINDMAN_DNS_REMOVAL_DELAY**: the delay in minutes to be applied to the removal of an DNS entry. The default is 10 minutes. This is to guarantee that in fact the removal should be processed.

7. `optional` **BINDMAN_MODE**: let the runtime know if the DEBUG mode is activated; useful for debugging the intermediary files created for sending `nsupdate` commands. Possible values: `DEBUG|PROD`. Empty defaults to `PROD`.

# Secure communication

On the `/keys` folder of the `bind` service, you will find the keys that enable secure communication between the manager and the Bind9 Server for the `test.com` zone.

For now, we support only `dnssec-keygen` generated keys. We used the following commands for the `test.com` zone:

```
dnssec-keygen -a HMAC-MD5 -b 512 -n HOST test.com
```

[Go here](http://www.firewall.cx/linux-knowledgebase-tutorials/system-and-network-services/831-linux-bind-ipadd-data-file.html) to understand a bit more about how to properly configure your BIND DNS server.
