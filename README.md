![Logo of the project](https://github.com/haynesherway/haynesbot/blob/master/logo.png?raw=true)

# Haynes Bot

A discord bot that does pokemon go related things, like calculate IVs and display moves and effects.

Calculations and pokemon data come from [haynesherway/pogo](https://www.github.com/haynesherway/pogo)

To add the bot to your server:  https://discordapp.com/oauth2/authorize?client_id=402854185072328714&scope=bot
**NOTE: Currently, if you add the bot to your server, it needs a restart to be able to send messages, so please send me a message here or in the discord below so I can restart it for you. If this becomes a more common issue, I will spend the time fixing it, but currently it is only added to around 1 server a week.**

HaynesBot help discord: https://discord.gg/CakVND

I made a patreon because people said they wanted to donate to the project, but it isn't necessary.
https://www.patreon.com/haynesbot

Update 4/5/2018: The !raidiv command now also works with research reward encounters (level 15)

Update 10/29/2018: !luckydate command added

[![GoDoc](https://godoc.org/github.com/haynesherway/haynesbot?status.svg)](https://godoc.org/github.com/haynesherway/haynesbot) [![Build Status](https://travis-ci.org/haynesherway/haynesbot.svg?branch=master)](https://travis-ci.org/haynesherway/haynesbot)


## Commands

* **!cp** {pokemon} {level} {attack iv} {defense iv} {stamina iv}  
		Get CP of a pokemon at a specified level with specified IVs  
		Example: !cp mewtwo 25 15 14 15  
* **!maxcp** {pokemon}  
		Get maximum CP of a pokemon with perfect IVs at level 40  
		Example: !maxcp latios  
* **!raidiv** {pokemon}  
		Get range of possible raid CPs for specified pokemon  
		Example: !raidcp groudon  
* **!raidiv** {pokemon} {cp}  
		Get possible IV combinations for specified raid (or research reward encounter) pokemon with specified IV  
		Example: !raidcp kyogre 2292  
* **!mewiv** {cp}  
		Get possible IV combinations for mew 
		Example: !mewiv 1306		
* **!raidchart** {pokemon}  
		Get a chart with possible stats for specified pokemon at raid level above 90%  
		Example: !raidchart machamp  
* **!moves** {pokemon}
		Get a list of fast and charge moves for specified pokemon  
		Example: !moves rayquaza  
* **!type** {pokemon}  
		Get a list of types for a specified pokemon  
		Example: !type rayquaza  
* **!effect** {pokemon|type}  
		Get a list of type relations a specified pokemon or type has  
		Example: !effect pikachu or !effect electric  
* **!luckydate**
		Returns the date for pokemon to have been caught by for a higher change at luckies.
		Example: !luckydate
* **!normal**
		Returns an image of the normal version of the pokemon.
		Example: !normal pidgey
* **!shiny**
		Returns an image of the shiny version of the pokemon.
		Example: !shiny pidgey
		
## Server Owner Commands:

* **!add**  
		Add server management capabilities
* **!add teams**  
		Add pokemon go team management
* **!setprefix** {prefix}  
		Set prefix other than ! for your server  
		Example: !setprefix $
* **!setwelcome** {message}  
		Set welcome message for when new members join your server  
		You can mention user with {mention}, print username with {user} and print server name with {guild}  
		Requires a channel #welcome  
		Example: !setwelcome Welcome to {guild}, {mention}!  
* **!setgoodbye** {message}  
		Set goodbye message for when users leave your guild  
		You can mention user with {mention}, print username with {user} and print server name with {guild}  
		Requires a channel #welcome  
		Example: !setgoodbye Goodbye {user}, we won't miss you!  


## Configuration

You will need to put your discord bot token in the config.json file

## Examples

!wat

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/wat.png?raw=true" width="500" height="700" title="Wat">

!wat iv

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/watdetail.png?raw=true" width="500" height="354" title="WatIV">

!moves

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/moves.png?raw=true" width="500" height="235" title="IV">

!effect pokemon

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/effect.png?raw=true" width="500" height="285" title="IV">

!effect type

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/effecttype.png?raw=true" width="500" height="282" title="IV">

!raidiv (!raidcp)

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/raidiv.png?raw=true" width="300" height="290" title="RaidIV">

!iv

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/iv.png?raw=true" width="300" title="IV">

!maxcp

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/maxcp.png?raw=true" width="300" height="236" title="IV">

!type

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/type.png?raw=true" width="300" height="225" title="IV">

!normal

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/normal.png?raw=true" width="300" title="Normal">

!shiny 

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/shiny.png?raw=true" width="300" title="Shiny">

!raidchart

![Raidchart Mewtwo-A](https://github.com/haynesherway/haynesbot/blob/master/examples/RAIDCHART-mewtwo-a.png?raw=true) ![Raidchart Deoxys](https://github.com/haynesherway/haynesbot/blob/master/examples/RAIDCHART-deoxys.png?raw=true) ![Raidchart Jirachi](https://github.com/haynesherway/haynesbot/blob/master/examples/RAIDCHART-jirachi.png?raw=true)

