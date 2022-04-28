pyang-plugin:
 - Use YANG dependencies from NSO: Find arguments for pyang
 - Handle extensions: Ignore or just filter...

transcript:
 - Think about concept of vendor, platform. Get rid of vendor?
 - Rename csr1000v -> ios

options:
 - Debug/Trace

Debug/Trace:
 - Function to log timestamp, node + arb. text.

Read transcript:
 - Validation of transcript: existance, format, ...
 - Simplify reading. Strange list/map
 - Search context [1] must always exits: Add checks

Native commands:
 - Add support for terminal width command.

Command parsing:
 - Cisco NX (at least) accepts semi-colon (;) separated commands:
   "terminal length 0 ; terminal width 511 ; show version ; show inventory"
