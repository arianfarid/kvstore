require "socket"

SOCKET_PATH = "/tmp/kvstore.sock"
ERROR = "ERROR"
NOT_FOUND = "NOTFOUND"
OK = "OK"
VALUE = "VALUE"
PATH = "./"
VALUES_NAME = "keystore_values.kv"
COMMANDS = {
  "GET" => "GET",
  "PUT" => "PUT",
  "DELETE" => "DELETE",
  "LIST" => "LIST",
}
TOMBSTONE = "<TOMBSTONE>"

class Opts
  property file_path
  property values_name

  def initialize(@file_path : String = PATH, @values_name : String = VALUES_NAME)
  end

  def get_values_path()
    self.file_path + self.values_name
  end
end

class Store
  property kv_file : File
  property index : Hash(String, String) = {} of String => String
  def initialize(path : String)
     @kv_file = File.open(path, "a+")
     @kv_file.rewind
     kv_file.each_line do |line|
      key, value = line.rstrip.split('\t', 2)
      next unless value
      if value == TOMBSTONE
        self.index.delete(key)
      else
        self.index[key] = value
      end
     end
     @kv_file.sync = true
  end
  def write_string(str : String)
    self.kv_file.puts str
  end
  def put(key : String, value : String)
    self.write_string("#{key}\t#{value}")
    self.index[key] = value
  end
  def get(key : String)
    index[key]?
  end
  def delete(key : String)
    self.write_string("#{key}\t#{TOMBSTONE}")
    self.index.delete(key)
  end
  def list()
    self.index.keys.join(" ")
  end
end

class KVServer
  def initialize(@store : Store, @socket_path : String)
    File.delete(@socket_path) if File.exists?(@socket_path)
  end

  def handle_client(client : UNIXSocket)
    while line = client.gets
      response = handle_command(line)
      client.puts response
    end
  rescue ex
    client.puts "#{ERROR} #{ex.message}"
  ensure
    client.close
  end

  def run
    server = UNIXServer.new(@socket_path)
    puts "Listening on #{@socket_path}"

    loop do
      client = server.accept
      spawn handle_client(client)
    end
  end

  def handle_command(line : String) : String
      parts = line.rstrip.split(' ', 3)
      cmd = parts[0]?

      case cmd
      when COMMANDS["GET"]
        key = parts[1]?
        return ERROR unless key

        if value = @store.get(key)
          "#{VALUE} #{value}"
        else
          NOT_FOUND
        end

      when COMMANDS["PUT"]
        key = parts[1]?
        value = parts[2]?
        return ERROR unless key && value

        @store.put(key, value)
        OK

      when COMMANDS["DELETE"]
        key = parts[1]?
        return ERROR unless key

        if @store.get(key)
          @store.delete(key)
          OK
        else
          NOT_FOUND
        end

      when COMMANDS["LIST"]
        @store.list()

      else
        ERROR
      end
    end
end

opts = Opts.new
store = Store.new(opts.get_values_path())

server = KVServer.new(store, SOCKET_PATH)
server.run