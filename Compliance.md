# APIKit compliance with the OpenAPIv2 specification

Only functional features are listed. Purely descriptive items (e.g. the info object) are directly copied to the output.

| Section | Feature | Scoped  | Implemented | Testcase |
| --- |---|:---:|:---:|:---:|
| Data Types | type: format |
|| integer: int32 | x |
|| integer: int64 | x |
|| number: float | x |
|| number: double | x |
|| string | x |
|| string: byte | x |
|| string: binary | x |
|| boolean | x |
|| string: date | x |
|| string: date-time | x |
|| string: password | x |
| Schemes || x
|| http / https | x
|| ws / wss | -
| Client: Consumes | (request to server) |
||  application/json | x
||  application/hal+json | x
||  application/x-www-form-urlencoded | x
||  multipart/form-data | x
| Client: Produces | (response from server) |
||  application/json | x
||  application/hal+json | x
||  application\/xml | x
||  file types as stream | x
| Server: Consumes | (request from client) |
||  application/json | x
||  application/x-www-form-urlencoded | x
||  multipart/form-data | x
| Server: Produces | (response to client) |
||  application/json | x
||  application/hal+json | x
||  file types as stream | x
| Definitions | (see Schema) | x
| Paths || x
|| $ref | x
|| GET, PUT, POST, DELETE, OPTIONS, HEAD, PATCH | x
|| Parameters (see Parameter) | x
| Operations || x
|| Individual Consumes | x
|| Individual Produces | x
|| Parameters (see Parameter) | x
|| Responses (see Response) | x
|| Schemes | -
|| Security (see Security Definitions) | x
| Parameters || x
|| in Path (all types) | x
|| in Query | x
|| in Header | x
|| Body | x
|| Required attribute | x
|| Allow empty value | x
| Parameters and Items || x
|| $ref | x
|| CSV array | x
|| SSV array | -
|| TSV array | -
|| Pipes array | -
|| Multi array (in Form and Query) | ?
|| Maps | x
|| Default values | x
|| Maximum | x
|| Exclusive maximum | x
|| Minimum | x
|| Exclusive minimum | x
|| Maximimum length | x
|| Minimum length | x
|| Pattern | x
|| Maximium items | x
|| Minimum items | x
|| Unique items | x
|| Enums | x
|| Multiple of | x
| Global parameter definitions | | x
| Response || x
|| $ref | x
|| Schema (see Schema) | x
|| Headers | ?
| Schema || x
|| $ref | x
|| All data types and formats | x
|| Default value | x
|| Maximum | x
|| Exclusive maximum | x
|| Minimum | x
|| Exclusive minimum | x
|| Maximimum length | x
|| Minimum length | x
|| Pattern | x
|| Maximium items | x
|| Minimum items | x
|| Unique items | x
|| Enums | x
|| Multiple of | -
|| All of | -
|| Maximimum properties | -
|| Minimum properties | -
|| Properties | -
|| Additional properties | -
|| Discriminator | -
|| Read only | -
| Global response definitions | | x
| Security Definitions ||
|| Basic Auth | x
|| Api Key | x
|| OAuth 2 | -
| Global security definitions || x
