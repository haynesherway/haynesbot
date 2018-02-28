![Logo of the project](https://github.com/haynesherway/haynesbot/blob/master/logo.png?raw=true)

# Haynes Bot

A discord bot that does pokemon go related things, like calculate IVs 


## Commands

* **!cp** {pokemon} {level} {attack iv} {defense iv} {stamina iv}  
		Get CP of a pokemon at a specified level with specified IVs  
		Example: !cp mewtwo 25 15 14 15  
* **!maxcp** {pokemon}  
		Get maximum CP of a pokemon with perfect IVs at level 40  
		Example: !maxcp latios  
* **!raidcp** {pokemon}  
		Get range of possible raid CPs for specified pokemon  
		Example: !raidcp groudon  
* **!raidcp** {pokemon} {cp}  
		Get possible IV combinations for specified raid pokemon with specified IV  
		Example: !raidcp kyogre 2292  
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


## Configuration

You will need to put your discord bot token in the config.json file

## Examples

!wat

<img src="https://github.com/haynesherway/haynesbot/blob/master/examples/wat.png?raw=true" width="500" height="572" title="Wat">

![wat](https://github.com/haynesherway/haynesbot/blob/master/examples/wat.png?raw=true)

!wat iv

![wat iv](https://github.com/haynesherway/haynesbot/blob/master/examples/watdetail.png?raw=true)

!raidchart rayquaza

![Raidchart Rayquaza](https://github.com/haynesherway/haynesbot/blob/master/examples/RAIDCHART-Rayquaza.png?raw=true)
