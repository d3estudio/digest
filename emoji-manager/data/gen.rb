require 'net/http'
require 'uri'
require 'json'

data = Net::HTTP.get(URI.parse('https://raw.githubusercontent.com/iamcal/emoji-data/master/emoji.json'))
data = JSON.parse(data)

def map_unified(item)
    item.split('-').map { |i| [i.hex].pack('U') }.join('')
end

items = {}
data.each do |raw_emoji|
    unless items.include? raw_emoji['unified']
        repr = raw_emoji['unified']
        repr = raw_emoji['variations'][0] if raw_emoji['variations'].length == 1
        items[raw_emoji['unified']] = {
            emoji: "#{map_unified(repr)}",
            aliases: []
        }
    end
    raw_emoji['short_names'].each { |a| items[raw_emoji['unified']][:aliases].push(a) }
end

result = items.values()

File.write(File.expand_path('../db', __FILE__), JSON.dump(result))
