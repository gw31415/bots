#!/usr/bin/env python
from sys import stdin, stdout
from proto import cmdin_pb2, cmdout_pb2, help_pb2
from sympy import latex, parse_expr
stdin = stdin.buffer.read()
in_pb = cmdin_pb2.Input()
in_pb.ParseFromString(stdin)
out_msg = cmdout_pb2.BotMsg()
for media in in_pb.media:
    buf = cmdout_pb2.OutputMedia()
    if media.type == 2:
        buf.type = 2
        res = str(latex(parse_expr(media.data.decode(encoding='utf-8'))))
        buf.data = res.encode(encoding='utf-8')
    else:
        buf.type == 2
        buf.data = "invalid type".encode(encoding='utf-8')
        buf.error = 1
    out_msg.medias.append(buf)
out_pb = cmdout_pb2.Output()
out_pb.msgs.append(out_msg)
stdout.buffer.write(out_pb.SerializeToString())
