# Ragnarok Online Packet Analyser
Program for analysing incoming packets for Ragnarok Online (iRO-Chaos). 
This application is still under development.

## Install Guide
Clone source code from repository, compile it and run.  
Application will try to dial to destination ```tcp://localhost:10101``` and accept incoming traffic from this connection.  
Incoming traffic should send to this application sniffed packets from ragnarok client. 


### Functionality:
1. Reads changes of character
1. Reads changes of homunculus
1. Detects when user / homunculus gains experience
1. Detects and notifies when homunculus hunger becomes lower than 25
1. Detects when you are summoning homunculus
1. Detects homunculus intimacy level
1. ETA until levelUP (base / job)
1. Exp per hour
1. Amount of monsters killed

### TODO:
1. Move settings to command line parameters
1. Picked up items
1. DPS Meter
1. Track all dropped items
1. Track all met characters
1. Zeny per hour
1. Dropped items from monster
1. Fill game object and keep all information about character in web UI
