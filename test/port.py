#encoding=utf-8

from BaseHTTPServer import BaseHTTPRequestHandler
from BaseHTTPServer import HTTPServer 


class RequestHandler(BaseHTTPRequestHandler):
    # def __init__(self):

    # 基本的echo应答
    def do_GET(self):
        print ("get one message")
        # echostr = "hello world\n"
        # self.send_response(200)
        # self.send_header("Content-Length", str(len(echostr)))
        # self.end_headers()
        # self.wfile.write(echostr.encode());
        return

    def do_POST(self):
        print("post one message")

        length = self.headers.get('Content-Length')
        if length is None:
            self.send_error(400, "Empty Content")
            return
        post_data = self.rfile.read(int(length))
        self.log_message("post data %s", post_data)

        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        # self.wfile.write(json.dumps(post_data))
        return


def main():
    process = HTTPServer(('0.0.0.0', 1935), RequestHandler)
    process.serve_forever()

if __name__ == '__main__':
    main()
