pyang-plugin:
 - Use YANG dependencies from NSO: Find arguments for pyang
 - Handle extensions: Ignore or just filter...

transcript:
 - Think about concept of vendor, platform. Get rid of vendor?
 - Rename csr1000v -> ios

options:
 - Debug/Trace

bad code:
 - Fix variables for numListeners and startintPort
   numListeners is actually upper port number.
   Do not return pointer to int.

Debug/Trace:
 - Function to log timestamp, node + arb. text.

Read transcript:
 - Validation of transcript: existance, format, ...
 - Simplify reading. Strange list/map
 - Search context [1] must always exits: Add checks

Command parsing:
 - Cisco NX (at least) accepts semi-colon (;) separated commands:
   "terminal length 0 ; terminal width 511 ; show version ; show inventory"

Modules:
 - Change version of golang.org/x/term to unique local version to avoid
   conflict.

Wish list:
 - Writable + compare config
 - 
