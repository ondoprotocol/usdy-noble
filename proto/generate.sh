cd proto
buf generate --template buf.gen.gogo.yaml
cd ..

cp -r github.com/ondoprotocol/usdy-noble/v2/* ./
rm -rf github.com
