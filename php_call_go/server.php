<?php
try {
    @unlink('/tmp/keyword_match.sock');
    $server = socket_create(AF_UNIX, SOCK_STREAM, 0);
    socket_set_option($server, SOL_SOCKET, SO_REUSEADDR, 1);
    socket_bind($server, '/tmp/keyword_match.sock');
    socket_listen($server);
} catch (\Exception $e) {
    var_dump(error_get_last());
    exit;
}
$socket_list = array(
    $server
);
$sockets = array($server);

while (true) {
    $socket_list = $sockets;
    $result = socket_select($socket_list, $write, $except, NULL);
    if ($result === false) {
        exit('error');
    }

    foreach ($socket_list as $key => $read) {
        // 服务器socket可读，处理连接
        if ($read == $server) {
            $client = socket_accept($read);
            $sockets[] = $client;
            continue;
        }

        $data = socket_read($read, 1024);
        // 读到空即客户端断开连接，从监听列表中去除此socket
        if (!$data) {
            foreach ($sockets as $k => $socket) {
                if ($socket == $read) {
                    unset($sockets[$k]);
                }
            }
            continue;
        }

        socket_write($read, $data, strlen($data));
    }

    sleep(1);
}