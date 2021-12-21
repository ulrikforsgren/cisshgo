"""cisshgo output plugin

This plugin takes a YANG data model and produces a YAML formatted context
search map to be loaded by the cisshgo network device simulator.
The YANG model must contain tailf-common:cli* statements to be able to produce
a working  context tree.
"""

import json
import pprint as pp

from pyang import plugin, error, types, statements
from pyang.util import unique_prefixes

pprint = pp.PrettyPrinter(indent=4).pprint

def pyang_plugin_init():
    plugin.register_plugin(CisshgoPlugin())

class CisshgoPlugin(plugin.PyangPlugin):
    def add_output_format(self, fmts):
        self.multiple_modules = True
        fmts['cisshgo'] = self

    def setup_fmt(self, ctx):
        ctx.implicit_errors = False

    def emit(self, ctx, modules, fd):
        """Main control function.
        """
        for epos, etag, eargs in ctx.errors:
            if error.is_error(error.err_level(etag)):
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
            self.process_children(module)
        #pprint(tree)

    def process_children(self, node, indent=0, commands=None):
        """Process all children of `node`, except "rpc" and "notification".
        """
        cli_mode_name = None
        if indent == 0:
          print("(config)#")
          indent+= 4
          commands = []
        for st in node.substmts:
            if type(st.keyword) is tuple:
                m, k = st.keyword
                if m == "tailf-common":
                    if k == "cli-mode-name":
                        cli_mode_name = st.arg

        if node.keyword == 'list':
          self.key_data(node, commands)
        if cli_mode_name:
          modes = f"{' '*indent}{cli_mode_name}"
          print(f"{modes:<70} {commands}")
          indent += 4
          commands = []
        for ch in node.i_children:
            if ch.keyword in ["rpc", "notification"]:
                continue
            if ch.keyword in ["choice", "case"]:
                self.process_children(ch, indent, commands.copy())
                continue
            if ch.keyword == "container":
                cmds = commands.copy()
                cmds.append(ch.arg)
                self.process_children(ch, indent, cmds)
            elif ch.keyword == "list":
                cmds = commands.copy()
                cmds.append(ch.arg)
                self.process_children(ch, indent, cmds)
            elif ch.keyword in ["leaf", "leaf-list"]:
                continue

    def key_data(self, node, commands):
      key = []
      for k in node.i_key:
          #print(k, type(k))
          assert(type(k) is statements.LeafLeaflistStatement)
          kp = self.type_data(k)
          key.append(kp)
      commands.append(f"{{{','.join(key)}}}")

    def type_data(self, k):
      tp = []
      for x in k.substmts:
          kp = ""
          st = type(x)
          if st is statements.Statement:
            pass
            # Ignore extension statements, e.g.
            # tail:info "XYZ";
          elif st is statements.TypeStatement:
            ts = type(x.i_type_spec)
            t = x.i_type_spec.name
            if ts is types.RangeTypeSpec:
              assert(t in ["int8", "int16", "int32", "int64",
                          "uint8", "uint16", "uint32", "uint64"])
              #print("RANGE:", x.i_ranges)
              # Not required
              # Only Regexp for generic numerical
              # pattern is required
            elif ts is types.PatternTypeSpec:
              assert(t == "string")
              #print("PATTERN:", x.i_type_spec.res)
              # Should be compiled to Regexp
              # Most important when combined with key
              # Is the len always 1?
              kp = str(x.i_type_spec.res[0])
            elif ts is types.EnumTypeSpec:
              assert(t == "enumeration")
              #print("ENUMS:", x.i_type_spec.enums)
              # Can easily be compiled into a Regexp.
              # First implementation can accept
              # anything.
              kp = "/".join([e for e,_ in x.i_type_spec.enums])
              #kp = x.i_type_spec.enums
            elif ts is types.StringTypeSpec:
              assert(t == "string")
            elif ts is types.UnionTypeSpec:
              # Investigate how to combine to a single
              # Regexp? or create two regexps
              #kp = kp or "union"
              kp = kp or self.type_data(x)

            elif ts is types.IntTypeSpec:
              assert(t in ["int8", "int16", "int32", "int64",
                          "uint8", "uint16", "uint32", "uint64"])
              # Not required.
            elif ts is types.LengthTypeSpec:
              assert(t == "string")
              # Not required.
            else:
              print(f"UNKNOWN TYPE: {ts}")
            kp = kp or t
            tp.append(kp)
      return ' or '.join(tp)

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
