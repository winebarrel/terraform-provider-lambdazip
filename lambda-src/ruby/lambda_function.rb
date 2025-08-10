#!/usr/bin/env ruby
require 'ruby_figlet'
using RubyFiglet

def handler(event:, context:)
  'meow...'.art
end

if File.basename(__FILE__) == File.basename($0)
  puts handler(event: nil, context: nil)
end
