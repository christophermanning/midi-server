#!/bin/bash
set -euo pipefail

DEVICEID=20:0

# aseqdump sample output
#
#Waiting for data. Press Ctrl+C to end.
#Source  Event                  Ch  Data
# 20:0   Note on                 0, note 70, velocity 126
# 20:0   Note off                0, note 70, velocity 0

aseqdump -p $DEVICEID |
  grep --line-buffered 'Note' |
  awk -W interactive '{ printf "{\"Timestamp\":0,\"Status\":%s,\"Data1\":%s,\"Data2\":%s}\n",substr($0,14,2)=="on"?144:128, substr($0,41,2), substr($0,54,3) }'
