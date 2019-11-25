#!/usr/bin/env bash

screen=${1:-0}

# wait 20s for the display manager service to start and attach to screen 
sleep 20

/usr/bin/xset -display :$screen s off 		# deactivate screen saver
/usr/bin/xset -display :$screen -dpms 		# disable DPMS
/usr/bin/xset -display :$screen s noblank 	# disable screen blanking

