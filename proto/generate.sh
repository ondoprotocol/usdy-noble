cd proto
buf generate --template buf.gen.gogo.yaml
cd ..

cp -r github.com/ondoprotocol/usdy-noble/* ./
rm -rf github.com
