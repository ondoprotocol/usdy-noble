cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r github.com/ondoprotocol/usdy-noble/v2/* ./
cp -r api/aura/* api/
find api/ -type f -name "*.go" -exec sed -i 's|github.com/ondoprotocol/usdy-noble/v2/api/aura|github.com/ondoprotocol/usdy-noble/v2/api|g' {} +

rm -rf github.com
rm -rf api/aura
rm -rf aura
