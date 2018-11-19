# Sandman Bind9 Manager

This repository defines the Sandman component that manages a Bind9 instance through nsupdates dispatched from a Sandman DNS Listener.

A listener is responsible for calling a manager whenever DNS Binding updates are identified. 

The REST APIs contracts are defined by the [Sandman DNS Webhook](https://github.com/labbsr0x/sandman-dns-webhook) project