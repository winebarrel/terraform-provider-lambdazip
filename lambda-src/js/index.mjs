#!/usr/bin/env node
import figlet from "figlet";

export const handler = async (_event, _context) => {
  const data = await figlet("Hello, lambdazip!");
  return data;
};

if (import.meta.url === `file://${process.argv[1]}`) {
  console.log(await handler());
}
