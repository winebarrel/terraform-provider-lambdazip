#!/usr/bin/env node
const util = require("util");
const figlet = util.promisify(require("figlet"));

exports.handler = async (event, context) => {
  const data = await figlet("Hello, lambdazip!");
  console.log(data);
};

if (require?.main === module) {
  exports.handler();
}
