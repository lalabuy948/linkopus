# bundle config set path 'vendor/bundle'
# bundle install

require 'sinatra'
require 'json'
require 'logger'

log = Logger.new('pipeline.log')

def verify_signature(payload_body)
    signature = 'sha1=' + OpenSSL::HMAC.hexdigest(OpenSSL::Digest.new('sha1'), ENV['SECRET_TOKEN'], payload_body)
    return halt 500, "Signatures didn't match!" unless Rack::Utils.secure_compare(signature, request.env['HTTP_X_HUB_SIGNATURE'])
end

post '/webhook' do
    payload_body = request.body.read
    verify_signature(payload_body)

    push = JSON.parse(payload_body)
    log.debug "webhook recieved... commit [#{push["head_commit"]["id"]}]"
    if push["ref"] == "refs/heads/master"
        log.debug "starting pipeline..."
        pid = spawn("./pipeline.sh")
        Process.detach(pid)
    end
end
