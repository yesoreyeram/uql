import { XMLParser, X2jOptionsOptional } from "fast-xml-parser";
import { Command, CommandResult, type_parse_arg } from "../../types";

export const parseXML = (pv: CommandResult, cv: Extract<Command, { type: "parse-xml" }>): CommandResult => {
  let output = pv.output;
  if (typeof pv.output === "string") {
    let xml_parser_options = get_parse_xml_options(cv.args);
    let parser = new XMLParser(xml_parser_options);
    output = parser.parse(pv.output);
  }
  return { ...pv, output };
};

const get_parse_xml_options = (args: type_parse_arg[][]): X2jOptionsOptional => {
  let options: X2jOptionsOptional = { ignoreAttributes: false, allowBooleanAttributes: true, commentPropName: "#comments" };
  if (args[0] && args[0].length > 0) {
    args[0].forEach((arg) => {
      switch (arg.identifier) {
        case "ignoreAttributes":
        case "allowBooleanAttributes":
        case "alwaysCreateTextNode":
        case "preserveOrder":
        case "parseTagValue":
        case "parseAttributeValue":
        case "trimValues":
          options[arg.identifier] = arg.value.toLowerCase() === "true";
          break;
        case "commentPropName":
          options["commentPropName"] = arg.value || "#comment";
          break;
        case "attributeNamePrefix":
          options["attributeNamePrefix"] = arg.value || "@_";
          break;
        default:
          // @ts-ignore
          options[arg.identifier] = arg.value.toLowerCase() === "true" ? true : arg.value.toLowerCase() === "false" ? false : arg.value;
          break;
      }
    });
  }
  return options;
};
