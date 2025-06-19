# use in n/vim to restart on save:
# :autocmd BufWritePost * silent! !./autoload.sh
#!/bin/bash
pkill lakecitydisruptor || true
go build -o lakecitydisruptor
echo http://localhost:11028
./lakecitydisruptor >> log.txt 2>&1 &
