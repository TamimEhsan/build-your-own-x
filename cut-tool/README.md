
## Run locally
```bash
> go run . -f=2 -d="," fourchords.csv | uniq | wc -l
  155
```
```bash
> tail -n5 fourchords.csv| go run . -d="," -f=1 -f=2 -
 "Young Volcanoes",Fall Out Boy
 "You Found Me",The Fray
 "You'll Think Of Me",Keith Urban
 "You're Not Sorry",Taylor Swift
 "Zombie",The Cranberries
```