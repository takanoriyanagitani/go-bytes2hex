#!/bin/sh

rtm=wazero
wsm="./bytes2hex.wasm"

ex1(){
	echo helo
	printf helo | $rtm run "${wsm}"
	echo
	echo
}

ex2(){
	echo 4095 bytes

	pass='generate-insecure-random-number'

	dd if=/dev/zero bs=4095 count=1 status=none |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:"${pass}" |
		tail --bytes=16 |
		xxd -ps

	dd if=/dev/zero bs=4095 count=1 status=none |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:"${pass}" |
		$rtm run "${wsm}" |
		tail --bytes=32

	echo
	echo
}

ex3(){
	echo 1048575 bytes

	pass='generate-insecure-random-number'

	dd if=/dev/zero bs=1048575 count=1 status=none |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:"${pass}" |
		tail --bytes=16 |
		xxd -ps

	dd if=/dev/zero bs=1048575 count=1 status=none |
		openssl \
			enc \
			-nosalt \
			-aes-256-ctr \
			-pass pass:"${pass}" |
		$rtm run "${wsm}" |
		tail --bytes=32

	echo
}

ex1
ex2
ex3
