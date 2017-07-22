<?php
$socket = socket_create(AF_UNIX, SOCK_STREAM, 0);
socket_connect($socket, '/tmp/keyword_match.sock');

$msg = "我爱北京天安门！我是个PHP和GO开发者！";
socket_send($socket, $msg, strlen($msg), 0);
$response = socket_read($socket, 1024);
var_dump($response);

