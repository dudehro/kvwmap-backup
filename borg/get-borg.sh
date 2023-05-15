#!/bin/bash
wget borg https://github.com/borgbackup/borg/releases/download/1.2.3/borg-linux64
mv borg-linux64 /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
