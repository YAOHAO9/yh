
for i in {1..10}
do
  ts-node test-client/src/index.ts $i > $i &
done