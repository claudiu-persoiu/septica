# API

## Format
Send:
```json
{
  "action": "string",
  "data": "string"
}
```

Result:
```
{
  "success": true,
  "action": "string",
  "data": "string"
}
```

## Actions

|  Action 	|   Send   	|      Receive      	|          Description         	|
|:-------:	|:--------:	|:-----------------:	|:----------------------------:	|
| start   	| nil      	| game_key          	| Start a new game             	|
| join    	| game_key 	| wait/invalid/full 	| Join an already started game 	|
| begin   	|          	|                   	|                              	|
| joined  	|          	|                   	|                              	|
| dropped 	|          	|                   	|                              	|
| joining 	|          	|                   	|                              	|


Generated with [https://www.tablesgenerator.com/markdown_tables]()
