/*
This is a simple search query parser.

Search query is usually get from a browser's input box, and used for searching in a database. After receive the search query, we need to do some basic process on it and transform it to a structured variable. This package is for doing the job.

Syntax of search query:

    1. All queries have a key and one or more values.
    2. Key and values are separate by a colon (":"). And there's no space before and after the colon. Eg. key:value
    3. If a key or value contains punctuations or spaces, they should be placed in quotation marks. Eg. from:"one@example.com"
    4. Single quotes could placed within double quotes, and double quotes could placed within single quotes. Eg. "'key'":'"value"'
    5. One key could have more than one values which are separate by comma. Eg. key:value1,value2,"value3"
    6. Put a minus before key means negative. Eg. -name:"not this"
    7. If a query does not contains a key, it's key is supposed to space. Eg. "a value not contains a key"

After process by function Parse, the search query will be transformd to a type of "Nodes" variable.

Examples:
    key:value                   -> &[{key [value] false}]
    key1:value1 key2:value2     -> &[{key1 [value1] false} {key2 [value2] false}]
    key:v1,v2                   -> &[{key [v1 v2] false}]
    "key has space":v           -> &[{key has space [v] false}]
    key:"value's space"         -> &[{key [value's space] false}]
    -key:"negative value"       -> &[{key [negative value] true}]
    "only value"                -> &[{ [only value] false}]
    two values                  -> &[{ [two] false} { [values] false}]
    k1:v1 k1:v1,v2              -> &[{k1 [v1] false} {k1 [v1 v2] false}]
*/
package queryparser
