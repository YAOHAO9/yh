go run example.go -s="dwc" -i="connector-1" -k="connector" -H="127.0.0.1" -P="3110" -c -t="ksYNdrAo" -l="1" --zhost="127.0.0.1" --zport="2181" &
go run example.go -s="dwc" -i="connector-2" -k="connector" -P="3111" -c -t="ksYNdrAo" &
go run example.go -s="dwc" -i="ddz-1" -k="ddz" -P="3121" -t="ksYNdrAo" &
go run example.go -s="dwc" -i="ddz-2" -k="ddz" -P="3122" -t="ksYNdrAo" &
go run example.go -s="dwc" -i="ddz-3" -k="ddz" -P="3123" -t="ksYNdrAo" &
go run example.go -s="dwc" -i="ebg-1" -k="ebg" -P="3130" -t="ksYNdrAo" &