# Dead Man's Switch

This little daemon watchs a DNS record periodically, if it contains a pre-defined value or is missing, triggers a set of actions.

Currently supported triggers:
* Execute programs or scripts
* Delete files

## Installation

There are precompiled binaries for Linux in the [releases](https://github.com/Jamesits/DeadManSwitch/releases) page.

Requires:
* Go 1.10 or later (only needed for compilation)
* Linux or Windows (other OSes are not tested; on Windows you need to figure out how to install yourself)

There is currently no formal package. A `install.sh` can be used to compile from source and install to your localhost, and a `package.sh` can be used to generate a binary tarball.

## Config

By default (I mean, if you use the systemd service provided) the config is at `/etc/dmswitch/config.toml`, and the scripts or programs placed in `/etc/dmswitch/hooks` will be run once triggered. 

The config file is self-explanatory. 

* Records are analyzed using substrings, so:
  * **DO NOT** use your domain or record type as a trigger
  * **DO NOT** make the trigger string a substring of the normal string, or vice versa
* All relative paths in the config is relative to the config file itself. 
* Programs will be executed in alphabet order. 
* File deletion happens after program execution.

## Usage

There is a systemd service installed by default. You can use `systemctl enable --now dmswitch` to start it on boot.

If you prefer launching the binary directly, use `-conf path/to/config/file` to point it to the config file. If this parameter is missing, it searches for a `config.toml` in its working directory.

## Design

DNS is a perfect one-way channel for C&C. It doesn't require the program to connect to a specific server (recursion works most of the time), there are many free providers available, and there are APIs everywhere.
