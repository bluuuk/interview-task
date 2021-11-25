# interview-task

> First write a token generator that creates a file with 10 million random tokens, one per line, each consisting of seven lowercase letters a-z. Then write a token reader that reads the file and stores the tokens in your DB. Naturally some tokens will occur more than once, so take care that these aren't duplicated in the DB, but do produce a list of all non-unique tokens and their frequencies. Find a clever way to do it efficiently in terms of network I/O, memory, and time and include documentation inline with your code or as txt file, describing your design decisions.

## Back of the Envelope Calculation

- $10 \cdot 10^6 \cdot 7 = 7 \cdot 10^7$ total bytes, so roughly `66` MB. Should fit into main memory
- $|\{a,\cdots,z\}| = 26 \implies 26^7$ combinations. Thinking of the Birthday problem, there is a good chance that we will have a collision of two random strings.

## Implementation thoughts

### Go

1. For a simple implementation, a line by line read in should do the job.
  - Further io/performance improvement with `bufio`, but needs testing
2. Checking for duplicates

  > $n = 10^7$

  - A `Map[string]int` can do the job with low programming complexity, because access is O(1) and filling it is O(n)
    - We can even hint `make` with $10^7$ to reduce the number of map resizing
  - Radix sort with 7 iterations will be O(n * 7)
    - But should be very memory intensiv, because we have to copy $10^7$ strings around in every bucket   
  - Couting sort will be O(n + 26^7) which does not look promosing because of the possible huge input space

  $\implies$ https://www.youtube.com/watch?v=kVgy1GSDHG8
  
### Database
  
> I have never used Postgres, so I will give it a try

The lib (PGX)[https://github.com/jackc/pgx] seems to be up to date, so I'll use it.

#### Schema

We just have to save tokens. (Documentation)[https://www.postgresql.org/docs/9.3/datatype-character.html] and (SO)[https://dba.stackexchange.com/questions/126003/index-performance-for-char-vs-varchar-postgres] suggest using `character varying` to save the token. We do not need to mark the token with `PRIMARY KEY` because it implies `UNIQUE` which we already check beforehand. In addition, filtering duplicates reduces the calls to insert rows into the database.

```sql
CREATE TABLE TOKENS(
  token VARCAHR(7) NOT NULL
)
```


