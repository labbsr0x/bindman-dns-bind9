# Sandman Bind9 Manager

This repository defines the Sandman component that manages a Bind9 instance through nsupdates dispatched from a Sandman DNS Listener.

A listener is responsible for calling a manager whenever DNS Binding updates are identified. 

The REST APIs contracts are defined by the [Sandman DNS Webhook](https://github.com/labbsr0x/sandman-dns-webhook) project

# Configuration

The manager is setup with the help of environment variables and volume mapping in the following way: 

## Volume Mapping

The manager needs to keep a store of records being managed. Hence, a `/data` volume must be mapped to the host. There we also expect to find the `.private` keyfile for communication with the actual `nameserver`

## Environment variables

1. `mandatory` **SANDMAN_NAMESERVER_ADDRESS**: address of the nameserver that a instance of a Sandman Bind9 Manager will manage

2. `mandatory` **SANDMAN_NAMESERVER_KEYFILE**: the keyfile name that will be used to authenticate with the nameserver. **MUST** match the regexp `K.*\.\+157\.\+.*\.key` and **MUST** be inside the `/data` volume

3. `optional` **SANDMAN_NAMESERVER_PORT**: custom port for communication with the nameserver; defaults to `53`

4. `optional` **SANDMAN_DNS_TTL**: the dns recording rule expiration time (or time-to-live). By default, the TTL is **3600 seconds**.

5. `optional` **SANDMAN_DNS_REMOVAL_DELAY**: the delay in minutes to be applied to the removal of an DNS entry. The default is 10 minutes. This is to guarantee that in fact the removal should be processed.

