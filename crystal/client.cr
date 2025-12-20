require "socket"

SOCKET_PATH = "/tmp/kvstore.sock"

if ARGV.empty?
    STDERR.puts "Usage: kv <COMMAND> [key] [value]"
  exit 1
end

command = ARGV.join(" ")

begin
    socket = UNIXSocket.new(SOCKET_PATH)
    socket.puts command

    if response = socket.gets
        puts response.rstrip
    end
ensure
    socket.close if socket
end
