# use in n/vim to restart on save:
# :autocmd BufWritePost * silent! !./autoload.sh
#!/bin/bash
pkill lakecityd || true
go build -o lakecitydisruptor.com
echo http://localhost:11028
./lakecitydisruptor.com >> log.txt 2>&1 &
