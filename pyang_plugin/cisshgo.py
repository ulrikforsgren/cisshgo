"""cisshgo output plugin

This plugin takes a YANG data model and produces a YAML formatted context
search map to be loaded by the cisshgo network device simulator.
The YANG model must contain tailf-common:cli* statements to be able to produce
a working  context tree.
"""

import json
import pprint as pp
import sys

from pyang import plugin, error, types, statements
from pyang.util import unique_prefixes

pprint = pp.PrettyPrinter(indent=4).pprint

def pyang_plugin_init():
    plugin.register_plugin(CisshgoPlugin())


hid = 3

def replace_in_regexp(s, f, r):
  while (i:=s.find(f)) != -1:
    lb = s.rfind(r'[', 0, i)
    rb = s.rfind(r']', 0, i)
    if lb<=rb:
      s = s.replace(f, f'[{r}]', 1)
    else:
      s = s.replace(f, r, 1)
  return s

def replace_unicode_patterns(p):
  if r'\p{N}' in p:
    p = replace_in_regexp(p, r'\p{N}', r'0-9')
  if r'\p{L}' in p:
    p = replace_in_regexp(p, r'\p{L}', r'a-zA-Z')
  return p

class CisshgoPlugin(plugin.PyangPlugin):
    def add_output_format(self, fmts):
        self.multiple_modules = True
        fmts['cisshgo'] = self

    def setup_fmt(self, ctx):
        ctx.implicit_errors = False
        ctx.identifier_state = self
        self.nodes = {}

    def emit(self, ctx, modules, fd):
        """Main control function.
        """
        for epos, etag, eargs in ctx.errors:
            if error.is_error(error.err_level(etag)):
                if etag != 'EXTENSION_ARGUMENT_PRESENT':
                    raise error.EmitError("CISSHGO plugin needs a valid module")
        tree = {}
        mods = {}
        annots = {}
        for m,p in unique_prefixes(ctx).items():
            mods[m.i_modulename] = [p, m.search_one("namespace").arg]
        for module in modules:
            for ann in module.search(("ietf-yang-metadata", "annotation")):
                typ = ann.search_one("type")
                annots[module.arg + ":" + ann.arg] = (
                    "string" if typ is None else self.base_type(ann, typ))
        for module in modules:
            self.process_children(module, 3)


    ### TODO: Refactor for better structure ( only proccess children...)
    def process_children(self, node, hparent, commands=None):
        """Process all children of `node`, except "rpc" and "notification".
        """
        global hid
        cli_mode_name = None
        cli_exit_command = None
        cli_allow_join_with_key = False
        cli_drop_node_name = False
        commands = commands or []

        for st in node.substmts:
            if type(st.keyword) is tuple:
                m, k = st.keyword
                if m == "tailf-common":
                    if k == "cli-mode-name":
                        # TODO: Default is to use container/list
                        cli_mode_name = st.arg
                    elif k == "cli-exit-command":
                        cli_exit_command = st.arg
                    elif k == "cli-allow-join-with-key":
                        cli_allow_join_with_key = True
                    elif k == "cli-drop-node-name":
                        cli_drop_node_name = True
        """
        Other statements to_consider:
         - tailf:cli-suppress-mode
         - tailf:cli-add-mode

         - tailf:cli-flatten-container

        """

        if node.keyword == 'container':
          assert(cli_allow_join_with_key == False)
          if not cli_drop_node_name:
            commands.append(node.arg)
        elif node.keyword == 'list':
          if not cli_drop_node_name:
            join = "<J>" if cli_allow_join_with_key else ""
            commands.append(f"{join}{node.arg}")
          commands.append(self.key_data(node))
        if cli_mode_name:
          commands = [ self.regexp_command(c) for c in commands ]
          commands = [ c for c in self.flatten_list(commands) ]
          self.print_entry(commands, hparent, cli_mode_name, cli_exit_command)
          hparent = hid
          commands = []
        for ch in node.i_children:
            if ch.keyword in ["choice", "case"]:
                self.process_children(ch, hparent, commands.copy())
            elif ch.keyword == "container":
                self.process_children(ch, hparent, commands.copy())
            elif ch.keyword == "list":
                self.process_children(ch, hparent, commands.copy())

    def print_entry(self, commands, hparent, cli_mode_name, cli_exit_command):
        global hid
        hid += 1
        print("        -")
        print(f"            cmd: {' '.join(commands)}")
        print(f"            id: {hid}")
        print(f"            up: {hparent}")
        print(f"            mode: \"({cli_mode_name})#\"")
        hparent = hid
        if cli_exit_command:
          print(f"            exit-cmd: \"{cli_exit_command}\"")

    def key_data(self, node):
      key = []
      for k in node.i_key:
          assert(type(k) is statements.LeafLeaflistStatement)
          t = k.search_one("type")
          kp = self.type_data(t)
          key.append(kp)
      return key

    def flatten_list(self, lst):
      for i in lst:
        if type(i) is list:
          for j in i:
            yield j
        else:
          yield i

    def regexp_command(self, key):
        if type(key) is str: return key
        return [ f'<R>{self.regexp_string(l)}' for l in key ]

    def regexp_string(self, t):
        kp, p2 = t
        if kp == 'string':
            if not p2:
                return '.+'
            return '|'.join(p2)
        if kp[:4] == 'uint':
            return '[0-9]+'
        if kp[:3] == 'int':
            return '-{0,1}[0-9]+'
        if kp == 'union':
            r = [ self.regexp_string(p) for p in p2 ]
            return '|'.join(r)
        if kp == 'enumeration':
            #TODO: escape special regexp characters needed.
            return '|'.join(p2)
        print(f"ERROR: Unhandled data type {kp} {p2}")
        sys.exit(2)

    def type_data(self, t, union=False, level=0):
        ts = t.i_type_spec
        tst = type(ts)

        #print("   "*level, t, ts.name, tst)
        if tst is types.IntTypeSpec:
            return (ts.name, '')
        if tst is types.RangeTypeSpec:
            return (ts.name, '')
        elif tst is types.StringTypeSpec:
            return (ts.name, '')
        elif tst is types.LengthTypeSpec: # Strings with only a length restriction
            return (ts.name, '')
        elif tst is types.PatternTypeSpec: # Strings with patterns and optional length
            return (ts.name, [ replace_unicode_patterns(str(p)) for p in ts.res ])
        elif tst is types.EnumTypeSpec:
            return (ts.name, [e for e,_ in ts.enums])
        elif tst is types.PathTypeSpec:
            return (ts.name, '')
        elif tst is types.UnionTypeSpec:
            ta = []
            if t.i_typedef:
                t = t.i_typedef.search_one('type')
            subtypes = t.search("type")
            assert(len(subtypes) or t.i_typedef)
            for st in subtypes:
                ta.append(self.type_data(st, level=level+1))
            return (ts.name, ta)
        else:
            raise TypeError(f"Can't handle type: {t.arg} {tst}")


    def base_type(self, ch, of_type):
        """Return the base type of `of_type`."""
        while 1:
            if of_type.arg == "leafref":
                if of_type.i_module.i_version == "1":
                    node = of_type.i_type_spec.i_target_node
                else:
                    node = ch.i_leafref.i_target_node
            elif of_type.i_typedef is None:
                break
            else:
                node = of_type.i_typedef
            of_type = node.search_one("type")
        if of_type.arg == "decimal64":
            return [of_type.arg, int(of_type.search_one("fraction-digits").arg)]
        elif of_type.arg == "union":
            return [of_type.arg, [self.base_type(ch, x) for x in of_type.i_type_spec.types]]
        else:
            return of_type.arg
