#! /usr/bin/env node
var figlet = require("figlet");

figlet("Hello", function (_err, data) {
  console.log(data);
});
