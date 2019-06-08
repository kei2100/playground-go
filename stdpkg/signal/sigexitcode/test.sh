#!/bin/bash

cd $(cd $(dirname $0) && pwd)

go build -o sleep

./sleep
EXIT_CODE=$?

echo "送信したsignalとgoプログラムの終了デフォルト終了コード"
echo "=== NORMAL_EXIT:${EXIT_CODE}"

for SIG in "2:INT" "15:TERM" "9:KILL"; do
    ./sleep &
    PID=$!
    NUM=$(echo $SIG | cut -f 1 -d :)
    kill -$NUM $PID
    wait $PID
    EXIT_CODE=$?
    echo "=== $(echo $SIG | cut -f 2 -d :):$EXIT_CODE"
done

#  $ ./test.sh
#
#  送信したsignalとgoプログラムの終了デフォルト終了コード
#  === NORMAL_EXIT:0
#  === INT:0
#  ./test.sh: line 18: 56781 Terminated: 15          ./sleep
#  === TERM:143
#  ./test.sh: line 18: 56788 Killed: 9               ./sleep
#  === KILL:137


# unixにおいて、Exit Code 128 以降は、128 + {Kill Signal}
#
# 0〜127	アプリケーションの終了コード
#
# 128	発生しない
# 143	「kill -15」⇒ 128+15 = 143
# 137	「kill -9」⇒ 128+9 = 137
